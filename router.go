package gmvc

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"
)

var (
	regexPart           *regexp.Regexp = regexp.MustCompile("/((?:[^/{}]|{[^}]*})*)")
	regexGlob                          = regexp.MustCompile("\\*{1,2}|{(?:[^}]|\\\\})*}")
	regexHandlerPattern                = regexp.MustCompile("^(?:(\\S+)\\s+){0,1}(/\\S*)\\s*$")
)

var (
	handlerPatternSyntax = "[HttpMethods] <UrlPattern>"
)

type route interface {
	match(c *Context, urlpath string, vars PathVars) (bool, error)
}

type Router struct {
	filters []*filter
	routes  []route
}

func NewRouter() *Router {
	return &Router{
		filters: make([]*filter, 0),
		routes:  make([]route, 0),
	}
}

func (rt *Router) Subrouter(pattern string) (*Router, error) {
	if pattern == "" || pattern == "/" {
		return rt, nil
	}

	srt := NewRouter()
	sr, err := newSubroutes(pattern, srt)
	if err != nil {
		return nil, err
	}

	rt.routes = append(rt.routes, sr)
	return srt, nil
}

func (rt *Router) Filter(pattern string, filter Filter) error {
	f, err := newFilter(pattern, filter)
	if err != nil {
		return err
	}

	rt.filters = append(rt.filters, f)
	return nil
}

func (rt *Router) FilterFunc(pattern string, f FilterFunc) error {
	return rt.Filter(pattern, f)
}

func (rt *Router) Handle(pattern string, handler Handler) error {
	values := regexHandlerPattern.FindStringSubmatch(pattern)
	if values == nil {
		return fmt.Errorf("incorrect format pattern for handler: %s, syntax: %s", pattern, handlerPatternSyntax)
	}

	var methods []string
	var pathPattern string
	var route *handlerRoute

	if values[1] == "" {
		methods = []string{"*"}
	} else {
		methods = strings.Split(values[1], ",")
	}
	pathPattern = values[2]

	for _, r := range rt.routes {
		if hr, ok := r.(*handlerRoute); ok {
			if hr.pattern == pathPattern {
				route = hr
				break
			}
		}
	}

	if route == nil {
		hr, err := newHandlerRoute(pathPattern)
		if err != nil {
			return err
		}
		route = hr
		rt.routes = append(rt.routes, route)
	}

	for _, method := range methods {
		if method == "" {
			continue
		}
		method = strings.ToUpper(method)
		if err := route.handle(method, handler); err != nil {
			return err
		}
	}

	return nil
}

func (rt *Router) HandleFunc(pattern string, f HandlerFunc) error {
	return rt.Handle(pattern, f)
}

func (rt *Router) route(c *Context, urlpath string, vars PathVars) (bool, error) {
	return newChain(c, urlpath, rt.filters, rt.routes, vars).next()
}

type chain struct {
	context *Context
	urlpath string
	filters []*filter
	routes  []route
	vars    PathVars
	pos     int
	tail    bool
}

func newChain(context *Context, urlpath string, filters []*filter, routes []route, vars PathVars) *chain {
	return &chain{
		context: context,
		urlpath: urlpath,
		filters: filters,
		routes:  routes,
		vars:    vars,
	}
}

func (c *chain) next() (bool, error) {
	if c.tail {
		return false, nil
	}

	if c.pos == len(c.filters) {
		c.tail = true
		for _, route := range c.routes {
			if match, err := route.match(c.context, c.urlpath, c.vars); match {
				return true, err
			}
		}
		return false, nil
	}

	for c.pos < len(c.filters) {
		f := c.filters[c.pos]
		c.pos++

		if match, err := f.match(c); match {
			return true, err
		}
	}

	return c.next()
}

type filter struct {
	pattern string
	tpl     *pathTemplate
	filter  Filter
}

func newFilter(pattern string, f Filter) (*filter, error) {
	var tpl *pathTemplate

	if pattern != "" {
		t, err := newPathTemplate(pattern, false)
		if err != nil {
			return nil, err
		}
		tpl = t
	}

	return &filter{
		pattern: pattern,
		tpl:     tpl,
		filter:  f,
	}, nil
}

func (f *filter) match(chain *chain) (bool, error) {
	vars := make(PathVars)
	if f.tpl == nil || f.tpl.match(chain.urlpath, vars) {
		c := newFilterContext(chain, vars)
		if chain.vars != nil {
			for k, v := range chain.vars {
				if _, ok := c.Vars[k]; !ok {
					c.Vars[k] = v
				}
			}
		}
		return c.match, f.filter.DoFilter(c)
	}
	return false, nil
}

type FilterContext struct {
	chain   *chain
	Context *Context
	Vars    PathVars
	match   bool
	next    bool
}

func newFilterContext(chain *chain, vars PathVars) *FilterContext {
	return &FilterContext{
		chain:   chain,
		Context: chain.context,
		Vars:    vars,
	}
}

func (fc *FilterContext) Next() error {
	if fc.next {
		return errors.New("multiple FilterContext.Next calls")
	}

	fc.next = true
	match, err := fc.chain.next()
	fc.match = match
	if match {
		return err
	}
	return nil
}

type subroutes struct {
	tpl    *pathTemplate
	router *Router
}

func newSubroutes(pattern string, router *Router) (*subroutes, error) {
	tpl, err := newPathTemplate(pattern, true)
	if err != nil {
		return nil, err
	}

	return &subroutes{
		tpl:    tpl,
		router: router,
	}, nil
}

func (sr *subroutes) match(c *Context, urlpath string, vars PathVars) (bool, error) {
	if sr.tpl.hasVars && vars == nil {
		vars = make(PathVars)
	}

	match, suffix := sr.tpl.matchPrefix(urlpath, vars)
	if !match {
		return false, nil
	}
	urlpath = path.Join("/", suffix)
	return sr.router.route(c, urlpath, vars)
}

type handlerRoute struct {
	pattern  string
	tpl      *pathTemplate
	handlers map[string]Handler
}

func newHandlerRoute(pattern string) (*handlerRoute, error) {
	tpl, err := newPathTemplate(pattern, false)
	if err != nil {
		return nil, err
	}

	return &handlerRoute{
		tpl:      tpl,
		handlers: make(map[string]Handler),
	}, nil
}

func (hr *handlerRoute) handle(method string, handler Handler) error {
	if hr.handlers[method] != nil {
		return fmt.Errorf("Conflicting handler methods mapped for pattern '%s'", hr.pattern)
	}
	hr.handlers[method] = handler
	return nil
}

func (hr *handlerRoute) match(c *Context, urlpath string, vars PathVars) (bool, error) {
	if hr.tpl.hasVars && vars == nil {
		vars = make(PathVars)
	}

	if hr.tpl.match(urlpath, vars) {
		if vars != nil {
			c.Vars = vars
		}

		h := hr.handlers[strings.ToUpper(c.Request.Method)]
		if h == nil {
			h = hr.handlers["*"]
		}
		if h == nil {
			errorStatus(c, http.StatusMethodNotAllowed)
			return true, nil
		}

		return true, h.HandleRequest(c)
	}

	return false, nil
}

type pathTemplate struct {
	regex   *regexp.Regexp
	prefix  bool
	hasVars bool
	vars    map[string]int
}

func newPathTemplate(pattern string, prefix bool) (*pathTemplate, error) {
	submatches := regexPart.FindAllStringSubmatch(pattern, -1)
	parts := make([]string, len(submatches))

	for i, m := range submatches {
		parts[i] = m[1]
	}

	buf := new(bytes.Buffer)
	buf.WriteString("^")

	for _, part := range parts {
		if part == "" {
			continue
		}

		buf.WriteString("/")
		locs := regexGlob.FindAllStringIndex(part, -1)
		if len(locs) == 0 {
			buf.WriteString(regexp.QuoteMeta(part))
		} else {
			var s, e int
			for j, loc := range locs {
				buf.WriteString(regexp.QuoteMeta(part[e:loc[0]]))

				s, e = loc[0], loc[1]
				g := part[s:e]

				switch {
				case g == "?":
					buf.WriteString("[^/]+")

				case g == "**":
					buf.WriteString(".+")

				case strings.HasPrefix(g, "{") && strings.HasSuffix(g, "}"):
					g = strings.Replace(g, "\\}", "}", -1)
					var k, v string
					i := strings.Index(g, ":")
					if i == -1 {
						k = g[1 : len(g)-1]
						v = "[^/]+"
					} else {
						k = g[1:i]
						v = g[i+1 : len(g)-1]
					}

					if k == "" {
						buf.WriteString("(?:" + v + ")")
					} else {
						buf.WriteString("(?P<" + regexp.QuoteMeta(k) + ">" + v + ")")
					}
				}

				if j == len(locs)-1 {
					buf.WriteString(regexp.QuoteMeta(part[e:]))
				}
			}
		}
	}

	if prefix {
		buf.WriteString("(|\\/.*)$")
	} else {
		buf.WriteString("\\/?$")
	}

	regex, err := regexp.Compile(buf.String())
	if err != nil {
		return nil, err
	}

	vars := make(map[string]int)
	for i, n := range regex.SubexpNames() {
		if n != "" {
			vars[n] = i
		}
	}

	return &pathTemplate{
		regex:   regex,
		prefix:  prefix,
		hasVars: len(vars) > 0,
		vars:    vars,
	}, nil
}

func (t *pathTemplate) match(urlpath string, vars PathVars) bool {
	if vars == nil {
		return t.regex.MatchString(urlpath)
	}

	submatch := t.regex.FindStringSubmatch(urlpath)
	if submatch == nil {
		return false
	}

	if t.hasVars && vars != nil {
		for k, v := range t.vars {
			vars[k] = submatch[v]
		}
	}

	return true
}

func (t *pathTemplate) matchPrefix(urlpath string, vars PathVars) (bool, string) {
	if !t.prefix {
		return false, ""
	}

	submatch := t.regex.FindStringSubmatch(urlpath)
	if submatch == nil {
		return false, ""
	}

	if t.hasVars && vars != nil {
		for k, v := range t.vars {
			vars[k] = submatch[v]
		}
	}

	return true, submatch[len(submatch)-1]
}

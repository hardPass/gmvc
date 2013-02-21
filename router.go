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
	part *regexp.Regexp = regexp.MustCompile("/((?:[^/{}]|{[^}]*})*)")
	glob                = regexp.MustCompile("\\*{1,2}|{(?:[^}]|\\\\})*}")
)

type route interface {
	do(c *Context, urlpath string) (bool, error)
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

func (rt *Router) route(c *Context, urlpath string) (bool, error) {
	return newRouteChain(c, urlpath, rt.filters, rt.routes).next()
}

func (rt *Router) Subrouter(pattern string) (*Router, error) {
	if pattern == "" || pattern == "/" {
		return rt, nil
	}

	subrt := NewRouter()
	sm, err := newSubroutes(pattern, subrt)
	if err != nil {
		return nil, err
	}

	rt.routes = append(rt.routes, sm)
	return subrt, nil
}

func (rt *Router) Filter(pattern string, filter Filter) error {
	fr, err := newFilter(pattern, filter)
	if err != nil {
		return err
	}

	rt.filters = append(rt.filters, fr)
	return nil
}

func (rt *Router) FilterFunc(pattern string, f FilterFunc) error {
	return rt.Filter(pattern, f)
}

func (rt *Router) Handle(pattern string, handler Handler) error {
	var methods []string
	var pathPattern string
	var route *handlerRoute

	s := strings.SplitN(pattern, " ", 2)

	if len(s) > 1 {
		methods = strings.Split(s[0], ",")
		pathPattern = s[1]
	} else {
		methods = []string{"*"}
		pathPattern = s[0]
	}

	pathPattern = path.Join("/", strings.Trim(pathPattern, " "))

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

type FilterContext struct {
	chain   *routeChain
	Context *Context
	Vars    PathVars
	next    bool
	hit     bool
}

func newFilterContext(chain *routeChain) *FilterContext {
	return &FilterContext{
		chain:   chain,
		Context: chain.context,
		Vars:    make(PathVars),
	}
}

func (fc *FilterContext) Next() error {
	if fc.next {
		return errors.New("multiple FilterContext.Next calls")
	}

	fc.next = true
	hit, err := fc.chain.next()
	fc.hit = hit
	if hit {
		return err
	}
	return nil
}

type routeChain struct {
	context *Context
	urlpath string
	filters []*filter
	routes  []route
	pos     int
	tail    bool
}

func newRouteChain(context *Context, urlpath string, filters []*filter, routes []route) *routeChain {
	return &routeChain{
		context: context,
		urlpath: urlpath,
		filters: filters,
		routes:  routes,
	}
}

func (rc *routeChain) next() (bool, error) {
	if rc.tail {
		return false, nil
	}

	if rc.pos == len(rc.filters) {
		rc.tail = true
		for _, route := range rc.routes {
			if hit, err := route.do(rc.context, rc.urlpath); hit {
				return true, err
			}
		}
		return false, nil
	}

	for rc.pos < len(rc.filters) {
		fr := rc.filters[rc.pos]
		rc.pos++

		if hit, err := fr.do(rc); hit {
			return true, err
		}
	}

	return rc.next()
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

func (fr *filter) do(chain *routeChain) (bool, error) {
	fc := newFilterContext(chain)
	if fr.tpl == nil || fr.tpl.match(chain.urlpath, fc.Vars) {
		err := fr.filter.DoFilter(fc)
		return fc.hit, err
	}
	return false, nil
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

func (sr *subroutes) do(c *Context, urlpath string) (bool, error) {
	match, suffix := sr.tpl.matchPrefix(urlpath, c.Vars)
	if !match {
		return false, nil
	}
	urlpath = path.Join("/", suffix)
	return sr.router.route(c, urlpath)
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

func (m *handlerRoute) handle(method string, handler Handler) error {
	if m.handlers[method] != nil {
		return fmt.Errorf("Conflicting handler methods mapped for pattern '%s'", m.pattern)
	}
	m.handlers[method] = handler
	return nil
}

func (m *handlerRoute) do(c *Context, urlpath string) (bool, error) {
	if m.tpl.match(urlpath, c.Vars) {
		method := strings.ToUpper(c.Request.Method)
		h := m.handlers[method]
		if h == nil {
			h = m.handlers["*"]
		}
		if h == nil {
			status := http.StatusMethodNotAllowed
			c.ErrorStatus(errors.New(http.StatusText(status)), status)
			return true, nil
		}
		err := h.HandleHTTP(c)
		return true, err
	}

	return false, nil
}

type pathTemplate struct {
	regex  *regexp.Regexp
	vars   map[string]int
	prefix bool
}

func newPathTemplate(pattern string, prefix bool) (*pathTemplate, error) {
	partsubs := part.FindAllStringSubmatch(pattern, -1)
	parts := make([]string, len(partsubs))

	for i, sub := range partsubs {
		parts[i] = sub[1]
	}

	buf := new(bytes.Buffer)
	buf.WriteString("^/")

	for i, part := range parts {
		locs := glob.FindAllStringIndex(part, -1)

		if len(locs) == 0 {
			buf.WriteString(regexp.QuoteMeta(part))
		} else {
			var s, e int
			for j, loc := range locs {
				buf.WriteString(regexp.QuoteMeta(part[e:loc[0]]))

				s, e = loc[0], loc[1]
				g := part[s:e]

				switch {
				case g == "*":
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

		if i < len(parts)-1 {
			buf.WriteString("/")
		}
	}

	if prefix {
		buf.WriteString("(|\\/.*)$")
	} else {
		buf.WriteString("\\/*$")
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

	t := &pathTemplate{
		regex:  regex,
		vars:   vars,
		prefix: prefix,
	}

	return t, nil

}

func (t *pathTemplate) match(urlpath string, vars PathVars) bool {
	if vars == nil {
		return t.regex.MatchString(urlpath)
	}

	sub := t.regex.FindStringSubmatch(urlpath)
	if sub == nil {
		return false
	}

	for k, v := range t.vars {
		vars[k] = sub[v]
	}
	return true
}

func (t *pathTemplate) matchPrefix(urlpath string, vars PathVars) (bool, string) {
	if !t.prefix {
		return false, ""
	}

	sub := t.regex.FindStringSubmatch(urlpath)
	if sub == nil {
		return false, ""
	}

	if vars != nil {
		for k, v := range t.vars {
			vars[k] = sub[v]
		}
	}

	suffix := sub[len(sub)-1]
	return true, suffix
}

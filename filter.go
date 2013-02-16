package gmvc

import (
	"errors"
)

type Filter interface {
	DoFilter(fc *FilterContext) error
}

type FilterFunc func(*FilterContext) error

func (f FilterFunc) DoFilter(fc *FilterContext) error {
	return f(fc)
}

type FilterContext struct {
	chain   *routeChain
	Context *Context
	Vars    PathVars
	next    bool
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
	ok, err := fc.chain.next()
	if ok {
		return err
	}
	return nil
}

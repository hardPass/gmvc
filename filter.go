package gmvc

type Filter interface {
	DoFilter(fc *FilterContext) error
}

type FilterFunc func(*FilterContext) error

func (f FilterFunc) DoFilter(fc *FilterContext) error {
	return f(fc)
}

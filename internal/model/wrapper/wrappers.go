package wrapper

type Wrappers map[WrapperType]*Wrapper

func (w Wrappers) Add(key WrapperType, value *Wrapper) {
	w[key] = value
}

func (w Wrappers) Get(key WrapperType) (*Wrapper, bool) {
	wrapper, ok := w[key]
	return wrapper, ok
}

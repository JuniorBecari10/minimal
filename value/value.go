package value

type Value interface {
	value()
}

// ---

type ValueNumber struct {
	Value float64
}

type ValueBool struct {
	Value bool
}

// ---

func (x ValueNumber) value() {}
func (x ValueBool) value()   {}

package value

type Upvalue struct {
	Locals *[]Value
	Index  int // filled when open.

	ClosedValue Value // filled when closed.
	IsClosed    bool
}

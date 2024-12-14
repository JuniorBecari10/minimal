package value

type Upvalue struct {
	Location    *Value // filled if open.
	ClosedValue Value  // filled when closed.

	IsClosed bool
}

package value

type Upvalue struct {
	LocalsIndex int
	Index       int // filled when open.

	ClosedValue Value // filled when closed.
	IsClosed    bool
}

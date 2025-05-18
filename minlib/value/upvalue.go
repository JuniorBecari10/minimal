package value

type Upvalue struct {
	LocalsIndex  int   // index in the locals array (if open)
	UpvalueIndex int   // index in the upvalue list (filled when open)

	ClosedValue  Value // captured value (used when closed)
	IsClosed     bool  // whether the variable is closed (moved to heap, which is the field above)
}


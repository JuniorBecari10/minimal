package chunk

import (
	"vm-go/token"
	"vm-go/value"
)

// Because of import cycle, Chunk is a type alias, so it can be redeclared without importing it
type Chunk = struct {
	Code      []byte
	Constants []value.Value

	Positions []token.Position
}

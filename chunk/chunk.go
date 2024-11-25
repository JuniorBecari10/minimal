package chunk

import (
	"vm-go/token"
	"vm-go/value"
)

type Chunk struct {
	Code      []byte
	Constants []value.Value

	Positions []token.Position
}

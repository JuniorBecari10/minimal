package value

import "minlib/token"

type Chunk struct {
	Code      []byte
	Constants []Value

	Metadata []Metadata
}

type Metadata struct {
	Position token.Position
	Length uint32
}

func NewMetaLen1(pos token.Position) Metadata {
	return Metadata{
		Position: pos,
		Length: 1,
	}
}


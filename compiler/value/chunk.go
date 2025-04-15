package value

import "vm-go/token"

type Chunk struct {
	Code      []byte
	Constants []Value

	Metadata []ChunkMetadata
}

type ChunkMetadata struct {
	Position token.Position
	Length int
}

func NewMetaLen1(pos token.Position) ChunkMetadata {
	return ChunkMetadata{
		Position: pos,
		Length: 1,
	}
}

// ---

func (c *Chunk) Serialize() []byte {
	return []byte{}
}

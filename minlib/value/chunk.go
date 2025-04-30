package value

import "minlib/token"

type Chunk struct {
	Code      []byte
	Constants []Value

	Metadata []ChunkMetadata
}

type ChunkMetadata struct {
	Position token.Position
	Length uint32
}

func NewMetaLen1(pos token.Position) ChunkMetadata {
	return ChunkMetadata{
		Position: pos,
		Length: 1,
	}
}


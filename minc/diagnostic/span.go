package diagnostic

import "minlib/token"

// it will only support one line
type Span struct {
	Pos token.Position
	Length int
}

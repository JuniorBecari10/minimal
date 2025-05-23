package diagnostic

import "minlib/token"

type Span struct {
	Start token.Position
	End token.Position
}

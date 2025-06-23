package types

import "minlib/token"

// Use this with caution.
func DummyType(data TypeData) Type {
	return Type{
		Token: token.StartToken(),
		Data: data,
	}
}

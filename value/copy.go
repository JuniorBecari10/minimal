package value

import (
	"fmt"
	"reflect"
	"vm-go/util"
)

// Deep-copies a Value.
func CopyValue(value Value) Value {
	switch v := value.(type) {
		case ValueNumber:
			return ValueNumber{Value: v.Value}
		case ValueString:
			return ValueString{Value: v.Value}
		case ValueBool:
			return ValueBool{Value: v.Value}
		
		case ValueNil:
			return ValueNil{}
		case ValueVoid:
			return ValueVoid{}
			
		case ValueFunction: {
			return ValueFunction{
				Arity: v.Arity,
				Chunk: Chunk{
					Code: util.CopyList(v.Chunk.Code, func(b byte) byte {
						return b
					}),
					Constants: util.CopyList(v.Chunk.Constants, CopyValue),
					Metadata: util.CopyList(v.Chunk.Metadata, func(m ChunkMetadata) ChunkMetadata {
						return m
					}),
				},
				Name: v.Name,
			}
		}

		case ValueClosure: {
			// dereference the function first
			fn, _ := CopyValue(*v.Fn).(ValueFunction)

			return ValueClosure{
				Fn: &fn,
				Upvalues: util.CopyList(v.Upvalues, func(up *Upvalue) *Upvalue {
					// copy the pointer
					return up
				}),
			}
		}

		case ValueNativeFn: {
			return ValueNativeFn{
				Arity: v.Arity,
				Fn: v.Fn,
			}
		}

        case ValueRange: {
            return ValueRange{
                Start: v.Start,
                End: v.End,
                Step: v.Step,
            }
        }

		case ValueRecord: {
			return ValueRecord{
				FieldNames: util.CopyList(v.FieldNames, func(s string) string {
					return s
				}),
				Name: v.Name,
			}
		}

		case ValueInstance: {
			return ValueInstance{
				Fields: util.CopyList(v.Fields, CopyValue),
				Record: v.Record, // don't copy
			}
		}
		
		case ValueBoundMethod: {
			return ValueBoundMethod{
				Receiver: v.Receiver, // don't copy
				Method: CopyValue(v.Method).(ValueClosure),
			}
		}

		default:
			panic(fmt.Sprintf("Unknown value: '%v': %s", v, reflect.TypeOf(v)))
	}
}

package value

import (
	"bytes"
	"encoding/binary"
)

/*
	File Structure:
	[code len] [code] [constants len] [constants] [metadata len] [metadata]
*/
func (c *Chunk) Serialize() []byte {
	buf := new(bytes.Buffer)

	// Serialize Code
	binary.Write(buf, binary.LittleEndian, int32(len(c.Code)))
	buf.Write(c.Code)

	// Serialize Constants
	binary.Write(buf, binary.LittleEndian, int32(len(c.Constants)))
	for _, v := range c.Constants {
		serializeValue(buf, v)
	}

	// Serialize Metadata
	binary.Write(buf, binary.LittleEndian, int32(len(c.Metadata)))
	for _, m := range c.Metadata {
		binary.Write(buf, binary.LittleEndian, int32(m.Position.Line))
		binary.Write(buf, binary.LittleEndian, int32(m.Position.Col))
		binary.Write(buf, binary.LittleEndian, int32(m.Length))
	}

	return buf.Bytes()
}

func serializeValue(buf *bytes.Buffer, v Value) {
	switch val := v.(type) {
		case ValueNumber: {
			buf.WriteByte(1)
			binary.Write(buf, binary.LittleEndian, val.Value)
		}

		case ValueString: {
			buf.WriteByte(2)
			serializeString(buf, val.Value)
		}

		case ValueBool: {
			buf.WriteByte(3)

			if val.Value {
				buf.WriteByte(1)
			} else {
				buf.WriteByte(0)
			}
		}

		case ValueNil:
			buf.WriteByte(4)

		case ValueVoid:
			buf.WriteByte(5)

		case *ValueFunction: {
			buf.WriteByte(6)
			binary.Write(buf, binary.LittleEndian, int32(val.Arity))
			
			if val.Name != nil {
				buf.WriteByte(1)
				serializeString(buf, *val.Name)
			} else {
				buf.WriteByte(0)
			}
			
			chunkData := val.Chunk.Serialize()
			
			binary.Write(buf, binary.LittleEndian, int32(len(chunkData)))
			buf.Write(chunkData)
		}

		case *ValueClosure: {
			buf.WriteByte(7)
			serializeValue(buf, val.Fn)

			binary.Write(buf, binary.LittleEndian, int32(len(val.Upvalues)))
			for _, up := range val.Upvalues {
				serializeUpvalue(buf, up)
			}
		}

		case *ValueRange: {
			buf.WriteByte(8)
			binary.Write(buf, binary.LittleEndian, val.Start)
			binary.Write(buf, binary.LittleEndian, val.End)
			binary.Write(buf, binary.LittleEndian, val.Step)
			
			if val.Inclusive {
				buf.WriteByte(1)
			} else {
				buf.WriteByte(0)
			}
		}

		case *ValueRecord: {
			buf.WriteByte(9)
			serializeString(buf, val.Name)
			
			binary.Write(buf, binary.LittleEndian, int32(len(val.FieldNames)))
			for _, f := range val.FieldNames {
				serializeString(buf, f)
			}
			
			binary.Write(buf, binary.LittleEndian, int32(len(val.Methods)))
			for _, m := range val.Methods {
				serializeValue(buf, &m)
			}
		}

		case *ValueInstance: {
			buf.WriteByte(10)
			binary.Write(buf, binary.LittleEndian, int32(len(val.Fields)))
			
			for _, f := range val.Fields {
				serializeValue(buf, f)
			}
			
			serializeValue(buf, val.Record)
		}

		case *ValueBoundMethod: {
			buf.WriteByte(11)
			
			serializeValue(buf, val.Receiver)
			serializeValue(buf, &val.Method)
		}

		// ValueNativeFunction cannot be serialized,
		// but that's not a problem, because the VM needs to reimplement them
	}
}

func serializeString(buf *bytes.Buffer, s string) {
	binary.Write(buf, binary.LittleEndian, int32(len(s)))
	buf.WriteString(s)
}

func serializeUpvalue(buf *bytes.Buffer, up Upvalue) {
	binary.Write(buf, binary.LittleEndian, int32(up.LocalsIndex))
	binary.Write(buf, binary.LittleEndian, int32(up.Index))
	
	if up.IsClosed {
		buf.WriteByte(1)
		serializeValue(buf, up.ClosedValue)
	} else {
		buf.WriteByte(0)
	}
}


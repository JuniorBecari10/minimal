package value

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Upvalues do not need a code because they will only exist inside closures.
const (
	IntCode = iota
	FloatCode
	StringCode
	CharCode
	BoolCode
	NilCode
	VoidCode
	FunctionCode
	ClosureCode
	RangeCode
	RecordCode
	InstanceCode
	BoundMethodCode
)

/*
	File Structure (parts aren't separated by newlines):

	[header]
	[name len] [name]
	[code len] [code]
	[constants len] [constants]
	[metadata len] [metadata]
	[footer]

	Where [header] is:
	"MNML"

	Where [footer] is:
	<checksum of the entire file, including [header]>
*/
func (c *Chunk) Serialize() []byte {
	buf := new(bytes.Buffer)

	writeName(buf, c)
	writeCode(buf, c)
	writeConstants(buf, c)
	writeMetadata(buf, c)
	
	return buf.Bytes()
}

// ---

func writeName(buf *bytes.Buffer, c *Chunk) {
	binary.Write(buf, binary.LittleEndian, uint32(len(c.Name)))
	buf.Write([]byte(c.Name))
}

func writeCode(buf *bytes.Buffer, c *Chunk) {
	binary.Write(buf, binary.LittleEndian, uint32(len(c.Code)))
	buf.Write(c.Code)
}

func writeConstants(buf *bytes.Buffer, c *Chunk) {
	binary.Write(buf, binary.LittleEndian, uint32(len(c.Constants)))
	for _, v := range c.Constants {
		serializeValue(buf, v)
	}
}

func writeMetadata(buf *bytes.Buffer, c *Chunk) {
	binary.Write(buf, binary.LittleEndian, uint32(len(c.Metadata)))
	for _, m := range c.Metadata {
		serializeMetadata(buf, m)
	}
}

// ---

func serializeValue(buf *bytes.Buffer, v Value) {
	switch val := v.(type) {
		case ValueInt: {
			buf.WriteByte(IntCode)
			binary.Write(buf, binary.LittleEndian, val.Value)
		}
		
		case ValueFloat: {
			buf.WriteByte(FloatCode)
			binary.Write(buf, binary.LittleEndian, val.Value)
		}

		case ValueString: {
			buf.WriteByte(StringCode)
			serializeString(buf, val.Value)
		}
		
		case ValueChar: {
			buf.WriteByte(CharCode)
			binary.Write(buf, binary.LittleEndian, val.Value)
		}

		case ValueBool: {
			buf.WriteByte(BoolCode)
			if val.Value {
				buf.WriteByte(1)
			} else {
				buf.WriteByte(0)
			}
		}

		case ValueNil:
			buf.WriteByte(NilCode)

		case ValueVoid:
			buf.WriteByte(VoidCode)

		case ValueFunction: {
			buf.WriteByte(FunctionCode)
			binary.Write(buf, binary.LittleEndian, val.Arity)

			if val.Name != nil {
				buf.WriteByte(1)
				serializeString(buf, *val.Name)
			} else {
				buf.WriteByte(0)
			}

			buf.Write(val.Chunk.Serialize())
		}

		case ValueClosure: {
			buf.WriteByte(ClosureCode)
			serializeValue(buf, val.Fn)
			
			binary.Write(buf, binary.LittleEndian, uint32(len(val.Upvalues)))
			
			for _, up := range val.Upvalues {
				serializeUpvalue(buf, up)
			}
		}

		case ValueRange: {
			buf.WriteByte(RangeCode)
			
			binary.Write(buf, binary.LittleEndian, val.Start)
			binary.Write(buf, binary.LittleEndian, val.End)
			binary.Write(buf, binary.LittleEndian, val.Step)
			
			if val.Inclusive {
				buf.WriteByte(1)
			} else {
				buf.WriteByte(0)
			}
		}

		case ValueRecord: {
			buf.WriteByte(RecordCode)
			
			serializeString(buf, val.Name)
			binary.Write(buf, binary.LittleEndian, uint32(len(val.FieldNames)))
			
			for _, f := range val.FieldNames {
				serializeString(buf, f)
			}
			
			binary.Write(buf, binary.LittleEndian, uint32(len(val.Methods)))
			for i := range val.Methods {
				serializeValue(buf, &val.Methods[i])
			}
		}

		case ValueInstance: {
			buf.WriteByte(InstanceCode)
			binary.Write(buf, binary.LittleEndian, uint32(len(val.Fields)))
			
			for _, f := range val.Fields {
				serializeValue(buf, f)
			}
			
			serializeValue(buf, val.Record)
		}

		case ValueBoundMethod: {
			buf.WriteByte(BoundMethodCode)
			
			serializeValue(buf, val.Receiver)
			serializeValue(buf, &val.Method)
		}

		// ValueNativeFunction cannot be serialized,
		// but that's not a problem, because the VM needs to reimplement them

		default:
			panic(fmt.Sprintf("Unknown type: %T\n", v))
	}
}

func serializeString(buf *bytes.Buffer, s string) {
	err := binary.Write(buf, binary.LittleEndian, uint32(len(s)))
	if err != nil {
		panic(fmt.Sprintf("failed to write string length: %v", err))
	}
	buf.WriteString(s)
}

func serializeUpvalue(buf *bytes.Buffer, up Upvalue) {
	err := binary.Write(buf, binary.LittleEndian, uint32(up.LocalsIndex))
	if err != nil {
		panic(fmt.Sprintf("failed to write upvalue LocalsIndex: %v", err))
	}
	err = binary.Write(buf, binary.LittleEndian, uint32(up.Index))
	if err != nil {
		panic(fmt.Sprintf("failed to write upvalue Index: %v", err))
	}

	if up.IsClosed {
		buf.WriteByte(1)
		serializeValue(buf, up.ClosedValue)
	} else {
		buf.WriteByte(0)
	}
}

func serializeMetadata(buf *bytes.Buffer, m Metadata) {
	binary.Write(buf, binary.LittleEndian, uint32(m.Position.Line))
	binary.Write(buf, binary.LittleEndian, uint32(m.Position.Col))
	binary.Write(buf, binary.LittleEndian, uint32(m.Length))
}


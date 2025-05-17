package value

import (
	"bytes"
	"io"
	"encoding/binary"
	"fmt"
	"minlib/token"
)

func Deserialize(data []byte) Chunk {
	return readChunk(bytes.NewBuffer(data))
}

func readChunk(r io.Reader) Chunk {
	chunk := Chunk{}

	readName(r, &chunk)
	readCode(r, &chunk)
	readConstants(r, &chunk)
	readMetadata(r, &chunk)

	return chunk
}

// ---

func readName(r io.Reader, chunk *Chunk) {
	var nameLen uint32
	binary.Read(r, binary.LittleEndian, &nameLen)
	
	name := make([]byte, nameLen)
	r.Read(name)

	chunk.Name = string(name)
}

func readCode(r io.Reader, chunk *Chunk) {
	var codeLen uint32
	binary.Read(r, binary.LittleEndian, &codeLen)

	chunk.Code = make([]byte, codeLen)
	r.Read(chunk.Code)
}

func readConstants(r io.Reader, chunk *Chunk) {
	var constCount uint32
	
	binary.Read(r, binary.LittleEndian, &constCount)
	chunk.Constants = make([]Value, constCount)
	
	for i := range chunk.Constants {
		chunk.Constants[i] = deserializeValue(r)
	}
}

func readMetadata(r io.Reader, chunk *Chunk) {
	var metaCount uint32
	
	binary.Read(r, binary.LittleEndian, &metaCount)
	chunk.Metadata = make([]Metadata, metaCount)
	
	for i := range chunk.Metadata {
		chunk.Metadata[i] = deserializeMetadata(r)
	}
}

// ---

func deserializeValue(r io.Reader) Value {
	buf := r.(*bytes.Buffer)
	tag, _ := buf.ReadByte()

	switch tag {
		case IntCode: {
			var i int32
			binary.Read(buf, binary.LittleEndian, &i)

			return ValueInt{Value: i}
		}
		
		case FloatCode: {
			var f float64
			binary.Read(buf, binary.LittleEndian, &f)

			return ValueFloat{Value: f}
		}
		
		case StringCode:
			return ValueString{Value: deserializeString(buf)}
		
		case CharCode: {
			var c uint8
			binary.Read(buf, binary.LittleEndian, &c)

			return ValueChar{Value: c}
		}
		
		case BoolCode: {
			var b bool
			binary.Read(buf, binary.LittleEndian, &b)
			
			return ValueBool{Value: b}
		}
		
		case NilCode:
			return ValueNil{}
		
		case VoidCode:
			return ValueVoid{}
		
		case FunctionCode: {
			var arity uint32
			binary.Read(buf, binary.LittleEndian, &arity)

			hasName, _ := buf.ReadByte()
			var name *string
			
			if hasName != 0 {
				s := deserializeString(buf)
				name = &s
			}
			
			fchunk := readChunk(buf)
			return ValueFunction{Arity: arity, Chunk: fchunk, Name: name}
		}

		case ClosureCode: {
			fn := deserializeValue(buf).(ValueFunction)
			
			var upcount int32
			binary.Read(buf, binary.LittleEndian, &upcount)
			
			upvalues := make([]Upvalue, upcount)
			for i := range upvalues {
				var li, idx int32
				var closed bool
				
				binary.Read(buf, binary.LittleEndian, &li)
				binary.Read(buf, binary.LittleEndian, &idx)
				binary.Read(buf, binary.LittleEndian, &closed)
				
				up := Upvalue{
					LocalsIndex: int(li),
					Index:       int(idx),
					IsClosed:    closed,
				}
				
				if closed {
					up.ClosedValue = deserializeValue(buf)
				}
				
				upvalues[i] = up
			}
			
			return ValueClosure{Fn: &fn, Upvalues: upvalues}
		}

		case RecordCode: {
			name := deserializeString(buf)
			var fieldCount uint32
			
			binary.Read(buf, binary.LittleEndian, &fieldCount)
			fields := make([]string, fieldCount)
			
			for i := range fields {
				fields[i] = deserializeString(buf)
			}
			
			var methodCount uint32
			
			binary.Read(buf, binary.LittleEndian, &methodCount)
			methods := make([]ValueClosure, methodCount)
			
			for i := range methods {
				methods[i] = deserializeValue(buf).(ValueClosure)
			}
			
			return ValueRecord{Name: name, FieldNames: fields, Methods: methods}
		}

		case RangeCode: {
			var start, end, step float64
			
			binary.Read(buf, binary.LittleEndian, &start)
			binary.Read(buf, binary.LittleEndian, &end)
			binary.Read(buf, binary.LittleEndian, &step)
			
			var inclusive bool
			binary.Read(buf, binary.LittleEndian, &inclusive)
			
			return ValueRange{Start: start, End: end, Step: step, Inclusive: inclusive}
		}

		case InstanceCode: {
			record := deserializeValue(buf).(ValueRecord)
			var fieldCount uint32
			
			binary.Read(buf, binary.LittleEndian, &fieldCount)
			fields := make([]Value, fieldCount)
			
			for i := range fields {
				fields[i] = deserializeValue(buf)
			}
			
			return ValueInstance{Fields: fields, Record: &record}
		}

		case BoundMethodCode: {
			receiver := deserializeValue(buf)
			method := deserializeValue(buf).(ValueClosure)
			
			return ValueBoundMethod{Receiver: receiver, Method: method}
		}

		default:
			panic(fmt.Sprintf("Unknown value tag: %d", tag))
	}
}

func deserializeString(r io.Reader) string {
	buf := r.(*bytes.Buffer)

	var strlen uint32
	binary.Read(buf, binary.LittleEndian, &strlen)
	
	strbytes := make([]byte, strlen)
	buf.Read(strbytes)
	
	return string(strbytes)
}

func deserializeMetadata(r io.Reader) Metadata {
	buf := r.(*bytes.Buffer)
	var line, col, length uint32

	binary.Read(buf, binary.LittleEndian, &line)
	binary.Read(buf, binary.LittleEndian, &col)
	binary.Read(buf, binary.LittleEndian, &length)
	
	return Metadata{Position: token.Position{Line: uint32(line), Col: uint32(col)}, Length: uint32(length)}
}

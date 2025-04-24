package value

import (
	"bytes"
	"io"
	"encoding/binary"
	"fmt"
	"minlib/token"
)

func Deserialize(data []byte) Chunk {
	buf := bytes.NewBuffer(data)
	chunk := Chunk{}

	// Read code
	var codeLen int32
	
	binary.Read(buf, binary.LittleEndian, &codeLen)
	chunk.Code = make([]byte, codeLen)
	buf.Read(chunk.Code)

	// Read constants
	var constCount int32
	
	binary.Read(buf, binary.LittleEndian, &constCount)
	chunk.Constants = make([]Value, constCount)
	
	for i := range chunk.Constants {
		chunk.Constants[i] = deserializeValue(buf)
	}

	// Read metadata
	var metaCount int32
	
	binary.Read(buf, binary.LittleEndian, &metaCount)
	chunk.Metadata = make([]ChunkMetadata, metaCount)
	
	for i := range chunk.Metadata {
		chunk.Metadata[i] = readMetadata(buf)
	}

	return chunk
}

func deserializeValue(r io.Reader) Value {
	buf := r.(*bytes.Buffer)
	tag, _ := buf.ReadByte()

	switch tag {
		case 1: { // Number
			var f float64
			binary.Read(buf, binary.LittleEndian, &f)

			return ValueNumber{Value: f}
		}
		
		case 2: // String
			return ValueString{Value: deserializeString(buf)}
		
		case 3: { // Bool
			var b bool
			binary.Read(buf, binary.LittleEndian, &b)
			
			return ValueBool{Value: b}
		}
		
		case 4: return ValueNil{}
		case 5: return ValueVoid{}
		
		case 6: { // Function
			var arity int32
			
			binary.Read(buf, binary.LittleEndian, &arity)
			fchunk := readChunk(buf)

			var hasName byte
			
			buf.ReadByte() // discard or check for error
			hasName, _ = buf.ReadByte()
			
			var name *string
			
			if hasName != 0 {
				s := deserializeString(buf)
				name = &s
			}
			
			return &ValueFunction{Arity: int(arity), Chunk: fchunk, Name: name}
		}

		case 7: { // Closure
			fn := deserializeValue(buf).(*ValueFunction)
			
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
			
			return ValueClosure{Fn: fn, Upvalues: upvalues}
		}

		case 8: {// Record
			name := deserializeString(buf)
			var fieldCount int32
			
			binary.Read(buf, binary.LittleEndian, &fieldCount)
			fields := make([]string, fieldCount)
			
			for i := range fields {
				fields[i] = deserializeString(buf)
			}
			
			var methodCount int32
			
			binary.Read(buf, binary.LittleEndian, &methodCount)
			methods := make([]ValueClosure, methodCount)
			
			for i := range methods {
				methods[i] = deserializeValue(buf).(ValueClosure)
			}
			
			return ValueRecord{Name: name, FieldNames: fields, Methods: methods}
		}

		case 9: {// Range
			var start, end, step float64
			
			binary.Read(buf, binary.LittleEndian, &start)
			binary.Read(buf, binary.LittleEndian, &end)
			binary.Read(buf, binary.LittleEndian, &step)
			
			var inclusive bool
			binary.Read(buf, binary.LittleEndian, &inclusive)
			
			return ValueRange{Start: start, End: end, Step: step, Inclusive: inclusive}
		}

		case 10: { // Instance
			record := deserializeValue(buf).(*ValueRecord)
			var fieldCount int32
			
			binary.Read(buf, binary.LittleEndian, &fieldCount)
			fields := make([]Value, fieldCount)
			
			for i := range fields {
				fields[i] = deserializeValue(buf)
			}
			
			return ValueInstance{Fields: fields, Record: record}
		}

		case 11: { // BoundMethod
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

	var strlen int32
	binary.Read(buf, binary.LittleEndian, &strlen)
	
	strbytes := make([]byte, strlen)
	buf.Read(strbytes)
	
	return string(strbytes)
}
func readChunk(r io.Reader) Chunk {
	buf := r.(*bytes.Buffer)

	var chunkLen int32
	binary.Read(buf, binary.LittleEndian, &chunkLen)

	chunkData := make([]byte, chunkLen)
	buf.Read(chunkData)

	return Deserialize(chunkData)
}
func readMetadata(r io.Reader) ChunkMetadata {
	buf := r.(*bytes.Buffer)
	var line, col, length int32

	binary.Read(buf, binary.LittleEndian, &line)
	binary.Read(buf, binary.LittleEndian, &col)
	binary.Read(buf, binary.LittleEndian, &length)
	
	return ChunkMetadata{Position: token.Position{Line: int(line), Col: int(col)}, Length: int(length)}
}

package vm

import (
	"fmt"
	"vm-go/compiler"
	"vm-go/util"
	"vm-go/value"
)

type InterpretResult int

const (
	STATUS_OK InterpretResult = iota
	STATUS_STACK_EMPTY
	STATUS_OUT_OF_BOUNDS
	STATUS_DIV_ZERO
	STATUS_TYPE_ERROR
	STATUS_INCORRECT_ARITY
	STATUS_PROPERTY_DOESNT_EXIST
    STATUS_UNREACHABLE_RANGE
)

type VM struct {
	currentChunk *value.Chunk
	topLevel     value.Chunk

	stack     []value.Value
	globals   []value.Value
	callStack []CallFrame
	openUpvalues  []*value.Upvalue // can be a linked list also, and it owns them.
 
	ip    int
	oldIp int

	hadError bool
	fileData *util.FileData
}

func NewVM(chunk value.Chunk, fileData *util.FileData) *VM {
	vm := VM{
		currentChunk: &chunk,
		topLevel: chunk,

		stack:     []value.Value{},
		globals:   []value.Value{},
		callStack: []CallFrame{},
		openUpvalues:  []*value.Upvalue{},

		ip:        0,
		oldIp:     0,

		hadError:  false,
		fileData: fileData,
	}

	vm.includeNativeFns()
	return &vm
}

func (v *VM) Run() InterpretResult {
	for !v.isAtEnd() && !v.hadError {
		v.oldIp = v.ip
 		i := v.nextByte()

		switch i {
			case compiler.OP_PUSH_CONST:
				v.push(v.currentChunk.Constants[v.getInt()])

			case compiler.OP_PUSH_CLOSURE: {
				// This won't panic, because the compiler only emits this instruction
				// if the constant is a function.
				fn := v.currentChunk.Constants[v.getInt()].(value.ValueFunction)

				upvalues := []*value.Upvalue{}
				upvalueCount := v.getInt()

				// Decode upvalue data from the instruction and put it into the object.
				for range upvalueCount {
					isLocal := v.getByte()
					index := v.getInt()

					if isLocal == 1 {
						// If it's a local, create an upvalue and put it there.
						up := v.captureUpvalue(len(v.callStack) - 1, index)
						upvalues = append(upvalues, up)
					} else {
						// If it's not, get it from the enclosing function's upvalue list.
						upvalues = append(upvalues, v.callStack[len(v.callStack) - 1].function.Upvalues[index])
					}
				}
				
				v.push(value.ValueClosure{
					Fn: &fn,
					Upvalues: upvalues,
				})
			}

			case compiler.OP_APPEND_METHODS: {
				count := v.getInt()
				methodsValues := v.getArguments(count) // actually it's collecting the methods
				methods := make([]value.ValueClosure, 0, len(methodsValues))

				for _, method := range methodsValues {
					methods = append(methods, method.(value.ValueClosure))
				}

				// the record is at the top of the stack now.
				record := v.pop().(value.ValueRecord)
				record.Methods = methods

				v.push(record)
			}

			// TODO: add a separated opcode for concatenating strings when typechecking is added
			case compiler.OP_ADD: {
				if !typesEqual(v.peek(0), v.peek(1)) {
					v.error(
						fmt.Sprintf(
							"Operands types must be equal when adding/concatenating. (left: '%s' (type '%s'), right: '%s' (type '%s'))",
							v.peek(1).String(),
							v.peek(1).Type(),

							v.peek(0).String(),
							v.peek(0).Type(),
						),
					)
					return STATUS_TYPE_ERROR
				}

				// if one operand is string, the other should be too
				if isString(v.peek(0)) {
					status := v.concatenateStrs()

					if status != STATUS_OK {
						return status
					}
				} else if isNumber(v.peek(0)) {
					status := v.binaryNum(i)

					if status != STATUS_OK {
						return status
					}
				} else {
					v.error(
						fmt.Sprintf(
							"Operands must be numbers or strings when adding/concatenating. (left: '%s' (type '%s'), right: '%s' (type '%s'))",
							v.peek(1).String(),
							v.peek(1).Type(),

							v.peek(0).String(),
							v.peek(0).Type(),
						),
					)
					return STATUS_TYPE_ERROR
				}
			}

			case compiler.OP_SUB, compiler.OP_MUL, compiler.OP_DIV, compiler.OP_MODULO: {
				status := v.binaryNum(i)

				if status != STATUS_OK {
					return status
				}
			}

			case compiler.OP_DEF_LOCAL:
				v.callStack[len(v.callStack)-1].locals = append(v.callStack[len(v.callStack)-1].locals, v.pop())

			case compiler.OP_GET_LOCAL:
				v.push(v.callStack[len(v.callStack)-1].locals[v.getInt()])

			case compiler.OP_SET_LOCAL:
				v.callStack[len(v.callStack)-1].locals[v.getInt()] = v.peek(0)

			case compiler.OP_GET_UPVALUE: {
				slot := v.getInt()
				v.push(v.getUpvalueValue(v.callStack[len(v.callStack)-1].function.Upvalues[slot]))
			}

			case compiler.OP_SET_UPVALUE: {
				slot := v.getInt()
				v.setUpvalueValue(v.callStack[len(v.callStack)-1].function.Upvalues[slot], v.peek(0))
			}

			case compiler.OP_DEF_GLOBAL:
				v.globals = append(v.globals, v.pop())

			case compiler.OP_GET_GLOBAL:
				v.push(v.globals[v.getInt()])

			case compiler.OP_SET_GLOBAL:
				v.globals[v.getInt()] = v.peek(0)

			case compiler.OP_GET_PROPERTY: {
				obj := v.pop()
				index := v.getInt()
				
				res := v.getProperty(obj, index)

				if res != STATUS_OK {
					return res
				}
			}

			case compiler.OP_SET_PROPERTY: {
				val := v.pop()
				obj := v.pop()
				index := v.getInt()
				
				res := v.setProperty(obj, index, val)

				if res != STATUS_OK {
					return res
				}
			}

			case compiler.OP_CLOSE_UPVALUE: {
				v.closeUpvalue(len(v.callStack) - 1, len(v.callStack[len(v.callStack)-1].locals) - 1)

				// pop the variable, as it's now safe to pop it,
				// since it's captured and put into the upvalue that captures it.
				v.callStack[len(v.callStack)-1].locals = v.callStack[len(v.callStack)-1].locals[:len(v.callStack[len(v.callStack)-1].locals) - 1]
			}

			case compiler.OP_POP:
				v.pop()

			case compiler.OP_POP_LOCAL:
				v.popVar()

			case compiler.OP_POPN_LOCAL:
				v.popnVar(v.getInt())

			case compiler.OP_JUMP:
				v.ip += v.getInt()

			case compiler.OP_JUMP_TRUE: {
				amount := v.getInt()

				// TODO: check for out of bounds by checking nil
				if b, ok := v.peek(0).(value.ValueBool); ok {
					if v.hadError {
						return STATUS_OUT_OF_BOUNDS
					}

					if b.Value {
						v.ip += amount
					}
				} else {
					v.error(fmt.Sprintf("Given expression ('%s') type is not 'bool'. Its type is '%s'.", v.peek(0).String(),  v.peek(0).Type()))
					return STATUS_TYPE_ERROR
				}
			}

			case compiler.OP_JUMP_FALSE: {
				amount := v.getInt()

				// TODO: check for out of bounds by checking nil
				if b, ok := v.peek(0).(value.ValueBool); ok {
					if v.hadError {
						return STATUS_OUT_OF_BOUNDS
					}

					if !b.Value {
						v.ip += amount
					}
				} else {
					v.error(fmt.Sprintf("Given expression ('%s') type is not 'bool'. Its type is '%s'.", v.peek(0).String(),  v.peek(0).Type()))
					return STATUS_TYPE_ERROR
				}
			}

			case compiler.OP_LOOP:
				v.ip -= v.getInt()

			case compiler.OP_EQUAL: {
				b := v.pop()
				a := v.pop()

				if !typesEqual(a, b) {
					v.error(
						fmt.Sprintf(
							"Types must be the same when comparing. (left: '%s' (type: '%s'), right: '%s', (type: '%s'))",
							a.String(),
							a.Type(),

							b.String(),
							b.Type(),
						),
					)
					return STATUS_TYPE_ERROR
				}

				v.push(value.ValueBool{ Value: valuesEqual(a, b) })
			}

			case compiler.OP_NOT_EQUAL: {
				b := v.pop()
				a := v.pop()

				if !typesEqual(a, b) {
					v.error(
						fmt.Sprintf(
							"Types must be the same when comparing. (left: '%s' (type: '%s'), right: '%s', (type: '%s'))",
							a.String(),
							a.Type(),

							b.String(),
							b.Type(),
						),
					)
					return STATUS_TYPE_ERROR
				}

				v.push(value.ValueBool{ Value: !valuesEqual(a, b) })
			}

			case compiler.OP_GREATER, compiler.OP_GREATER_EQUAL, compiler.OP_LESS, compiler.OP_LESS_EQUAL: {
				status := v.binaryComparison(i)

				if status != STATUS_OK {
					return status
				}
			}

			case compiler.OP_AND, compiler.OP_OR: {
				status := v.binaryBool(i)

				if status != STATUS_OK {
					return status
				}
			}

			case compiler.OP_NOT: {
				op := v.pop()

				if !isBool(op) {
					v.error(fmt.Sprintf("Given expression ('%s') type is not 'bool' to perform a logical not. Its type is '%s'.", op.String(),  op.Type()))
					return STATUS_TYPE_ERROR
				}

				opBool := op.(value.ValueBool)
				v.push(value.ValueBool{ Value: !opBool.Value })
			}

			case compiler.OP_NEGATE: {
				op := v.pop()

				if !isNumber(op) {
					v.error(fmt.Sprintf("Given expression ('%s') type is not 'num' to perform a number negation. Its type is '%s'.", op.String(), op.Type()))
					return STATUS_TYPE_ERROR
				}

				opNum := op.(value.ValueNumber)
				v.push(value.ValueNumber{ Value: -opNum.Value })
			}

            case compiler.OP_MAKE_RANGE: {
                step := v.pop()
                end := v.pop()
                start := v.pop()

                // Check if the given range is valid:

                // 1. Check if all three operands are numbers (or if 'step' is nil).
                // 2. Check if 'end' is reachable.
                // (if 'step' is positive, then 'end' must be greater than 'start', and vice versa. 'step' must never be equal to 0, unless 'start' is equal to 'end')

                if !isNumber(start) {
					v.error(fmt.Sprintf("Given 'start' expression ('%s') type is not 'num'. Its type is '%s'.", start.String(), start.Type()))
					return STATUS_TYPE_ERROR
                } else if !isNumber(end) {
					v.error(fmt.Sprintf("Given 'end' expression ('%s') type is not 'num'. Its type is '%s'.", end.String(), end.Type()))
					return STATUS_TYPE_ERROR
                } else if !isNumber(step) && !isNil(step) {
					v.error(fmt.Sprintf("Given 'step' expression ('%s') type is not 'num' or 'nil'. Its type is '%s'.", step.String(), step.Type()))
					return STATUS_TYPE_ERROR
                }

                startNum := start.(value.ValueNumber).Value
                endNum := end.(value.ValueNumber).Value
                stepNum := 1.0

                // Define 'step' if it hasn't been defined yet. (when its value is 'nil')
                if isNil(step) {
                    if endNum < startNum {
                        stepNum = -1
                    } else if endNum == startNum {
                        stepNum = 0
                    }
                    // if endNum > startNum, then stepNum = 1.
                } else {
                    // 'step' is defined, so we get it.
                    stepNum = step.(value.ValueNumber).Value
                }

                if endNum > startNum && stepNum < 0 ||
                   endNum < startNum && stepNum > 0 ||
                   endNum != startNum && stepNum == 0 {
					v.error("Range's end is unreachable if iterated over.")
					return STATUS_UNREACHABLE_RANGE
                }

                v.push(value.ValueRange{
                    Start: startNum,
                    End: endNum,
                    Step: stepNum,
                })
            }

			case compiler.OP_CALL: {
				arity := v.getInt()
				status := v.call(v.peek(arity), arity)

				if status != STATUS_OK {
					return status
				}
			}

			case compiler.OP_CALL_PROPERTY: {
				index := v.getInt()
				arity := v.getInt()

				obj := v.peek(arity)
				property, res := v.getPropertyValue(obj, index)

				if res != STATUS_OK {
					return res
				}
				
				res = v.call(property, arity)

				if res != STATUS_OK {
					return res
				}
			}

			case compiler.OP_RETURN: {
				v.closeUpvalues(len(v.callStack) - 1)
				frame := v.popFrame()
				var chunk value.Chunk

				if len(v.callStack) == 0 {
					chunk = v.topLevel
				} else {
					chunk = v.callStack[len(v.callStack) - 1].function.Fn.Chunk
				}

				v.ip = frame.oldIp
				v.currentChunk = &chunk
			}

			case compiler.OP_PUSH_TRUE: v.push(value.ValueBool{ Value: true })
			case compiler.OP_PUSH_FALSE: v.push(value.ValueBool{ Value: false })

			case compiler.OP_PUSH_NIL: v.push(value.ValueNil{})
			case compiler.OP_PUSH_VOID: v.push(value.ValueVoid{})

			case compiler.OP_ASSERT_BOOL: {
				if !isBool(v.peek(0)) {
					v.error(fmt.Sprintf("Given expression ('%s') type is not 'bool'. Its type is '%s'.", v.peek(0).String(),  v.peek(0).Type()))
				}
			}

			default:
				panic(fmt.Sprintf("Unknown instruction: '%d'", i))
		}
	}

	return STATUS_OK
}

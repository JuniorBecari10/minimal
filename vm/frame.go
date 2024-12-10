package vm

import "vm-go/value"

type CallFrame struct {
	function *value.ValueFunction
	oldIp int
	variableOffset int
	locals []value.Value
}

package vm

import "vm-go/value"

type CallFrame struct {
	function *value.ValueFunction
	oldIp int
	locals []value.Value
}

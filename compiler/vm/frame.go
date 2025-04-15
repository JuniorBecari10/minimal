package vm

import "vm-go/value"

type CallFrame struct {
	function *value.ValueClosure
	oldIp int
	locals []value.Value
}

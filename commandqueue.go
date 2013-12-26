package cl11

import (
	"strings"
)

type CommandQueueProperties uint8

// Bitfield.
const (
	QueueOutOfOrderExecModeEnable CommandQueueProperties = 1 << iota
	QueueProfilingEnable          CommandQueueProperties = 1 << iota
)

func (properties CommandQueueProperties) String() string {
	var propertiesStrings []string
	if properties&QueueOutOfOrderExecModeEnable != 0 {
		propertiesStrings = append(propertiesStrings, "CL_QUEUE_OUT_OF_ORDER_EXEC_MODE_ENABLE")
	}
	if properties&QueueProfilingEnable != 0 {
		propertiesStrings = append(propertiesStrings, "CL_QUEUE_PROFILING_ENABLE")
	}
	return "(" + strings.Join(propertiesStrings, "|") + ")"
}

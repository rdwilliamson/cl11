package cl11

import (
	"fmt"
	"strings"

	clw "github.com/rdwilliamson/clw11"
)

type CommandQueue struct {
	ID         clw.CommandQueue
	Context    *Context
	Device     *Device
	Properties CommandQueueProperties
}

func (cq CommandQueue) String() string {
	return fmt.Sprintf("%x", cq.ID)
}

type CommandQueueProperties struct {
	OutOfOrderExecution bool
	Profiling           bool
}

func (cqp CommandQueueProperties) String() string {
	var propertiesStrings []string
	if cqp.OutOfOrderExecution {
		propertiesStrings = append(propertiesStrings, "CL_QUEUE_OUT_OF_ORDER_EXEC_MODE_ENABLE")
	}
	if cqp.Profiling {
		propertiesStrings = append(propertiesStrings, "CL_QUEUE_PROFILING_ENABLE")
	}
	return "(" + strings.Join(propertiesStrings, "|") + ")"
}

func CreateCommandQueue(c *Context, d *Device, cqp CommandQueueProperties) (*CommandQueue, error) {
	var properties clw.CommandQueueProperties
	if cqp.OutOfOrderExecution {
		properties |= clw.QueueOutOfOrderExecModeEnable
	}
	if cqp.Profiling {
		properties |= clw.QueueProfilingEnable
	}

	commandQueue, err := clw.CreateCommandQueue(c.ID, d.ID, properties)
	if err != nil {
		return nil, err
	}
	return &CommandQueue{commandQueue, c, d, cqp}, nil
}

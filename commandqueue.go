package cl11

import (
	"strings"

	clw "github.com/rdwilliamson/clw11"
)

// Not thread safe.
type CommandQueue struct {
	ID         clw.CommandQueue
	Context    *Context
	Device     *Device
	Properties CommandQueueProperties

	// TODO some scratch space to avoid allocating memory when converting from
	// []cl.Event to []clw.Event. Should it be thread safe?
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

func (c *Context) CreateCommandQueue(d *Device, cqp CommandQueueProperties) (*CommandQueue, error) {
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
	return &CommandQueue{ID: commandQueue, Context: c, Device: d, Properties: cqp}, nil
}

func (cq *CommandQueue) Flush() error {
	return clw.Flush(cq.ID)
}

func (cq *CommandQueue) Finish() error {
	return clw.Finish(cq.ID)
}

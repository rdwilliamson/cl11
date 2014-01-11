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

func CreateCommandQueue(c *Context, d *Device, p CommandQueueProperties) (*CommandQueue, error) {
	var properties clw.CommandQueueProperties
	if p&QueueOutOfOrderExecModeEnable != 0 {
		properties |= clw.QueueOutOfOrderExecModeEnable
	}
	if p&QueueProfilingEnable != 0 {
		properties |= clw.QueueProfilingEnable
	}

	commandQueue, err := clw.CreateCommandQueue(clw.Context(c.ID), clw.DeviceID(d.ID), properties)
	if err != nil {
		return nil, err
	}
	return &CommandQueue{commandQueue, c, d, p}, nil
}

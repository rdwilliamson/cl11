package cl11

import (
	"strings"

	clw "github.com/rdwilliamson/clw11"
)

// Not thread safe.
type CommandQueue struct {
	id         clw.CommandQueue
	Context    *Context
	Device     *Device
	Properties CommandQueueProperties

	// Scratch space to avoid allocating memory when converting a wait list.
	eventsScratch []clw.Event
}

type CommandQueueProperties uint8

// Bitfield.
const (
	QueueOutOfOrderExecution = CommandQueueProperties(clw.QueueOutOfOrderExecModeEnable)
	QueueProfilingEnable     = CommandQueueProperties(clw.QueueProfilingEnable)
)

func (cqp CommandQueueProperties) String() string {
	var propertiesStrings []string
	if cqp&QueueOutOfOrderExecution != 0 {
		propertiesStrings = append(propertiesStrings, "out of order execution enable")
	}
	if cqp&QueueProfilingEnable != 0 {
		propertiesStrings = append(propertiesStrings, "profiling enable")
	}
	return "(" + strings.Join(propertiesStrings, "|") + ")"
}

func (c *Context) CreateCommandQueue(d *Device, cqp CommandQueueProperties) (*CommandQueue, error) {

	commandQueue, err := clw.CreateCommandQueue(c.id, d.id, clw.CommandQueueProperties(cqp))
	if err != nil {
		return nil, err
	}

	return &CommandQueue{id: commandQueue, Context: c, Device: d, Properties: cqp}, nil
}

func (cq *CommandQueue) Flush() error {
	return clw.Flush(cq.id)
}

func (cq *CommandQueue) Finish() error {
	return clw.Finish(cq.id)
}

func (cq *CommandQueue) toEvents(in []*Event) []clw.Event {

	if in == nil {
		return nil
	}

	if len(cq.eventsScratch) < len(in) {
		cq.eventsScratch = make([]clw.Event, len(in))
	}

	for i := range in {
		cq.eventsScratch[i] = in[i].id
	}

	return cq.eventsScratch[:len(in)]
}

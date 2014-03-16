package cl11

import (
	"strings"

	clw "github.com/rdwilliamson/clw11"
)

// The OpenCL functions that are submitted to a command-queue are enqueued in
// the order the calls are made but can be configured to execute in-order or
// out-of-order. In addition, a wait for events or a barrier command can be
// enqueued to the command-queue. The wait for events command ensures that
// previously enqueued commands identified by the list of events to wait for
// have finished before the next batch of commands is executed. The barrier
// command ensures that all previously enqueued commands in a command-queue have
// finished execution before the next batch of commands is executed.
type CommandQueue struct {
	id clw.CommandQueue

	// The context the command queue was created on.
	Context *Context

	// The device the command queue was created for.
	Device *Device

	// Bit-field list of properties for the command queue.
	Properties CommandQueueProperties

	// Scratch space to avoid allocating memory when converting a wait list.
	eventsScratch []clw.Event
}

type CommandQueueProperties int

// Bitfield.
const (
	QueueOutOfOrderExecution = CommandQueueProperties(clw.QueueOutOfOrderExecModeEnable)
	QueueProfilingEnable     = CommandQueueProperties(clw.QueueProfilingEnable)
)

func (cqp CommandQueueProperties) String() string {
	var propertiesStrings []string
	if cqp&QueueOutOfOrderExecution != 0 {
		propertiesStrings = append(propertiesStrings, "out of order execution")
	}
	if cqp&QueueProfilingEnable != 0 {
		propertiesStrings = append(propertiesStrings, "profiling")
	}
	return "{" + strings.Join(propertiesStrings, ", ") + "}"
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

func (cq *CommandQueue) Retain() error {
	return clw.RetainCommandQueue(cq.id)
}

func (cq *CommandQueue) Release() error {
	return clw.ReleaseCommandQueue(cq.id)
}

func (cq *CommandQueue) EnqueueMarker(e *Event) error {

	if e != nil {
		e.Context = cq.Context
		e.CommandType = CommandMarker
		e.CommandQueue = cq
	}

	return clw.EnqueueMarker(cq.id, &e.id)
}

func (cq *CommandQueue) EnqueueWaitForEvents(waitList []*Event) error {
	return clw.EnqueueWaitForEvents(cq.id, cq.toEvents(waitList))
}

func (cq *CommandQueue) EnqueueBarrier() error {
	return clw.EnqueueBarrier(cq.id)
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

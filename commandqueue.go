package cl11

import (
	"strings"

	clw "github.com/rdwilliamson/clw11"
)

// The OpenCL functions that are submitted to a command-queue are enqueued in
// the order the calls are made but can be configured to execute in-order or
// out-of-order. The properties argument in CreateCommandQueue can be used to
// specify the execution order. In addition, a wait for events or a barrier
// command can be enqueued to the command-queue. The wait for events command
// ensures that previously enqueued commands identified by the list of events to
// wait for have finished before the next batch of commands is executed. The
// barrier command ensures that all previously enqueued commands in a command-
// queue have finished execution before the next batch of commands is executed.
// Similarly, commands to read, write, copy or map memory objects that are
// enqueued after EnqueueNDRangeKernel, EnqueueTask or EnqueueNativeKernel
// commands are not guaranteed to wait for kernels scheduled for execution to
// have completed (if the out of order property is set). To ensure correct
// ordering of commands, the event object returned by EnqueueNDRangeKernel,
// EnqueueTask or EnqueueNativeKernel can be used to enqueue a wait for event or
// a barrier command can be enqueued that must complete before reads or writes
// to the memory object(s) occur.
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

// Create a command-queue on a specific device.
func (c *Context) CreateCommandQueue(d *Device, cqp CommandQueueProperties) (*CommandQueue, error) {

	commandQueue, err := clw.CreateCommandQueue(c.id, d.id, clw.CommandQueueProperties(cqp))
	if err != nil {
		return nil, err
	}

	return &CommandQueue{id: commandQueue, Context: c, Device: d, Properties: cqp}, nil
}

// Issues all previously queued OpenCL commands in a command-queue to the device
// associated with the command-queue.
//
// Flush only guarantees that all queued commands to command_queue get issued to
// the appropriate device. There is no guarantee that they will be complete
// after clFlush returns. Any blocking commands queued in a command-queue
// perform an implicit flush of the command-queue. To use event objects that
// refer to commands enqueued in a command-queue as event objects to wait on by
// commands enqueued in a different command-queue, the application must call a
// Flush or any blocking commands that perform an implicit flush of the command-
// queue where the commands that refer to these event objects are enqueued.
func (cq *CommandQueue) Flush() error {
	return clw.Flush(cq.id)
}

// Blocks until all previously queued OpenCL commands in a command-queue are
// issued to the associated device and have completed.
//
// Blocks until all previously queued OpenCL commands in command_queue are
// issued to the associated device and have completed. Finish does not return
// until all queued commands in command_queue have been processed and completed.
// Finish is also a synchronization point.
func (cq *CommandQueue) Finish() error {
	return clw.Finish(cq.id)
}

// Increments the command_queue reference count.
func (cq *CommandQueue) Retain() error {
	return clw.RetainCommandQueue(cq.id)
}

// Decrements the command_queue reference count.
func (cq *CommandQueue) Release() error {
	return clw.ReleaseCommandQueue(cq.id)
}

// Enqueues a marker command.
//
// Enqueues a marker command to the command queue. The marker command is not
// completed until all commands enqueued before it have completed. The marker
// command returns an event which can be waited on, i.e. this event can be
// waited on to ensure that all commands which have been queued before the
// market command have been completed. complete.
func (cq *CommandQueue) EnqueueMarker(e *Event) error {

	if e != nil {
		e.Context = cq.Context
		e.CommandType = CommandMarker
		e.CommandQueue = cq
	}

	return clw.EnqueueMarker(cq.id, &e.id)
}

// Enqueues a wait for a specific event or a list of events to complete before
// any future commands queued in the command-queue are executed.
//
// The context associated with events in event_list and command_queue must be
// the same.
func (cq *CommandQueue) EnqueueWaitForEvents(waitList []*Event) error {
	return clw.EnqueueWaitForEvents(cq.id, cq.toEvents(waitList))
}

// A synchronization point that enqueues a barrier operation.
//
// EnqueueBarrier is a synchronization point that ensures that all queued
// commands in command_queue have finished execution before the next batch of
// commands can begin execution.
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

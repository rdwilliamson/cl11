package cl11

import (
	"strings"
	"sync"
	"unsafe"

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

	// Pool used when converting a wait list.
	eventPool sync.Pool
}

const eventPoolThreshold = 8

type CommandQueueProperties uint

// Bit field.
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

	return &CommandQueue{id: commandQueue, Context: c, Device: d, Properties: cqp,
		eventPool: sync.Pool{New: func() interface{} { return make([]clw.Event, eventPoolThreshold) }}}, nil
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

// Return the command queue's reference count.
//
// The reference count returned should be considered immediately stale. It is
// unsuitable for general use in applications. This feature is provided for
// identifying memory leaks.
func (cq *CommandQueue) ReferenceCount() (int, error) {

	var count clw.Uint
	err := clw.GetCommandQueueInfo(cq.id, clw.QueueReferenceCount, clw.Size(unsafe.Sizeof(count)),
		unsafe.Pointer(&count), nil)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// Enqueues a marker command.
//
// Enqueues a marker command to the command queue. The marker command is not
// completed until all commands enqueued before it have completed. The marker
// command returns an event which can be waited on.
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

	events := cq.createEvents(waitList)
	err := clw.EnqueueWaitForEvents(cq.id, events)
	cq.releaseEvents(events)
	return err
}

// A synchronization point that enqueues a barrier operation.
//
// EnqueueBarrier is a synchronization point that ensures that all queued
// commands in command_queue have finished execution before the next batch of
// commands can begin execution.
func (cq *CommandQueue) EnqueueBarrier() error {
	return clw.EnqueueBarrier(cq.id)
}

func (cq *CommandQueue) createEvents(waitList []*Event) []clw.Event {
	var result []clw.Event
	if len(waitList) > eventPoolThreshold {
		result = make([]clw.Event, len(waitList))
	}
	result = cq.eventPool.Get().([]clw.Event)[:len(waitList)]
	for i, v := range waitList {
		result[i] = v.id
	}
	return result
}

func (cq *CommandQueue) releaseEvents(events []clw.Event) {
	if cap(events) == eventPoolThreshold {
		cq.eventPool.Put(events)
	}
}

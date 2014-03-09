package cl11

import (
	"fmt"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

// Event objects can be used to refer to a kernel execution command, read,
// write, map and copy commands on memory objects, EnqueueMarker, or user
// events.
//
// An event object can be used to track the execution status of a command. The
// API calls that enqueue commands to a command-queue create a new event object
// that is returned in the event argument.
type Event struct {
	id clw.Event

	// The context the event was created on.
	Context *Context

	// The command type the event represents.
	CommandType CommandType

	// The command queue the event was created on.
	CommandQueue *CommandQueue

	Queued int64 // Profiling information.
	Submit int64 // Profiling information.
	Start  int64 // Profiling information.
	End    int64 // Profiling information.
}

type CommandType int

const (
	CommandNDRangeKernel        = CommandType(clw.CommandNdrangeKernel)
	CommandTask                 = CommandType(clw.CommandTask)
	CommandNativeKernel         = CommandType(clw.CommandNativeKernel)
	CommandReadBuffer           = CommandType(clw.CommandReadBuffer)
	CommandWriteBuffer          = CommandType(clw.CommandWriteBuffer)
	CommandCopyBuffer           = CommandType(clw.CommandCopyBuffer)
	CommandReadImage            = CommandType(clw.CommandReadImage)
	CommandWriteImage           = CommandType(clw.CommandWriteImage)
	CommandCopyImage            = CommandType(clw.CommandCopyImage)
	CommandCopyImageToBuffer    = CommandType(clw.CommandCopyImageToBuffer)
	CommandCopyBufferToImage    = CommandType(clw.CommandCopyBufferToImage)
	CommandMapBuffer            = CommandType(clw.CommandMapBuffer)
	CommandMapImage             = CommandType(clw.CommandMapImage)
	CommandUnmapMemoryObject    = CommandType(clw.CommandUnmapMemoryObject)
	CommandMarker               = CommandType(clw.CommandMarker)
	CommandAcquireGlObjects     = CommandType(clw.CommandAcquireGlObjects)
	CommandReleaseGlObjects     = CommandType(clw.CommandReleaseGlObjects)
	CommandReadBufferRectangle  = CommandType(clw.CommandReadBufferRectangle)
	CommandWriteBufferRectangle = CommandType(clw.CommandWriteBufferRectangle)
	CommandCopyBufferRectangle  = CommandType(clw.CommandCopyBufferRectangle)
	CommandUser                 = CommandType(clw.CommandUser)
)

var commandTypeMap = map[CommandType]string{
	CommandNDRangeKernel:        "ND range kernel",
	CommandTask:                 "task",
	CommandNativeKernel:         "native kernel",
	CommandReadBuffer:           "read buffer",
	CommandWriteBuffer:          "write buffer",
	CommandCopyBuffer:           "copy buffer",
	CommandReadImage:            "read image",
	CommandWriteImage:           "write image",
	CommandCopyImage:            "copy image",
	CommandCopyImageToBuffer:    "copy image to buffer",
	CommandCopyBufferToImage:    "copy buffer to image",
	CommandMapBuffer:            "map buffer",
	CommandMapImage:             "map image",
	CommandUnmapMemoryObject:    "unmap memory object",
	CommandMarker:               "marker",
	CommandAcquireGlObjects:     "acquire GL objects",
	CommandReleaseGlObjects:     "release GL objects",
	CommandReadBufferRectangle:  "read buffer rectangle",
	CommandWriteBufferRectangle: "write buffer rectangle",
	CommandCopyBufferRectangle:  "copy buffer rectangle",
	CommandUser:                 "user",
}

func (ct CommandType) String() string {
	return commandTypeMap[ct]
}

type CommandExecutionStatus int

const (
	Queued    = CommandExecutionStatus(clw.Queued)
	Submitted = CommandExecutionStatus(clw.Submitted)
	Running   = CommandExecutionStatus(clw.Running)
	Complete  = CommandExecutionStatus(clw.Complete)
)

func (ces CommandExecutionStatus) String() string {
	switch ces {
	case Complete:
		return "complete"
	case Running:
		return "running"
	case Submitted:
		return "submitted"
	case Queued:
		return "queued"
	}
	return ""
}

// Return the execution status of the command identified by event.
//
// Statuses are: Queued (command has been enqueued in the command-queue),
// Submitted (enqueued command has been submitted by the host to the device
// associated with the command-queue), Running (device is currently executing
// this command), Complete (the command has completed).
//
// An eventErr is an error that cause the event to abnormally terminate. A
// getStatusErr is an error that occurred while trying to retrieve the event's
// status.
func (e *Event) Status() (ces CommandExecutionStatus, eventErr, getStatusErr error) {

	err := clw.GetEventInfo(e.id, clw.EventCommandExecutionStatus, clw.Size(unsafe.Sizeof(ces)),
		unsafe.Pointer(&ces), nil)
	if err != nil {
		return 0, nil, err
	}

	if ces < 0 {
		return 0, fmt.Errorf("event abnormally terminated: %s", clw.CodeToError(clw.Int(ces)).Error()), nil
	}

	return CommandExecutionStatus(ces), nil, nil
}

// Gets profiling information for the command associated with event if profiling
// is enabled.
//
// The 64-bit values can be used to measure the time in nano-seconds consumed by
// OpenCL commands.
//
// OpenCL devices are required to correctly track time across changes in device
// frequency and power states. The Device.ProfilingTimerResolution specifies the
// resolution of the timer i.e. the number of nanoseconds elapsed before the
// timer is incremented.
//
// Event objects can be used to capture profiling information that measure
// execution time of a command. Profiling of OpenCL commands can be enabled
// either by using a command-queue created with QueueProfilingEnable flag set in
// properties argument to CreateCommandQueue.
func (e *Event) GetProfilingInfo() error {

	var value clw.Ulong
	err := clw.GetEventProfilingInfo(e.id, clw.ProfilingCommandQueued, clw.Size(unsafe.Sizeof(value)),
		unsafe.Pointer(&value), nil)
	if err != nil {
		return err
	}
	e.Queued = int64(value)

	err = clw.GetEventProfilingInfo(e.id, clw.ProfilingCommandSubmit, clw.Size(unsafe.Sizeof(value)),
		unsafe.Pointer(&value), nil)
	if err != nil {
		return err
	}
	e.Submit = int64(value)

	err = clw.GetEventProfilingInfo(e.id, clw.ProfilingCommandStart, clw.Size(unsafe.Sizeof(value)),
		unsafe.Pointer(&value), nil)
	if err != nil {
		return err
	}
	e.Start = int64(value)

	err = clw.GetEventProfilingInfo(e.id, clw.ProfilingCommandEnd, clw.Size(unsafe.Sizeof(value)),
		unsafe.Pointer(&value), nil)
	if err != nil {
		return err
	}
	e.End = int64(value)

	return nil
}

// Registers a user callback function for a specific command execution status.
//
// The registered callback function will be called when the execution status of
// command associated with event changes to Complete.
//
// Each call to SetCallback registers the specified user callback
// function on a callback stack associated with event. The order in which the
// registered user callback functions are called is undefined.
//
// All callbacks registered for an event object must be called. All enqueued
// callbacks shall be called before the event object is destroyed. Callbacks
// must return promptly. The behavior of calling expensive system routines,
// OpenCL API calls to create contexts or command-queues, or blocking OpenCL
// operations, in a callback is undefined.
func (e *Event) SetCallback(callback func(e *Event, userData interface{}), userData interface{}) error {

	return clw.SetEventCallback(e.id, clw.Complete,
		func(event clw.Event, ces clw.CommandExecutionStatus, _userData interface{}) {
			callback(e, _userData)
		},
		userData)
}

// Waits on the host thread for the event to complete. See WaitForEvents for
// more info.
func (e *Event) Wait() error {
	return clw.WaitForEvents([]clw.Event{e.id})
}

// Waits on the host thread for commands identified by event objects to
// complete.
//
// Waits on the host thread for commands identified by event objects in
// event_list to complete. A command is considered complete if its execution
// status is "complete" or abnormally terminated (an error occured).
//
// If the cl_khr_gl_event extension is enabled, event objects can also be used
// to reflect the status of an OpenGL sync object. The sync object in turn
// refers to a fence command executing in an OpenGL command stream. This
// provides another method of coordinating sharing of buffers and images between
// OpenGL and OpenCL.
func WaitForEvents(events ...*Event) error {
	waitList := make([]clw.Event, len(events))
	for i := range events {
		waitList[i] = events[i].id
	}
	return clw.WaitForEvents(waitList)
}

// Creates a user event object.
//
// User events allow applications to enqueue commands that wait on a user event
// to finish before the command is executed by the device. The execution status
// of the user event object created is set to Submitted.
//
// Enqueued commands that specify user events must ensure that the status of the
// user events be set before any OpenCL APIs that release OpenCL objects except
// for event objects are called.
func (c *Context) CreateUserEvent() (*Event, error) {

	event, err := clw.CreateUserEvent(c.id)
	if err != nil {
		return nil, err
	}

	return &Event{id: event, Context: c, CommandType: CommandUser}, nil
}

// Sets the execution status of a user event object to Complete.
func (e *Event) SetComplete() error {
	return clw.SetUserEventStatus(e.id, clw.Int(clw.Complete))
}

// Sets the execution status of a user event object to an error state. All
// enqueued commands that wait on this user event will be terminated. Err must
// be negative.
func (e *Event) SetError(err int) error {
	if err >= 0 {
		return wrapError(fmt.Errorf("can not set event error code to a non-negative value"))
	}
	return clw.SetUserEventStatus(e.id, clw.Int(err))
}

// Increments the event reference count.
//
// The OpenCL commands that return an event perform an implicit retain.
func (e *Event) Retain() error {
	return clw.RetainEvent(e.id)
}

// Decrements the event reference count.
//
// The event object is deleted once the reference count becomes zero, the
// specific command identified by this event has completed (or terminated) and
// there are no commands in the command-queues of a context that require a wait
// for this event to complete.
func (e *Event) Release() error {
	return clw.ReleaseEvent(e.id)
}

// The event reference count.
//
// The reference count returned should be considered immediately stale. It is
// unsuitable for general use in applications. This feature is provided for
// identifying memory leaks.
func (e *Event) ReferenceCount() (int, error) {
	var param clw.Uint
	err := clw.GetEventInfo(e.id, clw.EventReferenceCount, clw.Size(unsafe.Sizeof(param)), unsafe.Pointer(&param), nil)
	return int(param), err
}

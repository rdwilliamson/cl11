package cl11

import (
	"fmt"
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

//
type Event struct {
	id           clw.Event
	Context      *Context
	CommandType  CommandType
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
	Complete  = CommandExecutionStatus(clw.Complete)
	Running   = CommandExecutionStatus(clw.Running)
	Submitted = CommandExecutionStatus(clw.Submitted)
	Queued    = CommandExecutionStatus(clw.Queued)
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

// Creates a user event object.
//
// User events allow applications to enqueue commands that wait on a user event
// to finish before the command is executed by the device. The execution status
// of the user event object created is set to Submitted.
func (c *Context) CreateUserEvent() (*Event, error) {

	event, err := clw.CreateUserEvent(c.id)
	if err != nil {
		return nil, err
	}

	return &Event{id: event, Context: c, CommandType: CommandUser}, nil
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

func (e *Event) SetCallback(callback func(e *Event, userData interface{}), userData interface{}) error {

	return clw.SetEventCallback(e.id, clw.Complete,
		func(event clw.Event, ces clw.CommandExecutionStatus, _userData interface{}) {
			callback(e, _userData)
		},
		userData)
}

func (e *Event) Wait() error {
	return clw.WaitForEvents([]clw.Event{e.id})
}

func WaitForEvents(events ...*Event) error {
	waitList := make([]clw.Event, len(events))
	for i := range events {
		waitList[i] = events[i].id
	}
	return clw.WaitForEvents(waitList)
}

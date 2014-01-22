package cl11

import (
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

type Event clw.Event

func (c *Context) CreateUserEvent() (Event, error) {
	event, err := clw.CreateUserEvent(c.ID)
	if err != nil {
		return nil, err
	}
	return Event(event), nil
}

type CommandType int

const (
	CommandNDRangeKernel        CommandType = iota
	CommandTask                 CommandType = iota
	CommandNativeKernel         CommandType = iota
	CommandReadBuffer           CommandType = iota
	CommandWriteBuffer          CommandType = iota
	CommandCopyBuffer           CommandType = iota
	CommandReadImage            CommandType = iota
	CommandWriteImage           CommandType = iota
	CommandCopyImage            CommandType = iota
	CommandCopyImageToBuffer    CommandType = iota
	CommandCopyBufferToImage    CommandType = iota
	CommandMapBuffer            CommandType = iota
	CommandMapImage             CommandType = iota
	CommandUnmapMemoryObject    CommandType = iota
	CommandMarker               CommandType = iota
	CommandAcquireGlObjects     CommandType = iota
	CommandReleaseGlObjects     CommandType = iota
	CommandReadBufferRectangle  CommandType = iota
	CommandWriteBufferRectangle CommandType = iota
	CommandCopyBufferRectangle  CommandType = iota
	CommandUser                 CommandType = iota
)

func (ct CommandType) String() string {
	switch ct {
	case CommandNDRangeKernel:
		return "ND range kernel"
	case CommandTask:
		return "task"
	case CommandNativeKernel:
		return "native kernel"
	case CommandReadBuffer:
		return "read buffer"
	case CommandWriteBuffer:
		return "write buffer"
	case CommandCopyBuffer:
		return "copy buffer"
	case CommandReadImage:
		return "read image"
	case CommandWriteImage:
		return "write image"
	case CommandCopyImage:
		return "copy image"
	case CommandCopyImageToBuffer:
		return "copy image to buffer"
	case CommandCopyBufferToImage:
		return "copy buffer to image"
	case CommandMapBuffer:
		return "map buffer"
	case CommandMapImage:
		return "map image"
	case CommandUnmapMemoryObject:
		return "unmap memory object"
	case CommandMarker:
		return "marker"
	case CommandAcquireGlObjects:
		return "acquire GL objects"
	case CommandReleaseGlObjects:
		return "release GL objects"
	case CommandReadBufferRectangle:
		return "read buffer rectangle"
	case CommandWriteBufferRectangle:
		return "write buffer rectangle"
	case CommandCopyBufferRectangle:
		return "copy buffer rectangle"
	case CommandUser:
		return "user"
	}
	panic("unknown command type")
}

type CommandExecutionStatus int

const (
	Complete  CommandExecutionStatus = iota
	Running   CommandExecutionStatus = iota
	Submitted CommandExecutionStatus = iota
	Queued    CommandExecutionStatus = iota
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
	panic("unknown command execution status")
}

func toEvents(in []Event) []clw.Event {
	return *(*[]clw.Event)(unsafe.Pointer(&in))
}

// Returns the events status, an error that caused the event to terminate, or an
// error that occurred trying to retrieve the event status.
func EventStatus(e Event) (CommandExecutionStatus, error, error) {
	// Not a method because the underlying C type is a pointer and throws a
	// "invalid receiver type *Event (Event is a pointer type)" error.

	var status clw.CommandExecutionStatus
	err := clw.GetEventInfo(clw.Event(e), clw.EventCommandExecutionStatus, clw.Size(unsafe.Sizeof(status)),
		unsafe.Pointer(&status), nil)
	if err != nil {
		return 0, nil, err
	}

	if status < clw.CommandExecutionStatus(0) {
		return 0, clw.CodeToError(clw.Int(status)), nil
	}

	switch status {
	case clw.Complete:
		return Complete, nil, nil
	case clw.Running:
		return Running, nil, nil
	case clw.Submitted:
		return Submitted, nil, nil
	case clw.Queued:
		return Queued, nil, nil
	}

	panic("unknown command execution status")
}

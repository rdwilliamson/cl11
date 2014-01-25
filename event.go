package cl11

import (
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

type Event struct {
	ID          clw.Event
	CommandType CommandType
}

func (c *Context) CreateUserEvent() (*Event, error) {
	event, err := clw.CreateUserEvent(c.ID)
	if err != nil {
		return nil, err
	}
	return &Event{ID: event}, nil
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

type CommandExecutionStatus int8

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
	panic("unknown command execution status")
}

func toEvents(in []*Event) []clw.Event {

	if in == nil {
		return nil
	}

	// TODO avoid allocating memory.
	out := make([]clw.Event, len(in))
	for i := range in {
		out[i] = in[i].ID
	}
	return out
}

// Returns the events status, an error that caused the event to terminate, or an
// error that occurred trying to retrieve the event status.
func (e *Event) Status() (CommandExecutionStatus, error, error) {
	var status clw.CommandExecutionStatus
	err := clw.GetEventInfo(e.ID, clw.EventCommandExecutionStatus, clw.Size(unsafe.Sizeof(status)),
		unsafe.Pointer(&status), nil)
	if err != nil {
		return 0, nil, err
	}

	if status < 0 {
		return 0, clw.CodeToError(clw.Int(status)), nil
	}

	return CommandExecutionStatus(status), nil, nil
}

package cl11

import (
	"unsafe"

	clw "github.com/rdwilliamson/clw11"
)

type Event struct {
	id           clw.Event
	Context      *Context
	CommandType  CommandType
	CommandQueue *CommandQueue
}

// Device counter times in nanoseconds.
type EventProfilingInfo struct {
	Queued int64
	Submit int64
	Start  int64
	End    int64
}

func (c *Context) CreateUserEvent() (*Event, error) {

	event, err := clw.CreateUserEvent(c.id)
	if err != nil {
		return nil, err
	}

	return &Event{id: event, Context: c, CommandType: CommandUser}, nil
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

// Returns the events status, an error that caused the event to terminate, or an
// error that occurred trying to retrieve the event status.
func (e *Event) Status() (CommandExecutionStatus, error, error) {

	var status clw.CommandExecutionStatus
	err := clw.GetEventInfo(e.id, clw.EventCommandExecutionStatus, clw.Size(unsafe.Sizeof(status)),
		unsafe.Pointer(&status), nil)
	if err != nil {
		return 0, nil, err
	}

	if status < 0 {
		return 0, clw.CodeToError(clw.Int(status)), nil
	}

	return CommandExecutionStatus(status), nil, nil
}

func (e *Event) ProfilingInfo() (*EventProfilingInfo, error) {

	var queued clw.Ulong
	err := clw.GetEventProfilingInfo(e.id, clw.ProfilingCommandQueued, clw.Size(unsafe.Sizeof(queued)),
		unsafe.Pointer(&queued), nil)
	if err != nil {
		return nil, err
	}

	var submit clw.Ulong
	err = clw.GetEventProfilingInfo(e.id, clw.ProfilingCommandSubmit, clw.Size(unsafe.Sizeof(submit)),
		unsafe.Pointer(&submit), nil)
	if err != nil {
		return nil, err
	}

	var start clw.Ulong
	err = clw.GetEventProfilingInfo(e.id, clw.ProfilingCommandStart, clw.Size(unsafe.Sizeof(start)),
		unsafe.Pointer(&start), nil)
	if err != nil {
		return nil, err
	}

	var end clw.Ulong
	err = clw.GetEventProfilingInfo(e.id, clw.ProfilingCommandEnd, clw.Size(unsafe.Sizeof(end)),
		unsafe.Pointer(&end), nil)
	if err != nil {
		return nil, err
	}

	return &EventProfilingInfo{int64(queued), int64(submit), int64(start), int64(end)}, nil
}

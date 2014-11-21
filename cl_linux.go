package cl11

import clw "github.com/rdwilliamson/clw11"

const (
	GLContext  = ContextProperties(clw.GLContext)
	EGLDisplay = ContextProperties(clw.EGLDisplay)
	GLXDisplay = ContextProperties(clw.GLXDisplay)
)

var InvalidGLSharegroupReference = clw.InvalidGLSharegroupReference

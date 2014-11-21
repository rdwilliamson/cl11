package cl11

import clw "github.com/rdwilliamson/clw11"

const (
	GLContext  = ContextProperties(clw.GLContext)
	EGLDisplay = ContextProperties(clw.EGLDisplay)
	WGLHDC     = ContextProperties(clw.WGLHDC)
)

var InvalidGLSharegroupReference = clw.InvalidGLSharegroupReference

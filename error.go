package cl11

import (
	"fmt"
	"runtime"
	"strings"
)

// Error with calling function.
type Error struct {
	Function string
	Err      error
}

func (err Error) Error() string {
	return fmt.Sprint(err.Function, ": ", err.Err.Error())
}

// Gets "package.function" from call stack for error.
func wrapError(err error) error {
	pc, _, _, _ := runtime.Caller(2)
	name := runtime.FuncForPC(pc).Name()
	last := strings.LastIndex(name, "/")
	if last == -1 {
		last = 0
	} else {
		last++
	}
	return &Error{name[last:], err}
}

package report

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type ErrorManager struct {
	errors []*DiagnosticItem
}

// todo: possibly we need to add mutex here

func NewErrorManager() *ErrorManager {
	return &ErrorManager{
		errors: make([]*DiagnosticItem, 0),
	}
}

func (e *ErrorManager) AddError(item *DiagnosticItem) {
	file, line := getCallerInfo(2) // Adjust skip to 2 to get caller of AddError.

	fmt.Printf("%s:%d %s line: %d\n", file, line, item.Msg, item.Pointer.Loc.Line)

	e.errors = append(e.errors, item)
}

func (e *ErrorManager) Errors() []*DiagnosticItem {
	return e.errors
}

// getCallerInfo retrieves the caller's file, line number, and function name.
func getCallerInfo(skip int) (file string, line int) {
	_, fullPath, line, ok := runtime.Caller(skip)
	if !ok {
		fullPath = "unknown"
		line = 0
	}
	file = filepath.Base(fullPath)
	return
}

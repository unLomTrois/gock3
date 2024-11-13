package report

import "log"

type ErrorManager struct {
	errors []*DiagnosticItem
}

func NewErrorManager() *ErrorManager {
	return &ErrorManager{
		errors: make([]*DiagnosticItem, 0),
	}
}

func (e *ErrorManager) AddError(item *DiagnosticItem) {
	log.Println(item.Msg, "line:", item.Pointer.Loc.Line)
	e.errors = append(e.errors, item)
}

func (e *ErrorManager) Errors() []*DiagnosticItem {
	return e.errors
}

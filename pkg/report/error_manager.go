package report

type ErrorManager struct {
	errors []*DiagnosticItem
}

func NewErrorManager() *ErrorManager {
	return &ErrorManager{
		errors: make([]*DiagnosticItem, 0),
	}
}

func (e *ErrorManager) AddError(item *DiagnosticItem) {
	e.errors = append(e.errors, item)
}

func (e *ErrorManager) Errors() []*DiagnosticItem {
	return e.errors
}

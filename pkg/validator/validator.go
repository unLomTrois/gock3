package validator

import "github.com/unLomTrois/gock3/pkg/report"

type Validator interface {
	Errors() []*report.DiagnosticItem
}

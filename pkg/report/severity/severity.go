package severity

type Severity int

const (
	Info Severity = iota
	Warning
	Error
	Critical
)

func (s Severity) String() string {
	switch s {
	case Info:
		return "Info"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	case Critical:
		return "Critical"
	default:
		return "Unknown"
	}
}

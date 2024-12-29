package data

type Entity interface {
	Name() string
	Location() string
	GetKind() string
}

package entity

type Entity interface {
	Name() string
	Location() string
	GetKind() EntityKind
}

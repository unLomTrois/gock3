package symboltable

import "github.com/unLomTrois/gock3/pkg/data"

type SymbolTableInterface interface {
	AddEntity(entity data.Entity)
	AddEntities(entities []data.Entity)
	Get(name string) data.Entity
	Contains(name string) bool
}

type SymbolTable struct {
	table map[string]data.Entity
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		table: make(map[string]data.Entity),
	}
}

func (s *SymbolTable) AddEntity(entity data.Entity) {
	name := entity.Name()

	s.table[name] = entity
}

func (s *SymbolTable) AddEntities(entities []data.Entity) {
	for _, entity := range entities {
		s.AddEntity(entity)
	}
}

func (s *SymbolTable) Get(name string) data.Entity {
	return s.table[name]
}

func (s *SymbolTable) Contains(name string) bool {
	_, ok := s.table[name]

	return ok
}

func (s *SymbolTable) Len() int {
	return len(s.table)
}

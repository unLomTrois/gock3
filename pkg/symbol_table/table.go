package symboltable

import "github.com/unLomTrois/gock3/pkg/data"

type SymbolTableInterface interface {
	AddEntity(entity data.Entity)
	AddEntities(entities []data.Entity)
	Get(name string) data.Entity
	Contains(name string) bool
}

type EntityKind = string

type SymbolTable struct {
	store map[EntityKind]map[string]data.Entity
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: make(map[EntityKind]map[string]data.Entity),
	}
}

func (st *SymbolTable) AddEntity(entity data.Entity) {
	kind := entity.GetKind()

	if _, exists := st.store[kind]; !exists {
		st.store[kind] = make(map[string]data.Entity)
	}

	name := entity.Name()
	st.store[kind][name] = entity
}

func (st *SymbolTable) AddEntities(entities []data.Entity) {
	for _, entity := range entities {
		st.AddEntity(entity)
	}
}

func (st *SymbolTable) Get(kind EntityKind, name string) (data.Entity, bool) {
	if entities, ok := st.store[kind]; ok {
		e, found := entities[name]
		return e, found
	}
	return nil, false
}

func (st *SymbolTable) Contains(kind EntityKind, name string) bool {
	if entities, ok := st.store[kind]; ok {
		_, found := entities[name]
		return found
	}
	return false
}

func (s *SymbolTable) Len() int {
	// iterate
	var count int
	for _, entities := range s.store {
		count += len(entities)
	}
	return count
}

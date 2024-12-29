package symboltable

import (
	"github.com/unLomTrois/gock3/pkg/entity"
)

type SymbolTableInterface interface {
	AddEntity(entity entity.Entity)
	AddEntities(entities []entity.Entity)
	Get(name string) entity.Entity
	Contains(name string) bool
}
type SymbolTable struct {
	store map[entity.EntityKind]map[string]entity.Entity
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: make(map[entity.EntityKind]map[string]entity.Entity),
	}
}

func (st *SymbolTable) AddEntity(item entity.Entity) {
	kind := item.GetKind()

	if _, exists := st.store[kind]; !exists {
		st.store[kind] = make(map[string]entity.Entity)
	}

	name := item.Name()
	st.store[kind][name] = item
}

func (st *SymbolTable) AddEntities(entities []entity.Entity) {
	for _, entity := range entities {
		st.AddEntity(entity)
	}
}

func (st *SymbolTable) Get(kind entity.EntityKind, name string) (entity.Entity, bool) {
	if entities, ok := st.store[kind]; ok {
		e, found := entities[name]
		return e, found
	}
	return nil, false
}

func (st *SymbolTable) Contains(kind entity.EntityKind, name string) bool {
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

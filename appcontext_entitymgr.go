package core

import "fmt"

func (appCtx *AppContextBehavior) GetEntity(id uint64) (Entity, bool) {
	entity, ok := appCtx.entityMap.Load(id)
	if !ok {
		return nil, false
	}
	return entity.(Entity), true
}

func (appCtx *AppContextBehavior) RangeEntities(fun func(entity Entity) bool) {
	if fun == nil {
		return
	}

	appCtx.entityMap.Range(func(key, value interface{}) bool {
		return fun(value.(Entity))
	})
}

func (appCtx *AppContextBehavior) AddEntity(entity Entity) {
	if entity == nil {
		panic("nil entity")
	}

	if entity.GetID() <= 0 {
		panic("entity id equal 0 invalid")
	}

	if _, loaded := appCtx.entityMap.LoadOrStore(entity.GetID(), entity); loaded {
		panic(fmt.Errorf("repeated entity '{%d}' in this app context", entity.GetID()))
	}
}

func (appCtx *AppContextBehavior) RemoveEntity(id uint64) {
	appCtx.entityMap.Delete(id)
}

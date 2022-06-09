package core

type EntityQuery interface {
	GetEntity(id uint64) (Entity, bool)
	RangeEntities(func(entity Entity) bool)
}

type EntityMgr interface {
	EntityQuery
	AddEntity(entity Entity)
	RemoveEntity(id uint64)
}

type EntityMgrEvents interface {
	EventEntityMgrAddEntity() IEvent
	EventEntityMgrRemoveEntity() IEvent
	EventEntityMgrEntityAddComponents() IEvent
	EventEntityMgrEntityRemoveComponent() IEvent
}

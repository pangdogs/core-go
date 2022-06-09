//go:generate go run github.com/pangdogs/core/eventcode -decl $GOFILE -package $GOPACKAGE -core ""
package core

type EventEntityMgrAddEntity[T any] interface {
	OnEntityMgrAddEntity(entityMgr T, entity Entity)
}

type EventEntityMgrRemoveEntity[T any] interface {
	OnEntityMgrRemoveEntity(entityMgr T, entity Entity)
}

type EventEntityMgrEntityAddComponents[T any] interface {
	OnEntityMgrEntityAddComponents(entityMgr T, entity Entity, components []Component)
}

type EventEntityMgrEntityRemoveComponent[T any] interface {
	OnEntityMgrEntityRemoveComponent(entityMgr T, entity Entity, component Component)
}

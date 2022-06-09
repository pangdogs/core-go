//go:generate go run github.com/pangdogs/core/eventcode -decl $GOFILE -package $GOPACKAGE -core "" -exportemit=false
package core

type EventCompMgrAddComponents[T any] interface {
	OnCompMgrAddComponents(compMgr T, components []Component)
}

type EventCompMgrRemoveComponent[T any] interface {
	OnCompMgrRemoveComponent(compMgr T, component Component)
}

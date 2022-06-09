//go:generate go run github.com/pangdogs/core/eventcode -decl $GOFILE -package $GOPACKAGE -core ""
package core

type EventEntityDestroySelf interface {
	OnEntityDestroySelf(entity Entity)
}

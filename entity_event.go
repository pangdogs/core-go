//go:generate go run github.com/pangdogs/core/eventcode -decl $GOFILE -core "" -emit_package $GOPACKAGE -export_emit=false
package core

type EventEntityDestroySelf interface {
	OnEntityDestroySelf(entity Entity)
}

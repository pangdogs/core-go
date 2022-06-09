//go:generate go run github.com/pangdogs/core/eventcode -decl $GOFILE -package $GOPACKAGE -core "" -exportemit=false
package core

type EventComponentDestroySelf interface {
	OnComponentDestroySelf(comp Component)
}

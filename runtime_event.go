//go:generate go run github.com/pangdogs/core/eventcode -decl $GOFILE -package $GOPACKAGE -core ""
package core

type EventUpdate interface {
	Update()
}

type EventLateUpdate interface {
	LateUpdate()
}

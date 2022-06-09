//go:generate go run github.com/pangdogs/core/eventcode -decl $GOFILE -package $GOPACKAGE -core ""
package core

type EventPushSafeCallSegment interface {
	OnPushSafeCallSegment(segment func())
}

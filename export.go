package core

import (
	"github.com/pangdogs/core/internal"
	"github.com/pangdogs/core/internal/foundation"
)

type Context = internal.Context

var NewContext = foundation.NewContext

type App = internal.App

var NewApp = foundation.NewApp

var NewAppOption = foundation.NewAppOption

type AppInheritorFoundation = foundation.AppInheritor

type Runtime = internal.Runtime

var NewRuntime = foundation.NewRuntime

var NewRuntimeOption = foundation.NewRuntimeOption

type Frame = internal.Frame

var NewFrame = foundation.NewFrame

type Entity = internal.Entity

var NewEntity = foundation.NewEntity

var NewEntityOption = foundation.NewEntityOption

type EntityInheritorFoundation = foundation.EntityInheritor

type Component = internal.Component

type ComponentFoundation = foundation.Component

type Hook = internal.Hook

type HookFoundation = foundation.Hook

type EventSource = internal.EventSource

type EventSourceFoundation = foundation.EventSource

var BindEvent = foundation.BindEvent

var UnbindEvent = foundation.UnbindEvent

var UnbindAllEventSource = foundation.UnbindAllEventSource

var UnbindAllHook = foundation.UnbindAllHook

var SendEvent = foundation.SendEvent

type SafeRet = internal.SafeRet

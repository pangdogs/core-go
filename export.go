package core

import (
	"github.com/pangdogs/core/internal/foundation"
	"github.com/pangdogs/core/internal/list"
)

type Context = foundation.Context

var NewContext = foundation.NewContext

type App = foundation.App

var NewApp = foundation.NewApp

var NewAppOption = foundation.NewAppOption

type AppFoundation = foundation.AppFoundation

type Runtime = foundation.Runtime

var NewRuntime = foundation.NewRuntime

var NewRuntimeOption = foundation.NewRuntimeOption

type RuntimeFoundation = foundation.RuntimeFoundation

type GC = foundation.GC

type Frame = foundation.Frame

var NewFrame = foundation.NewFrame

type Entity = foundation.Entity

var NewEntity = foundation.NewEntity

var NewEntityOption = foundation.NewEntityOption

type EntityFoundation = foundation.EntityFoundation

type Component = foundation.Component

type ComponentFoundation = foundation.ComponentFoundation

type ComponentInit = foundation.ComponentInit

type ComponentAwake = foundation.ComponentAwake

type ComponentEntityInit = foundation.ComponentEntityInit

type ComponentStart = foundation.ComponentStart

type ComponentUpdate = foundation.ComponentUpdate

type ComponentLateUpdate = foundation.ComponentLateUpdate

type ComponentEntityShut = foundation.ComponentEntityShut

type ComponentHalt = foundation.ComponentHalt

type ComponentShut = foundation.ComponentShut

type Hook = foundation.Hook

type HookFoundation = foundation.HookFoundation

type EventSource = foundation.EventSource

type EventSourceFoundation = foundation.EventSourceFoundation

var BindEvent = foundation.BindEvent

var UnbindEvent = foundation.UnbindEvent

var UnbindAllEventSource = foundation.UnbindAllEventSource

var UnbindAllHook = foundation.UnbindAllHook

var SendEvent = foundation.SendEvent

type SafeStack = foundation.SafeStack

var NewSafeStack = foundation.NewSafeStack

type SafeRet = foundation.SafeRet

var UnsafeCall = foundation.UnsafeCall

type List = list.List

type Element = list.Element

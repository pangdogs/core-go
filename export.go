package core

import (
	"github.com/pangdogs/core/internal/foundation"
	"github.com/pangdogs/core/internal/misc"
)

type Context = foundation.Context

var NewContext = foundation.NewContext

type App = foundation.App

var AppGetInheritor = foundation.AppGetInheritor

var NewApp = foundation.NewApp

var NewAppOption = foundation.NewAppOption

type AppFoundation = foundation.AppFoundation

type Runtime = foundation.Runtime

var RuntimeGetInheritor = foundation.RuntimeGetInheritor

var NewRuntime = foundation.NewRuntime

var NewRuntimeOption = foundation.NewRuntimeOption

type RuntimeFoundation = foundation.RuntimeFoundation

type GC = foundation.GC

type Frame = foundation.Frame

var NewFrame = foundation.NewFrame

type Entity = foundation.Entity

var EntityGetInheritor = foundation.EntityGetInheritor

var EntitySetLifecycleEntityInitFunc = foundation.EntitySetLifecycleEntityInitFunc

var EntitySetLifecycleStartFunc = foundation.EntitySetLifecycleStartFunc

var EntitySetLifecycleUpdateFunc = foundation.EntitySetLifecycleUpdateFunc

var EntitySetLifecycleLateUpdateFunc = foundation.EntitySetLifecycleLateUpdateFunc

var EntitySetLifecycleEntityShutFunc = foundation.EntitySetLifecycleEntityShutFunc

var EntityGetLifecycleEntityInitFunc = foundation.EntityGetLifecycleEntityInitFunc

var EntityGetLifecycleStartFunc = foundation.EntityGetLifecycleStartFunc

var EntityGetLifecycleUpdateFunc = foundation.EntityGetLifecycleUpdateFunc

var EntityGetLifecycleLateUpdateFunc = foundation.EntityGetLifecycleLateUpdateFunc

var EntityGetLifecycleEntityShutFunc = foundation.EntityGetLifecycleEntityShutFunc

var NewEntity = foundation.NewEntity

var NewEntityOption = foundation.NewEntityOption

type EntityFoundation = foundation.EntityFoundation

type Component = foundation.Component

var ComponentGetInheritor = foundation.ComponentGetInheritor

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

var InitHook = foundation.InitHook

type EventSource = foundation.EventSource

type EventSourceFoundation = foundation.EventSourceFoundation

var InitEventSource = foundation.InitEventSource

var AllocEventID = foundation.AllocEventID

var BindEvent = foundation.BindEvent

var UnbindEvent = foundation.UnbindEvent

var UnbindAllEventSource = foundation.UnbindAllEventSource

var UnbindAllHook = foundation.UnbindAllHook

var SendEvent = foundation.SendEvent

type EventRet = foundation.EventRet

const (
	EventRet_Continue    = foundation.EventRet_Continue
	EventRet_Break       = foundation.EventRet_Break
	EventRet_Unsubscribe = foundation.EventRet_Unsubscribe
)

type SafeStack = foundation.SafeStack

var NewSafeStack = foundation.NewSafeStack

type SafeRet = foundation.SafeRet

var UnsafeCall = foundation.UnsafeCall

type List = misc.List

type Element = misc.Element

type IFace = misc.IFace

var NilIFace = misc.NilIFace

type Cache = misc.Cache

var NewList = misc.NewList

var NewCache = misc.NewCache

type AspectJointPointTab = foundation.AspectJointPointTab

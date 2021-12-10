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

var ComponentSetNotAnalysisLifecycle = foundation.ComponentSetNotAnalysisLifecycle

var ComponentGetNotAnalysisLifecycle = foundation.ComponentGetNotAnalysisLifecycle

var ComponentSetLifecycleComponentInit = foundation.ComponentSetLifecycleComponentInit

var ComponentSetLifecycleComponentAwake = foundation.ComponentSetLifecycleComponentAwake

var ComponentSetLifecycleComponentEntityInit = foundation.ComponentSetLifecycleComponentEntityInit

var ComponentSetLifecycleComponentStart = foundation.ComponentSetLifecycleComponentStart

var ComponentSetLifecycleComponentUpdate = foundation.ComponentSetLifecycleComponentUpdate

var ComponentSetLifecycleComponentLateUpdate = foundation.ComponentSetLifecycleComponentLateUpdate

var ComponentSetLifecycleComponentEntityShut = foundation.ComponentSetLifecycleComponentEntityShut

var ComponentSetLifecycleComponentHalt = foundation.ComponentSetLifecycleComponentHalt

var ComponentSetLifecycleComponentShut = foundation.ComponentSetLifecycleComponentShut

var ComponentGetLifecycleComponentInit = foundation.ComponentGetLifecycleComponentInit

var ComponentGetLifecycleComponentAwake = foundation.ComponentGetLifecycleComponentAwake

var ComponentGetLifecycleComponentEntityInit = foundation.ComponentGetLifecycleComponentEntityInit

var ComponentGetLifecycleComponentStart = foundation.ComponentGetLifecycleComponentStart

var ComponentGetLifecycleComponentUpdate = foundation.ComponentGetLifecycleComponentUpdate

var ComponentGetLifecycleComponentLateUpdate = foundation.ComponentGetLifecycleComponentLateUpdate

var ComponentGetLifecycleComponentEntityShut = foundation.ComponentGetLifecycleComponentEntityShut

var ComponentGetLifecycleComponentHalt = foundation.ComponentGetLifecycleComponentHalt

var ComponentGetLifecycleComponentShut = foundation.ComponentGetLifecycleComponentShut

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

var AllocEventID = foundation.AllocEventID

type EventSource = foundation.EventSource

type EventSourceFoundation = foundation.EventSourceFoundation

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

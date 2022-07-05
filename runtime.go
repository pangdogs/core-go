package core

import (
	"github.com/pangdogs/core/container"
)

type Runtime interface {
	container.GC
	Runnable
	init(runtimeCtx RuntimeContext, opts *RuntimeOptions)
	getOptions() *RuntimeOptions
	GetID() uint64
	GetRuntimeCtx() RuntimeContext
}

func RuntimeGetOptions(runtime Runtime) RuntimeOptions {
	return *runtime.getOptions()
}

func RuntimeGetInheritor(runtime Runtime) Runtime {
	return runtime.getOptions().Inheritor
}

func NewRuntime(runtimeCtx RuntimeContext, optFuncs ...NewRuntimeOptionFunc) Runtime {
	opts := &RuntimeOptions{}
	NewRuntimeOption.Default()(opts)

	for _, optFun := range optFuncs {
		optFun(opts)
	}

	var runtime *RuntimeBehavior

	if opts.Inheritor != nil {
		opts.Inheritor.init(runtimeCtx, opts)
		return opts.Inheritor
	}

	runtime = &RuntimeBehavior{}
	runtime.init(runtimeCtx, opts)

	return runtime.opts.Inheritor
}

type RuntimeBehavior struct {
	id              uint64
	opts            RuntimeOptions
	ctx             RuntimeContext
	hooksMap        map[uint64][3]Hook
	processQueue    chan func()
	eventUpdate     Event
	eventLateUpdate Event
}

func (runtime *RuntimeBehavior) GC() bool {
	runtime.ctx.GC()
	runtime.eventUpdate.GC()
	runtime.eventLateUpdate.GC()
	return true
}

func (runtime *RuntimeBehavior) MarkGC() {
}

func (runtime *RuntimeBehavior) NeedGC() bool {
	return true
}

func (runtime *RuntimeBehavior) CollectGC(gc container.GC) {
}

func (runtime *RuntimeBehavior) init(runtimeCtx RuntimeContext, opts *RuntimeOptions) {
	if runtimeCtx == nil {
		panic("nil runtimeCtx")
	}

	if opts == nil {
		panic("nil opts")
	}

	runtime.opts = *opts

	if runtime.opts.Inheritor == nil {
		runtime.opts.Inheritor = runtime
	}

	runtime.id = runtimeCtx.GetAppCtx().genUID()
	runtime.ctx = runtimeCtx
	runtime.hooksMap = make(map[uint64][3]Hook)

	runtime.eventUpdate.Init(runtime.getOptions().EnableAutoRecover, runtimeCtx.GetReportError(), EventRecursion_Disallow, runtimeCtx.getOptions().HookCache, runtime)
	runtime.eventLateUpdate.Init(runtime.getOptions().EnableAutoRecover, runtimeCtx.GetReportError(), EventRecursion_Disallow, runtimeCtx.getOptions().HookCache, runtime)

	if opts.EnableAutoRun {
		runtime.opts.Inheritor.Run()
	}
}

func (runtime *RuntimeBehavior) getOptions() *RuntimeOptions {
	return &runtime.opts
}

func (runtime *RuntimeBehavior) GetID() uint64 {
	return runtime.id
}

func (runtime *RuntimeBehavior) GetRuntimeCtx() RuntimeContext {
	return runtime.ctx
}

func (runtime *RuntimeBehavior) OnEntityMgrAddEntity(runtimeCtx RuntimeContext, entity Entity) {
	if entityInit, ok := entity.(EntityInit); ok {
		entityInit.Init()
	}

	entity.RangeComponents(func(comp Component) bool {
		comp.setPrimary(true)
		return true
	})

	entity.RangeComponents(func(comp Component) bool {
		if !comp.getPrimary() {
			return true
		}

		if compAwake, ok := comp.(ComponentAwake); ok {
			compAwake.Awake()
		}

		return true
	})

	entity.RangeComponents(func(comp Component) bool {
		if !comp.getPrimary() {
			return true
		}

		if compStart, ok := comp.(ComponentStart); ok {
			compStart.Start()
		}

		return true
	})

	if entityInitFin, ok := entity.(EntityInitFin); ok {
		entityInitFin.InitFin()
	}

	runtime.connectEntity(entity)
}

func (runtime *RuntimeBehavior) OnEntityMgrRemoveEntity(runtimeCtx RuntimeContext, entity Entity) {
	runtime.disconnectEntity(entity)

	if entityShut, ok := entity.(EntityShut); ok {
		entityShut.Shut()
	}

	entity.RangeComponents(func(comp Component) bool {
		if compShut, ok := comp.(ComponentShut); ok {
			compShut.Shut()
		}
		return true
	})

	if entityShutFin, ok := entity.(EntityShutFin); ok {
		entityShutFin.ShutFin()
	}
}

func (runtime *RuntimeBehavior) OnEntityMgrEntityAddComponents(runtimeCtx RuntimeContext, entity Entity, components []Component) {
	for _, comp := range components {
		if compAwake, ok := comp.(ComponentAwake); ok {
			compAwake.Awake()
		}
	}

	for _, comp := range components {
		if compStart, ok := comp.(ComponentStart); ok {
			compStart.Start()
		}
	}

	for _, comp := range components {
		runtime.connectComponent(comp)
	}
}

func (runtime *RuntimeBehavior) OnEntityMgrEntityRemoveComponent(runtimeCtx RuntimeContext, entity Entity, component Component) {
	runtime.disconnectComponent(component)

	if compShut, ok := component.(ComponentShut); ok {
		compShut.Shut()
	}
}

func (runtime *RuntimeBehavior) OnEntityDestroySelf(entity Entity) {
	runtime.ctx.RemoveEntity(entity.GetID())
}

func (runtime *RuntimeBehavior) OnComponentDestroySelf(comp Component) {
	comp.GetEntity().RemoveComponentByID(comp.GetID())
}

func (runtime *RuntimeBehavior) connectEntity(entity Entity) {
	var hooks [3]Hook

	if entityUpdate, ok := entity.(EntityUpdate); ok {
		hooks[0] = BindEvent[EntityUpdate](&runtime.eventUpdate, entityUpdate)
	}

	if entityLateUpdate, ok := entity.(EntityLateUpdate); ok {
		hooks[1] = BindEvent[EntityLateUpdate](&runtime.eventLateUpdate, entityLateUpdate)
	}

	hooks[2] = BindEvent[EventEntityDestroySelf](entity.EventEntityDestroySelf(), runtime)

	runtime.hooksMap[entity.GetID()] = hooks

	entity.RangeComponents(func(comp Component) bool {
		runtime.connectComponent(comp)
		return true
	})
}

func (runtime *RuntimeBehavior) connectComponent(comp Component) {
	var hooks [3]Hook

	if compUpdate, ok := comp.(ComponentUpdate); ok {
		hooks[0] = BindEvent[ComponentUpdate](&runtime.eventUpdate, compUpdate)
	}

	if compLateUpdate, ok := comp.(ComponentLateUpdate); ok {
		hooks[1] = BindEvent[ComponentLateUpdate](&runtime.eventLateUpdate, compLateUpdate)
	}

	hooks[2] = BindEvent[EventComponentDestroySelf](comp.EventComponentDestroySelf(), runtime)

	runtime.hooksMap[comp.GetID()] = hooks
}

func (runtime *RuntimeBehavior) disconnectEntity(entity Entity) {
	entity.RangeComponents(func(comp Component) bool {
		runtime.disconnectComponent(comp)
		return true
	})

	hooks, ok := runtime.hooksMap[entity.GetID()]
	if ok {
		delete(runtime.hooksMap, entity.GetID())

		for _, hook := range hooks {
			hook.Unbind()
		}
	}
}

func (runtime *RuntimeBehavior) disconnectComponent(comp Component) {
	hooks, ok := runtime.hooksMap[comp.GetID()]
	if ok {
		delete(runtime.hooksMap, comp.GetID())

		for _, hook := range hooks {
			hook.Unbind()
		}
	}
}

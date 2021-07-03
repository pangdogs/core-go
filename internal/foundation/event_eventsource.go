package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal"
	"github.com/pangdogs/core/internal/list"
)

type EventSourceWhole interface {
	internal.EventSource
	internal.GC
	getRuntime() internal.Runtime
	addHook(hook internal.Hook, priority ...int) error
	removeHook(hookID uint64)
	rangeHooks(fun func(hook internal.Hook, priority int) bool)
}

func NewHookBundle(hook internal.Hook, priority ...int) (*HookBundle, error) {
	if hook == nil {
		return nil, errors.New("nil hook")
	}

	return &HookBundle{
		Hook: hook,
		Priority: func() int {
			if len(priority) > 0 {
				return priority[0]
			}
			return 0
		}(),
	}, nil
}

type HookBundle struct {
	Hook     internal.Hook
	Priority int
}

type EventSource struct {
	id         uint64
	runtime    internal.Runtime
	hookList   list.List
	hookMap    map[uint64]*list.Element
	hookGCList []*list.Element
}

func (es *EventSource) InitEventSource(rt internal.Runtime) {
	if rt == nil {
		panic("nil runtime")
	}

	es.id = rt.GetApp().(AppWhole).MakeUID()
	es.runtime = rt
	es.hookList.Init()
	es.hookMap = map[uint64]*list.Element{}
}

func (es *EventSource) GetEventSourceID() uint64 {
	return es.id
}

func (es *EventSource) getRuntime() internal.Runtime {
	return es.runtime
}

func (es *EventSource) addHook(hook internal.Hook, priority ...int) error {
	hb, err := NewHookBundle(hook, priority...)
	if err != nil {
		return err
	}

	if _, ok := es.hookMap[hook.GetHookID()]; ok {
		return errors.New("hook id already exists")
	}

	for e := es.hookList.Front(); e != nil; e = e.Next() {
		if hb.Priority < e.Value.(*HookBundle).Priority {
			es.hookMap[hook.GetHookID()] = es.hookList.InsertBefore(hb, e)
			return nil
		}
	}

	es.hookMap[hook.GetHookID()] = es.hookList.PushBack(hb)

	return nil
}

func (es *EventSource) removeHook(hookID uint64) {
	if e, ok := es.hookMap[hookID]; ok {
		delete(es.hookMap, hookID)
		e.SetMark(0, true)
		es.hookGCList = append(es.hookGCList, e)
		es.runtime.(RuntimeWhole).PushGC(es)
	}
}

func (es *EventSource) rangeHooks(fun func(hook internal.Hook, priority int) bool) {
	if fun == nil {
		return
	}

	es.hookList.UnsafeTraversal(func(e *list.Element) bool {
		if e.Escape() || e.GetMark(0) {
			return true
		}
		return fun(e.Value.(*HookBundle).Hook, e.Value.(*HookBundle).Priority)
	})
}

func (es *EventSource) GC() {
	for i := 0; i < len(es.hookGCList); i++ {
		es.hookList.Remove(es.hookGCList[i])
	}
	es.hookGCList = es.hookGCList[:0]
}

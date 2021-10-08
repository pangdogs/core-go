package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal/list"
)

type EventSource interface {
	GC
	InitEventSource(rt Runtime)
	GetEventSourceID() uint64
	getRuntime() Runtime
	addHook(hook Hook, priority ...int) error
	removeHook(hookID uint64)
	rangeHooks(fun func(hook Hook, priority int) bool)
}

func NewHookBundle(hook Hook, _priority ...int) (*HookBundle, error) {
	if hook == nil {
		return nil, errors.New("nil hook")
	}

	priority := 0
	if len(_priority) > 0 {
		priority = _priority[0]
	}

	return &HookBundle{
		Hook:     hook,
		Priority: priority,
	}, nil
}

type HookBundle struct {
	Hook     Hook
	Priority int
}

type EventSourceFoundation struct {
	id         uint64
	runtime    Runtime
	hookList   list.List
	hookMap    map[uint64]*list.Element
	hookGCList []*list.Element
}

func (es *EventSourceFoundation) GC() {
	for i := 0; i < len(es.hookGCList); i++ {
		es.hookList.Remove(es.hookGCList[i])
	}
	es.hookGCList = es.hookGCList[:0]
}

func (es *EventSourceFoundation) InitEventSource(rt Runtime) {
	if rt == nil {
		panic("nil runtime")
	}

	es.id = rt.GetApp().makeUID()
	es.runtime = rt
	es.hookList.Init()
	es.hookMap = map[uint64]*list.Element{}
}

func (es *EventSourceFoundation) GetEventSourceID() uint64 {
	return es.id
}

func (es *EventSourceFoundation) getRuntime() Runtime {
	return es.runtime
}

func (es *EventSourceFoundation) addHook(hook Hook, priority ...int) error {
	if hook == nil {
		return errors.New("nil hook")
	}

	if _, ok := es.hookMap[hook.GetHookID()]; ok {
		return errors.New("hook id already exists")
	}

	hb, err := NewHookBundle(hook, priority...)
	if err != nil {
		return err
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

func (es *EventSourceFoundation) removeHook(hookID uint64) {
	if e, ok := es.hookMap[hookID]; ok {
		delete(es.hookMap, hookID)
		e.SetMark(0, true)
		es.hookGCList = append(es.hookGCList, e)
		es.runtime.PushGC(es)
	}
}

func (es *EventSourceFoundation) rangeHooks(fun func(hook Hook, priority int) bool) {
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

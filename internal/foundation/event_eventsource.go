package foundation

import (
	"errors"
	"github.com/pangdogs/core/internal/list"
)

type EventSource interface {
	InitEventSource(rt Runtime)
	GetEventSourceID() uint64
	GetEventSourceRuntime() Runtime
	addHook(hook Hook, priority int32) error
	removeHook(hookID uint64)
	rangeHooks(fun func(hook interface{}, priority int32) bool)
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

func (es *EventSourceFoundation) GetEventSourceRuntime() Runtime {
	return es.runtime
}

func (es *EventSourceFoundation) addHook(hook Hook, priority int32) error {
	if hook == nil {
		return errors.New("nil hook")
	}

	if _, ok := es.hookMap[hook.GetHookID()]; ok {
		return errors.New("hook id already exists")
	}

	for e := es.hookList.Front(); e != nil; e = e.Next() {
		if priority < int32(e.Mark>>32) {
			ne := es.hookList.InsertBefore(hook, e)
			ne.Mark |= uint64(priority) << 32
			es.hookMap[hook.GetHookID()] = ne
			return nil
		}
	}

	ne := es.hookList.PushBack(hook)
	ne.Mark |= uint64(priority) << 32
	es.hookMap[hook.GetHookID()] = ne

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

func (es *EventSourceFoundation) rangeHooks(fun func(hook interface{}, priority int32) bool) {
	if fun == nil {
		return
	}

	es.hookList.UnsafeTraversal(func(e *list.Element) bool {
		if e.Escape() || e.GetMark(0) {
			return true
		}
		return fun(e.Value, int32(e.Mark>>32))
	})
}

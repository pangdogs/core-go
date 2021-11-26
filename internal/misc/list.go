// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package misc implements a doubly linked misc.
//
// To iterate over a misc (where l is a *List):
//	for e := l.Front(); e != nil; e = e.Next() {
//		// do something with e.Value
//	}
//
package misc

import "unsafe"

type IFace [2]unsafe.Pointer

// Element is an element of a linked misc.
type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element

	// The list to which this element belongs.
	list *List

	// The value stored with this element.
	Value [4]IFace

	// Mark 标记
	Mark [4]uint64
}

// Next returns the next misc element or nil.
func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous misc element or nil.
func (e *Element) Prev() *Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Escape 是否已从链表脱出
func (e *Element) Escape() bool {
	return e.list == nil
}

// SetMark 设置标记
func (e *Element) SetMark(bit int, v bool) {
	if v {
		e.Mark[bit/64] |= 1 << bit
	} else {
		e.Mark[bit/64] &= ^(1 << bit)
	}
}

// GetMark 获取标记
func (e *Element) GetMark(bit int) bool {
	return (e.Mark[bit/64]>>bit)&uint64(1) == 1
}

// SetValue 设置数据
func (e *Element) SetValue(index int, v interface{}) {
	e.Value[index] = *(*IFace)(unsafe.Pointer(&v))
}

// GetValue 获取数据
func (e *Element) GetValue(index int) interface{} {
	return *(*interface{})(unsafe.Pointer(&e.Value[index]))
}

// SetIFace 设置接口指针，用于提高接口转换效率
func (e *Element) SetIFace(index int, f IFace) {
	e.Value[index] = f
}

// GetIFace 获取接口指针，用于提高接口转换效率
func (e *Element) GetIFace(index int) unsafe.Pointer {
	return unsafe.Pointer(&e.Value[index])
}

// List represents a doubly linked misc.
// The zero value for List is an empty misc ready to use.
type List struct {
	cache *Cache  // 元素分配缓存，用于减轻GC压力
	root  Element // sentinel misc element, only &root, root.prev, and root.next are used
	len   int     // current misc length excluding (this) sentinel element
}

// Init initializes or clears misc l.
func (l *List) Init(cache *Cache) *List {
	l.cache = cache
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// NewList returns an initialized misc.
func NewList(cache *Cache) *List { return new(List).Init(cache) }

// Len returns the number of elements of misc l.
// The complexity is O(1).
func (l *List) Len() int { return l.len }

// Front returns the first element of misc l or nil if the misc is empty.
func (l *List) Front() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of misc l or nil if the misc is empty.
func (l *List) Back() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero List value.
func (l *List) lazyInit() {
	if l.root.next == nil {
		l.Init(nil)
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *List) insert(e, at *Element) *Element {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *List) insertValue(v interface{}, at *Element) *Element {
	e := l.cache.Alloc()
	e.SetValue(0, v)
	return l.insert(e, at)
}

// insertIFace is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *List) insertIFace(f IFace, at *Element) *Element {
	e := l.cache.Alloc()
	e.SetIFace(0, f)
	return l.insert(e, at)
}

// remove removes e from its misc, decrements l.len, and returns e.
func (l *List) remove(e *Element) *Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}

// move moves e to next to at and returns e.
func (l *List) move(e, at *Element) *Element {
	if e == at {
		return e
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e

	return e
}

// Remove removes e from l if e is an element of misc l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *List) Remove(e *Element) interface{} {
	if e.list == l {
		// if e.misc == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of misc l and returns e.
func (l *List) PushFront(v interface{}) *Element {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element e with value v at the back of misc l and returns e.
func (l *List) PushBack(v interface{}) *Element {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before Mark and returns e.
// If Mark is not an element of l, the misc is not modified.
// The Mark must not be nil.
func (l *List) InsertBefore(v interface{}, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after Mark and returns e.
// If Mark is not an element of l, the misc is not modified.
// The Mark must not be nil.
func (l *List) InsertAfter(v interface{}, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// PushIFaceFront inserts a new element e with value v at the front of misc l and returns e.
func (l *List) PushIFaceFront(f IFace) *Element {
	l.lazyInit()
	return l.insertIFace(f, &l.root)
}

// PushIFaceBack inserts a new element e with value v at the back of misc l and returns e.
func (l *List) PushIFaceBack(f IFace) *Element {
	l.lazyInit()
	return l.insertIFace(f, l.root.prev)
}

// InsertIFaceBefore inserts a new element e with value v immediately before Mark and returns e.
// If Mark is not an element of l, the misc is not modified.
// The Mark must not be nil.
func (l *List) InsertIFaceBefore(f IFace, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertIFace(f, mark.prev)
}

// InsertIFaceAfter inserts a new element e with value v immediately after Mark and returns e.
// If Mark is not an element of l, the misc is not modified.
// The Mark must not be nil.
func (l *List) InsertIFaceAfter(f IFace, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertIFace(f, mark)
}

// MoveToFront moves element e to the front of misc l.
// If e is not an element of l, the misc is not modified.
// The element must not be nil.
func (l *List) MoveToFront(e *Element) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, &l.root)
}

// MoveToBack moves element e to the back of misc l.
// If e is not an element of l, the misc is not modified.
// The element must not be nil.
func (l *List) MoveToBack(e *Element) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, l.root.prev)
}

// MoveBefore moves element e to its new position before Mark.
// If e or Mark is not an element of l, or e == Mark, the misc is not modified.
// The element and Mark must not be nil.
func (l *List) MoveBefore(e, mark *Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

// MoveAfter moves element e to its new position after Mark.
// If e or Mark is not an element of l, or e == Mark, the misc is not modified.
// The element and Mark must not be nil.
func (l *List) MoveAfter(e, mark *Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

// PushBackList inserts a copy of another misc at the back of misc l.
// The lists l and other may be the same. They must not be nil.
func (l *List) PushBackList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of another misc at the front of misc l.
// The lists l and other may be the same. They must not be nil.
func (l *List) PushFrontList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}

// SafeTraversal 安全遍历元素，性能较差，中途可以删除元素
func (l *List) SafeTraversal(visitor func(e *Element) bool) {
	if visitor == nil {
		return
	}

	snap := make([]*Element, 0, l.Len())

	for e := l.Front(); e != nil; e = e.Next() {
		snap = append(snap, e)
	}

	for i := 0; i < len(snap); i++ {
		if !visitor(snap[i]) {
			break
		}
	}
}

// UnsafeTraversal 不安全遍历元素，性能较好，中途不能删除元素
func (l *List) UnsafeTraversal(visitor func(e *Element) bool) {
	if visitor == nil {
		return
	}

	for e := l.Front(); e != nil; e = e.Next() {
		if !visitor(e) {
			break
		}
	}
}

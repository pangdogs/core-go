package container

// Element 元素
type Element[T any] struct {
	_next, _prev *Element[T]
	list         *List[T]
	escaped      bool
	Value        T
	GC           GC
}

// next 下一个元素，包含正在删除的元素
func (e *Element[T]) next() *Element[T] {
	if n := e._next; e.list != nil && n != &e.list.root {
		return n
	}
	return nil
}

// prev 前一个元素，包含正在删除的元素
func (e *Element[T]) prev() *Element[T] {
	if p := e._prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Next 下一个元素
func (e *Element[T]) Next() *Element[T] {
	for n := e.next(); n != nil; n = n.next() {
		if !n.escaped {
			return n
		}
	}
	return nil
}

// Prev 前一个元素
func (e *Element[T]) Prev() *Element[T] {
	for p := e.prev(); p != nil; p = p.prev() {
		if !p.escaped {
			return p
		}
	}
	return nil
}

// Escape 从链表中删除
func (e *Element[T]) Escape() {
	if e.list != nil {
		e.escaped = true
	}
}

// Escaped 是否已从链表中删除
func (e *Element[T]) Escaped() bool {
	return e.escaped
}

// NewList 创建链表
func NewList[T any](cache *Cache[T]) *List[T] {
	return new(List[T]).Init(cache)
}

// List 链表
type List[T any] struct {
	cache *Cache[T]
	root  Element[T]
	len   int
}

// Init 初始化
func (l *List[T]) Init(cache *Cache[T]) *List[T] {
	l.cache = cache
	l.root._next = &l.root
	l.root._prev = &l.root
	l.len = 0
	return l
}

// GC 执行GC
func (l *List[T]) GC() {
	for e := l.Front(); e != nil; e = e.next() {
		if e.escaped {
			l.remove(e)
		} else {
			if e.GC != nil {
				e.GC()
			}
		}
	}
}

// Len 链表长度
func (l *List[T]) Len() int {
	return l.len
}

// Front 链表头部
func (l *List[T]) Front() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root._next
}

// Back 链表尾部
func (l *List[T]) Back() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root._prev
}

// lazyInit 迟滞初始化
func (l *List[T]) lazyInit() {
	if l.root._next == nil {
		l.Init(nil)
	}
}

// insert 插入元素
func (l *List[T]) insert(e, at *Element[T]) *Element[T] {
	e._prev = at
	e._next = at._next
	e._prev._next = e
	e._next._prev = e
	e.list = l
	l.len++
	return e
}

// insertValue 插入数据
func (l *List[T]) insertValue(value T, at *Element[T]) *Element[T] {
	e := l.cache.alloc()
	e.Value = value
	return l.insert(e, at)
}

// remove 删除元素
func (l *List[T]) remove(e *Element[T]) *Element[T] {
	e._prev._next = e._next
	e._next._prev = e._prev
	e._next = nil
	e._prev = nil
	e.list = nil
	l.len--
	return e
}

// move 移动元素
func (l *List[T]) move(e, at *Element[T]) *Element[T] {
	if e == at {
		return e
	}
	e._prev._next = e._next
	e._next._prev = e._prev

	e._prev = at
	e._next = at._next
	e._prev._next = e
	e._next._prev = e

	return e
}

// Remove 删除元素
func (l *List[T]) Remove(e *Element[T]) T {
	if e.list == l {
		e.Escape()
	}
	return e.Value
}

// PushFront 在链表头部插入数据
func (l *List[T]) PushFront(value T) *Element[T] {
	l.lazyInit()
	return l.insertValue(value, &l.root)
}

// PushBack 在链表尾部插入数据
func (l *List[T]) PushBack(value T) *Element[T] {
	l.lazyInit()
	return l.insertValue(value, l.root._prev)
}

// InsertBefore 在链表指定位置前插入数据
func (l *List[T]) InsertBefore(value T, mark *Element[T]) *Element[T] {
	if mark.list != l {
		return nil
	}
	return l.insertValue(value, mark._prev)
}

// InsertAfter 在链表指定位置后插入数据
func (l *List[T]) InsertAfter(value T, mark *Element[T]) *Element[T] {
	if mark.list != l {
		return nil
	}
	return l.insertValue(value, mark)
}

// MoveToFront 移动元素至链表头部
func (l *List[T]) MoveToFront(e *Element[T]) {
	if e.list != l || l.root._next == e {
		return
	}
	l.move(e, &l.root)
}

// MoveToBack 移动元素至链表尾部
func (l *List[T]) MoveToBack(e *Element[T]) {
	if e.list != l || l.root._prev == e {
		return
	}
	l.move(e, l.root._prev)
}

// MoveBefore 移动元素至链表指定位置前
func (l *List[T]) MoveBefore(e, mark *Element[T]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark._prev)
}

// MoveAfter 移动元素至链表指定位置后
func (l *List[T]) MoveAfter(e, mark *Element[T]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

// PushBackList 在链表尾部插入其他链表
func (l *List[T]) PushBackList(other *List[T]) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.next() {
		l.insertValue(e.Value, l.root._prev)
	}
}

// PushFrontList 在链表头部插入其他链表
func (l *List[T]) PushFrontList(other *List[T]) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.prev() {
		l.insertValue(e.Value, &l.root)
	}
}

// Traversal 遍历元素
func (l *List[T]) Traversal(visitor func(e *Element[T]) bool) {
	if visitor == nil {
		return
	}

	for e := l.Front(); e != nil; e = e.next() {
		if !e.escaped && !visitor(e) {
			break
		}
	}
}

// TraversalAt 从指定位置开始遍历元素
func (l *List[T]) TraversalAt(visitor func(e *Element[T]) bool, mark *Element[T]) {
	if visitor == nil || mark.list != l {
		return
	}

	for e := mark; e != nil; e = e.next() {
		if !e.escaped && !visitor(e) {
			break
		}
	}
}

// ReverseTraversal 反向遍历元素
func (l *List[T]) ReverseTraversal(visitor func(e *Element[T]) bool) {
	if visitor == nil {
		return
	}

	for e := l.Back(); e != nil; e = e.prev() {
		if !e.escaped && !visitor(e) {
			break
		}
	}
}

// ReverseTraversalAt 从指定位置开始反向遍历元素
func (l *List[T]) ReverseTraversalAt(visitor func(e *Element[T]) bool, mark *Element[T]) {
	if visitor == nil || mark.list != l {
		return
	}

	for e := mark; e != nil; e = e.prev() {
		if !e.escaped && !visitor(e) {
			break
		}
	}
}

package list

// Element is an element of a linked list.
type Element[T any] struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element[T]

	// The list to which this element belongs.
	list *LinkedList[T]

	// The value stored with this element.
	value T
}

// Next returns the next list element or nil.
func (e *Element[T]) Next() *Element[T] {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Element[T]) Prev() *Element[T] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Value returns the value stored with this element.
func (e *Element[T]) Value() T {
	return e.value
}

// LinkedList represents a doubly linked list.
// The zero value for LinkedList is an empty list ready to use.
type LinkedList[T any] struct {
	root Element[T] // sentinel list element, only &root, root.prev, and root.next are used
	len  int        // current list length excluding sentinel
}

// Init initializes or clears list l.
func (l *LinkedList[T]) Init() *LinkedList[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// New returns an initialized list.
func New[T any]() *LinkedList[T] {
	return new(LinkedList[T]).Init()
}

// Len returns the number of elements of list l.
func (l *LinkedList[T]) Len() int { return l.len }

// Front returns the first element of list l or nil if the list is empty.
func (l *LinkedList[T]) Front() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty.
func (l *LinkedList[T]) Back() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero LinkedList value.
func (l *LinkedList[T]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *LinkedList[T]) insert(e, at *Element[T]) *Element[T] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element[T]{value: v}, at).
func (l *LinkedList[T]) insertValue(v T, at *Element[T]) *Element[T] {
	return l.insert(&Element[T]{value: v}, at)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *LinkedList[T]) PushBack(v T) *Element[T] {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *LinkedList[T]) PushFront(v T) *Element[T] {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

package linkedlist

type List[K comparable, V any] struct {
	head *Node[K, V]
	tail *Node[K, V]
	size int
}

func (l *List[K, V]) Size() int {
	return l.size
}

func (l *List[K, V]) PushNode(node *Node[K, V]) *Node[K, V] {
	node.prev = l.tail.prev
	node.next = l.tail
	l.tail.prev.next = node
	l.tail.prev = node
	l.size++
	return node
}

func (l *List[K, V]) DeleteNode(node *Node[K, V]) {
	node.prev.next = node.next
	node.next.prev = node.prev
	l.size--

	node.prev = nil
	node.next = nil
}

func (l *List[K, V]) DeleteHead() *Node[K, V] {
	node := l.head.next
	l.DeleteNode(node)
	return node
}

func New[K comparable, V any]() *List[K, V] {
	head := &Node[K, V]{}
	tail := &Node[K, V]{}
	head.next = tail
	tail.prev = head
	return &List[K, V]{
		head: head,
		tail: tail,
	}
}

package linkedlist

func NewNode[K comparable, V any](k K, v V) *Node[K, V] {
	return &Node[K, V]{
		key: k,
		val: v,
	}
}

type Node[K comparable, V any] struct {
	prev, next *Node[K, V]
	key        K
	val        V
}

func (n Node[K, V]) GetValue() V {
	return n.val
}

func (n Node[K, V]) SetValue(v V) {
	n.val = v
}

func (n Node[K, V]) GetKey() K {
	return n.key
}

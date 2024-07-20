package imdb

type IMDB[K comparable, V any] interface {
	Get(K) (V, bool)
	Put(K, V)
	Remove(K) bool
}

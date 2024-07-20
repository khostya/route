package output

type Message[T any] struct {
	t       string
	message T
}

func (m Message[T]) GetMessage() T {
	return m.message
}

func BuildMessageChan[T any](module string, messages <-chan T) <-chan Message[T] {
	out := make(chan Message[T])
	go func() {
		defer close(out)
		for m := range messages {
			out <- Message[T]{message: m, t: module}
		}
	}()

	return out
}

func FilterMessageChan[T any](module string, messages <-chan Message[T]) <-chan Message[T] {
	out := make(chan Message[T])
	go func() {
		defer close(out)
		for m := range messages {
			if m.t != module {
				continue
			}
			out <- m
		}
	}()

	return out
}

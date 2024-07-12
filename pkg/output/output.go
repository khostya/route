package output

type Controller[T any] struct {
	subscribers []chan<- T
}

func NewController[T any]() *Controller[T] {
	controller := &Controller[T]{}

	return controller
}

func (c *Controller[T]) Add(messages <-chan T) {
	go c.add(messages)
}

func (c *Controller[T]) add(messages <-chan T) {
	for m := range messages {
		c.Push(m)
	}
}

func (c *Controller[T]) Push(message T) {
	go c.push(message)
}

func (c *Controller[T]) push(message T) {
	for _, subscriber := range c.subscribers {
		subscriber <- message
	}
}

func (c *Controller[T]) Subscribe() chan T {
	out := make(chan T, 10)
	c.subscribers = append(c.subscribers, out)
	return out
}

func (c *Controller[T]) Close() {
	for _, subscriber := range c.subscribers {
		close(subscriber)
	}
}

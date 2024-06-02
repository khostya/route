package cli

type handler struct {
	name   string
	handle commandHandler
}

func newHandler(name string, handle func([]string) string) handler {
	return handler{name: name, handle: handle}
}

func newHandlers(service orderService) []handler {
	executor := Executor{service: service}

	return []handler{
		newHandler(refundOrder, executor.refundOrder),
		newHandler(issueOrder, executor.issueOrders),
		newHandler(returnOrder, executor.returnOrder),
		newHandler(deliverOrder, executor.deliverOrder),
		newHandler(listOrder, executor.listOrder),
		newHandler(listRefunded, executor.listRefunded),
	}
}

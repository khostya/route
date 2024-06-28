package cli

import "slices"

type (
	handler struct {
		name   string
		handle commandHandler
	}
	handlers []handler
)

func newHandler(name string, handle commandHandler) handler {
	return handler{name: name, handle: handle}
}

func newHandlers(service orderService) handlers {
	executor := newExecutor(service)

	return []handler{
		newHandler(refundOrder, executor.refundOrder),
		newHandler(issueOrders, executor.issueOrders),
		newHandler(returnOrder, executor.returnOrder),
		newHandler(deliverOrder, executor.deliverOrder),
		newHandler(listOrders, executor.listOrders),
		newHandler(listRefunded, executor.listRefunded),
	}
}

func (h handlers) mustFind(name string) handler {
	var handlers []handler = h

	return handlers[slices.IndexFunc(handlers, func(h handler) bool {
		return h.name == name
	})]
}

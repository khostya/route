package cli

import "fmt"

var (
	deliverOrderUsage = fmt.Sprintf("%s %s %s %s", deliverOrder, orderIdParam, userIdParam, expParam)
	returnOrderUsage  = fmt.Sprintf("%s %s", returnOrder, orderIdParam)
	issueOrdersUsage  = fmt.Sprintf("%s %s", issueOrders, ordersIdsParam)
	listOrdersUsage   = fmt.Sprintf("%s %s %s", listOrders, userIdParam, sizeParam)
	refundOrderUsage  = fmt.Sprintf("%s %s %s", refundOrder, orderIdParam, userIdParam)
	listRefundedUsage = fmt.Sprintf("%s %s %s", listRefunded, sizeParam, pageParam)
	procsUsage        = fmt.Sprintf("%s %s", procs, nParam)
)

const (
	orderIdParam   = "--id=<id заказа>"
	ordersIdsParam = "<id заказа 1> ... <id заказа N>"
	userIdParam    = "--user=<id>"
	expParam       = "--exp=2024-06-09T17:12:32+05:00"
	sizeParam      = "--size=20"
	pageParam      = "--page=10"
	nParam         = "--n=10"

	helpDescription = "Cправка"

	deliverOrderDescription = `На вход принимается ID заказа, ID получателя и срок хранения. Заказ нельзя принять дважды.`

	returnOrderDescription = `На вход принимается ID заказа. Метод должен удалять заказ из вашего файла. Можно вернуть только те заказы, у которых вышел срок хранения и если заказы находятся в пвз.`

	issueOrdersDescription = `Можно выдавать только те заказы, которые были приняты от курьера и чей срок хранения меньше текущей даты. Все ID заказов должны принадлежать только одному клиенту.`

	listOrdersDescription = `На вход принимается ID пользователя как обязательный параметр и опциональные параметры. Параметры позволяют получать только последние N заказов или заказы клиента, находящиеся в нашем ПВЗ.`

	refundOrderDescription = `На вход принимается ID пользователя и ID заказа. Заказ может быть возвращен в течение двух дней с момента выдачи.`

	listRefundedDescription = `Получить список всех заказов, которые вернули клиенты: Метод должен выдавать список пагинированно.`

	procsDescription = "Изменить максимальное количество горутин"

	exitDescription = `Завершить выполнение`
)

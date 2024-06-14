package cli

import (
	"fmt"
	"time"
)

var (
	deliverOrderUsage = fmt.Sprintf("%s %s %s %s", deliverOrder, orderIdParamUsage, userIdParamUsage, expParamUsage)
	returnOrderUsage  = fmt.Sprintf("%s %s", returnOrder, orderIdParamUsage)
	issueOrdersUsage  = fmt.Sprintf("%s %s", issueOrders, ordersIdsParamUsage)
	listOrdersUsage   = fmt.Sprintf("%s %s %s", listOrders, userIdParamUsage, sizeParamUsage)
	refundOrderUsage  = fmt.Sprintf("%s %s %s", refundOrder, orderIdParamUsage, userIdParamUsage)
	listRefundedUsage = fmt.Sprintf("%s %s %s", listRefunded, sizeParamUsage, pageParamUsage)
	workersUsage      = fmt.Sprintf("%s %s", workers, nParamUsage)

	orderIdParamUsage   = fmt.Sprintf("--%s=1", orderIdParam)
	userIdParamUsage    = fmt.Sprintf("--%s=1", userIdParam)
	expParamUsage       = fmt.Sprintf("--%s=%s", expParam, time.Now().Add(time.Hour*2).Format(time.RFC3339))
	sizeParamUsage      = fmt.Sprintf("--%s=20", sizeParam)
	pageParamUsage      = fmt.Sprintf("--%s=10", pageParam)
	nParamUsage         = fmt.Sprintf("--%s=10", nParam)
	ordersIdsParamUsage = "<id заказа 1> ... <id заказа N>"
)

const (
	nParam       = "n"
	pageParam    = "page"
	sizeParam    = "size"
	userIdParam  = "user"
	expParam     = "exp"
	orderIdParam = "id"

	helpDescription = "Cправка"

	deliverOrderDescription = `На вход принимается ID заказа, ID получателя и срок хранения. Заказ нельзя принять дважды.`

	returnOrderDescription = `На вход принимается ID заказа. Метод должен удалять заказ из вашего файла. Можно вернуть только те заказы, у которых вышел срок хранения и если заказы находятся в пвз.`

	issueOrdersDescription = `Можно выдавать только те заказы, которые были приняты от курьера и чей срок хранения меньше текущей даты. Все ID заказов должны принадлежать только одному клиенту.`

	listOrdersDescription = `На вход принимается ID пользователя как обязательный параметр и опциональные параметры. Параметры позволяют получать только последние N заказов или заказы клиента, находящиеся в нашем ПВЗ.`

	refundOrderDescription = `На вход принимается ID пользователя и ID заказа. Заказ может быть возвращен в течение двух дней с момента выдачи.`

	listRefundedDescription = `Получить список всех заказов, которые вернули клиенты: Метод должен выдавать список пагинированно.`

	workersDescription = "Изменить максимальное количество горутин"

	exitDescription = `Завершить выполнение`
)

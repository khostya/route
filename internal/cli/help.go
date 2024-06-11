package cli

const (
	helpDescription = "Cправка"

	deliverOrderUsage       = "deliver --id=<id заказа> --user=<id> --exp=2024-06-09T17:12:32+05:00"
	deliverOrderDescription = `На вход принимается ID заказа, ID получателя и срок хранения. Заказ нельзя принять дважды.`

	returnOrderUsage       = "return --id=<id заказа>"
	returnOrderDescription = `На вход принимается ID заказа. Метод должен удалять заказ из вашего файла. Можно вернуть только те заказы, у которых вышел срок хранения и если заказы находятся в пвз.`

	issueOrdersUsage       = "issue <id заказа 1> ... <id заказа N>"
	issueOrdersDescription = `Можно выдавать только те заказы, которые были приняты от курьера и чей срок хранения меньше текущей даты. Все ID заказов должны принадлежать только одному клиенту.`

	listOrdersUsage       = "list --user=<id>"
	listOrdersDescription = `На вход принимается ID пользователя как обязательный параметр и опциональные параметры. Параметры позволяют получать только последние N заказов или заказы клиента, находящиеся в нашем ПВЗ.`

	refundOrderUsage       = "refund --id=<id заказа> --user=<id>"
	refundOrderDescription = `На вход принимается ID пользователя и ID заказа. Заказ может быть возвращен в течение двух дней с момента выдачи.`

	listRefundedUsage       = "refunded --size=20 --page=10"
	listRefundedDescription = `Получить список всех заказов, которые вернули клиенты: Метод должен выдавать список пагинированно.`

	procsUsage       = "procs --n=10"
	procsDescription = "Изменить максимальное количество горутин"

	exitDescription = `Завершить выполнение`
)

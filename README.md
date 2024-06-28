## Примеры команд
```
help
```
```
deliver --id=1 --user=1 --exp=2024-06-25T03:32:02+05:00 --wrapper=<[box package stretch]> --weight_in_kg=10.3 --price_in_rub=10.3
```
```
list --user=1
```
```
issue 1
```
```
refund --id=1 --user=1
```
```
refunded --size=20 --page=1
```
```
return --id=1
```
```
procs --n=10
```
```
exit
```
## Домашнее задание №4 «Выдача заказов в разной упаковке»
Описание UML https://www.uml-diagrams.org/
Открывать в https://app.diagrams.net/
https://drive.google.com/file/d/1eI_kM_SfUT-Yih2sY3-p-XprEcHXEsOO/view?usp=sharing
### Цель:

Модифицировать ваш сервис, добавить возможность в ПВЗ выдавать заказы в любой из трех различных упаковок

### Основное задание:

- **Модифицируйте Go-приложение**, добавьте в метод "Принять заказ от курьера" возможность передавать параметр упаковки
- Всего есть три вида упаковки: пакет, коробка, пленка
- Реализуйте функционал так, чтобы в будущем можно было просто добавить еще один вид упаковки
- При выборе пакета необходимо проверять, что вес заказа меньше 10 кг, если нет, то возвращаем информативную ошибку
- При выборе пакета стоимость заказа увеличивается на 5 рублей
- При выборе коробки необходимо проверить, что вес заказа меньше 30 кг, если нет, то возвращаем информативную ошибку
- При выборе коробки стоимость заказа увеличивается на 20 рублей
- При выборе пленки дополнительных проверок не требуется
- При выборе пленки стоимость заказа увеличивается на 1 рубль

### Дополнительное задание:

- Опишите архитектуру своего решения любым известным стандартом (например, UML)
- В MR вложите файл с описанием архитектуры
- При выборе стандарта необходимо описать, какой был выбран стандарт и дать ссылку на его документацию
- Запрещается использовать генерилки диаграмм, а также инструменты генерации связей между таблицами в БД в качестве описания

### Дедлайны сдачи и проверки задания:
- 22 июня 23:59 (сдача) / 25 июня, 23:59 (проверка)

## Запрос
```postgresql
SELECT id, recipient_id, status, status_updated_at, expiration_date, hash, created_at
FROM ozon.orders
where status = 'delivered' and recipient_id = '1' and id = any('{1}')
order by created_at desc
limit 1
offset 0;
```

```
Без индексов
```
https://explain.tensor.ru/archive/explain/87fac1780751587d62913933111db0b9:0:2024-06-15

```postgresql
create index on ozon.orders using btree(created_at);
```
https://explain.tensor.ru/archive/explain/feadc442b3d6aadc89ea8d8337e4ea98:0:2024-06-15

```postgresql
create index on ozon.orders using btree(recipient_id);
```
https://explain.tensor.ru/archive/explain/f29d8fa9c4e1e076fbe690af341aefae:0:2024-06-15

```postgresql
create index on ozon.orders using btree(recipient_id);
create index on ozon.orders using btree(created_at);
```
https://explain.tensor.ru/archive/explain/e3b604ef4f363f75d602f06f58bfbd58:0:2024-06-15

```postgresql
create index on ozon.orders using btree(status, recipient_id, id);
```
https://explain.tensor.ru/archive/explain/efcc29899aa66af31134dd3df00a957f:0:2024-06-15

Домашнее задание №3 «Рефакторинг слоя базы данных»

Основное задание

Цель:
Модифицируйте приложение, написанное в "Домашнее задание №2", чтобы взаимодействие с хранением данных было через Postgres, а не через файл.

Задание:

Переведите ваше приложение с хранения данных в файле на Postgres.
Реализуйте миграцию для DDL операторов.
Используйте транзакции.


Дополнительное задание:

Проанализируйте запросы в БД. Приложите результаты анализа в README.md. Добавьте индексы, где это необходимо.


Подсказки

Помните, что в одном файле миграции должен находиться один DDL оператор.
Для анализа плана запросов используйте Explain Tensor.


Дедлайны сдачи и проверки задания:

15 июня 23:59 (сдача) / 18 июня, 23:59 (проверка)

## Домашнее задание №2 «Ускорение обработки данных»
### Цель:

Модифицируйте приложение, написанное в "Домашнем задании №1", чтобы процесс обработки данных занимал меньше времени.

### Задание:

- Приложение должно работать в независимом режиме. Пользователь может передавать следующую команду, не дожидаясь результатов предыдущей.
- Ускорьте работу вашего приложения. Реализуйте обработку заданий в два потока, используя один из изученных паттернов конкурентной разработки.
- Приложение должно использовать механизм блокировки (например, мьютекс) для синхронизации доступа к данным между процессами.

### Дополнительные задания:

- Добавьте обработку системных сигналов, таких как SIGINT (Ctrl+C) и SIGTERM, чтобы корректно завершить работу приложения. При получении сигнала приложение должно завершить все задачи до выхода из консоли (graceful shutdown).
- Добавьте управление кол-вом горутин через отдельную команду для выполнения. Т.е. кол-во горутин должно изменяться динамически, без рестарта приложения.
- Добавить нотификации о начале и окончании обработки команды для каждой горутины. Реализовать в отдельной горутине.

### Подсказки:

- Для реагирования на системные сигналы вы можете использовать пакет os/signal в Go.
- Для синхронизации доступа к файлу между процессами записи и чтения вы можете использовать пакет sync и его мьютексы.
- Для работы с файлами в Go вы можете использовать пакет os или io/ioutil.

### Дедлайны сдачи и проверки задания:
- 8 июня 23:59 (сдача) / 11 июня, 23:59 (проверка)

## Домашнее задание №1 «Утилита для управления ПВЗ»
Необходимо реализовать консольную утилиту для менеджера ПВЗ.

Программа должна обладать командой help, благодаря которой можно получить список доступных команд с кратким описанием.

Список команд для реализации:

1. **Принять заказ от курьера**
   На вход принимается ID заказа, ID получателя и срок хранения. Заказ нельзя принять дважды. Если срок хранения в прошлом, приложение должно выдать ошибку. Список принятых заказов необходимо сохранять в файл. Формат файла остается на выбор автора.
2. **Вернуть заказ курьеру**
   На вход принимается ID заказа. Метод должен удалять заказ из вашего файла. Можно вернуть только те заказы, у которых вышел срок хранения и если заказы не были выданы клиенту.
3. **Выдать заказ клиенту**
   На вход принимается список ID заказов. Можно выдавать только те заказы, которые были приняты от курьера и чей срок хранения меньше текущей даты. Все ID заказов должны принадлежать только одному клиенту.
4. **Получить список заказов**
   На вход принимается ID пользователя как обязательный параметр и опциональные параметры.
   Параметры позволяют получать только последние N заказов или заказы клиента, находящиеся в нашем ПВЗ.
5. **Принять возврат от клиента**
   На вход принимается ID пользователя и ID заказа. Заказ может быть возвращен в течение двух дней с момента выдачи. Также необходимо проверить, что заказ выдавался с нашего ПВЗ.
6. **Получить список возвратов**
   Метод должен выдавать список пагинированно.


### Подсказки
- `os.OpenFile`, `os.Create`, `os.ReadALL`

### Дедлайны сдачи и проверки задания:
- 1 июня 23:59 (сдача) / 4 июня, 23:59 (проверка)

### Дополнительные задания. Реализация для получения 10 баллов

1. Программа запускается и ожидает поступление команд
2. Пользователь набирает команду и нажимает клавишу Enter
3. После нажатия клавиши Enter, программа выполняет введенную команду и записывает изменения в файл
4. Программа ожидает ввода новой команды
5. Программа завершает выполнение после нажатия комбинации клавиш ctrl+c или ввода команды exit
6. При выполнении команды в выходной файл записывается/обновляется дополнительное поле "hash"
7. Дополнительное поле "hash" считается с помощью функции `hash.GenerateHash()`  
   Функция GenerateHash() - это сторонняя функция, которую нужно импортировать и использовать, реализовывать её не нужно, она предоставляется преподавателям или тьютором.
   Функция GenerateHash() работает "долго", например, несколько секунд. Использование функции замедляет выполнение команды и не позволяет ввести следующую команду.
8. Придумайте свой/свои интерфейсы и используйте их в вашей программе. Продемонстрируйте усвоение учебного материала об интерфейсах

# Микросервис транзакций для avito-tech

## Задача
Необходимо реализовать микросервис для работы с балансом пользователей 
(зачисление средств, списание средств, перевод средств от пользователя к пользователю, а также метод получения баланса пользователя).
Сервис должен предоставлять HTTP API и принимать/отдавать запросы/ответы в формате JSON.

## Дополнительные задачи
Реализовать способ получения информации о транзакциях пользователя.

# Реализация
![image](https://user-images.githubusercontent.com/78440425/158826602-e4d4f3fa-9831-453f-8468-c8cda9b3d829.png)
Приложение реализует 4 метода:
 Узнать баланс пользователя.
 Узнать последние транзакции пользователя.
 Совершить перевод денег от одного пользователя к другому.
 Изменить баланс пользователя.

База данных представляет из себя 2 таблицы: таблица пользователей(wallets) и таблица транзакций(transactions).
Таблица пользователей состоит из 2 колонок: id, balance - идентификатор кошелька и его баланс.
Таблица транзакций состоит из 4 колонок: SenderId, ReceiverId, Sum, Status.
SenderId - идентификатор кошелька отправителя, с него деньи снимаются.
ReceiverId - идентификатор кошелька получателя, на него деньи зачисляются.
Sum - сумма перевода.
Status - более понятный статус ошибки в формате varchar, string,
по этому статусу можно понять, на каком именно этапе обработки произошла ошибка и как завершился каждый из этапов танзакции.

![image](https://user-images.githubusercontent.com/78440425/158829745-c99f1bd2-8945-4211-9c16-f0fd890f7245.png)

## Получение баланса пользователя, получение информации о транзакциях
![image](https://user-images.githubusercontent.com/78440425/158827342-d2cb3a3b-a7be-4e6e-9562-63ee15fec34e.png)
#### Получение баланса
your_server/getbalance - адрес метода.
Метод : HTTP POST с json-телом:

  {
      "Id":"your_id"
  }

Пример ответа:

![image](https://user-images.githubusercontent.com/78440425/158831404-0a298751-1186-4822-989b-3827b5ef3163.png)

#### Получение истории транзакций
your_server/gettransactions - адрес метода.
Метод : HTTP POST с json-телом:

  {
      "Id":"your_id"
  }

Пример ответа:

![image](https://user-images.githubusercontent.com/78440425/158831286-2b34150c-6fc3-4e18-b6b3-299368ea72c9.png)

## Переводы между кошельками, изменение баланса пользователя

![image](https://user-images.githubusercontent.com/78440425/158832119-b220d630-5d9a-41cd-91f8-8986a5f5d844.png)

#### Транзакции
your_server/maketransaction - адрес метода.
Метод: HTTP POST с json-телом:

  {
      "Sender": "2",
      "Receiver": "1",
      "Sum": 30,
      "Status": "not implemented"
  }

Пример ответа: 

![image](https://user-images.githubusercontent.com/78440425/158832685-01640ed3-91a8-4ba7-b953-fedbf3771868.png)

#### Изменение баланса
your_server/changebalance - адрес метода.
Метод: HTTP POST с json-телом:

  {
      "Sender": "1",
      "Sum": 2000,
      "Status": "not implemented"
  }

"Sender" это id того пользователя, с балансом которого мы будем работать.
"Sum" может быть как положительным, так и отрицательным.

Пример ответа:

![image](https://user-images.githubusercontent.com/78440425/158833359-e1dae6a3-5860-4675-9c86-5276844baa2b.png)

## Дополнение
В работе использовал язык Go, а также PostgreSQL.
В Go использовал такие модули как:
 net/http для работы http запросами и ответами, их обработкой и созданием.
 gorilla/mux для маршрутизации.
 encoding/json для работы с json, поскольку это формат запросов и ответов.
 database/sql для работы с базой данных.

Написал Unit-тесты для всех функций из файла "app/functions/functions.go"
Не сделал хорошего способа авторизации, авторизация по заголовку Key в header-е запроса.
Не упаковал в контейнер, проблемы с установкой докера на компьютер.

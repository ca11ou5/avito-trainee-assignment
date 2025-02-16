# **Тестовое задание для стажёра Backend-направления (зимняя волна 2025)**

## Магазин мерча

Был реализован внутренний магазин мерча, где сотрудники могут приобретать товары за монеты (coin). Каждому новому сотруднику выделяется 1000 монет, которые можно использовать для покупки товаров. Кроме того, монеты можно передавать другим сотрудникам в знак благодарности или как подарок.

## Инструкция по запуску
С конфигурацией docker-compose файла можно ознакомиться по ./configs/envs/local.env.  
Сам docker-compose файл располагается в ./deployments.  
В проекте описан Makefile, облегчающий выполнение инструкций.

### Запустить приложение локально ###
Для запуска приложения необходим свободный :8080 порт.
```
make compose-up
```

### Запустить настроенный линтер ###
(Конфигурация располагается в .golangci.yml).
```
make lint
```

### Запустить Unit-тесты ###
(Исключая E2E).
```
make tests
```

### Вывести отчет по покрытию кода Unit-тестами ###
(Исключая E2E). Покрытие тестами составило 45.9%.
```
make testscover
```

### Запустить E2E тесты ###
E2E тесты были реализованы для каждого эндпоинта, но только для успешных кейсов, неудачные и corner-кейсы постарался затронуть в Unit'ах. Необходимо подождать около 30ти секунд запуска приложения локально.
```
make e2e-tests
```

## Вопросы/проблемы ##
Хочу уточнить что для поднятия миграций в базе данных, использовалась утилита migrate, которая запускается отдельным сервисом в docker-compose файле и связана bridge сетью только с самой базой данных.    

С проектированием схемы базы данных можно ознакомиться в файле миграций
'./internal/adapters/secondary/postgres/migrations/init.up.sql'.  

В качестве вспомогательных библиотек для тестирования и мокирование использовались:
a) [github.com/stretchr/testify](https://github.com/stretchr/testify) - для мокирования своих интерфейсов.  
b) net/http/httptest - для тестирования обработчиков.  
с) [github.com/DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) - для мокирования базы данных.

Сильно серьезных проблем или вопросов не возникло, но некоторые заставили призадуматься. Вот самые основные:
- Как в случае покупки сотрудником мерча или переводом монет другому сотруднику атомарно сохранять данные в базе?
  - Решение: совершать такие операции в рамках одной транзакции, в случае чего делать Rollback, также использовать CTE (Common Table Expressions) для уменьшения количества запросов в БД и повышения атомарности операции.
- Как поступить с некоторыми данными, на которые очевидно требуется валидация?
  - Решение: сделал валидацию на свой вкус, как посчитал нужным (0 < len(username) <= 255 / 8 < len(password) <= 72 / etc.
- Проводить ли работу с трассировкой и логированием в целом?
  - Решение: использовал свой когда-то написанный пакет для трассировки и логирования, с возможностью записи в Graylog, поставил на это HTTP Middleware, но в логике приложения логи нигде не ставил, для меня все было понятно и без них, но в случае чего, это легко можно подкрутить, трассировка организуется через контекст, который прокидывается в ходе работы всего обработчика запроса.
- Стоит ли декомпозировать работу с JWT и хешированием в отдельный слой, или оставить в слое бизнес-логики?
  - Решение: сначала я реализовал это в слое бизнес-логики, потому что уточнений по этому поводу не было, но это неправильно если работа с авторизацией будет проходить в этом же сервисе. После для более комфортного мокирования зависимостей перенес эту работу в слой Adapter'ов (Repository).

## Вывод

В рамках поставленной задачи удалось реализовать поставленное задание, включая дополнительные задания, кроме нагрузочного тестирования (немного не хватило времени).
Используемый язык программирования - Go 1.23, база данных - PostgreSQL. Все проблемы и corner-кейсы, показавшиеся при тестировании и без него, были решены.

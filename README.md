## Задание
Необходимо написать веб сервис который занимается хранением профилей пользователей и их авторизацией.

Зависимости сервиса:
  - Для выполнения запросов к БД (github.com/jmoiron/sqlx v1.4.0)
  - Драйвер БД (github.com/lib/pq v1.10.9)
  - Для генерации uuid (github.com/google/uuid v1.6.0)
  - Для хеширования паролей (golang.org/x/crypto v0.40.0)

Для тестов:
  - gofakeit для генерации тестовых данных
  - testify для assert/require
  - gomock для генерации моков для юнит тестов

Зависимости команд из Makefile:
  - tsp 1.2.1 [ссылка](https://typespec.io/)
  - ogen 1.14.0 [ссылка](https://ogen.dev/blog/ogen-intro/)
  - mockgen v0.5.2 [ссылка](https://github.com/uber-go/mock)
  - goose v3.24.3 [ссылка](https://github.com/pressly/goose)

Документация:
  - API описан в **spec/main.tsp**
  - openapi спецификация: **openapi.yaml**
  - postman коллекция: **test_task.postman_collection.json**

### Запуск

```bash
git clone https://github.com/liriquew/test_task.git && \
cd test_task && \
docker compose up -d --build
```

## Модель
Профиль имеет набор полей:
1. id (uuid, unique)
2. email
3. username (unique)
4. password
5. admin (bool)

### Ограничения:
 - Длина имени пользователя должна быть больше 8 и состоять из английских букв и цифр
 - Длина пароля должна быть больше 8 и включать английские строчные и заглавные буквы и цифры
 - Почта должна быть валидной почтой

## API
У сервиса должeн быть набор ручек (rest, json):
```
GET /user (выдача листинга пользователей)
```

Создание, выдача, изменение и удаление профиля, доступно только пользователям с **Admin:true**
```
POST /users (создание пользователя)
GET /user/{id} (выдача пользователя)
PATCH /users/{id} (частичная замена в модели пользователя)
PUT /users/{id} (полная замена модели пользователя)
DELETE /users/{id} (удаление пользователя)
```
id - строковое представление uuid, например **"019763a9-7fc4-7e1a-9756-41c2ec1b998"**

## Механизм аутентификации
Сeрвис использует basic access authentication [ссылка](https://en.wikipedia.org/wiki/Basic_access_authentication)

При запросе, клиент указывает в заголовке Authorization следующую схему:
```
Basic base64encode(username:password)
```

base64: [ссылка](https://www.base64encode.org/)

### Админ
Изначально, при старте, в хранилище создается админ
```go
Username: "admin",
Email:    "admin@admin.ru",
Password: "admin",
```

Для аутентификации под этим админом надо указать в заголовке:
```
Authorization: Basic YWRtaW46YWRtaW4=
```

## Хранение
Для хранения используется postgres 17.5

### Интерфейс, требуемый для работы сервиса

```go
// internal/service/service.go
type Repository interface {
	ListUsers(context.Context, int64) ([]domain.User, error)

	CreateUser(context.Context, *domain.User) (*domain.UUID, error)
	GetUserById(context.Context, domain.UUID) (*domain.User, error)
	UpdateUser(context.Context, *domain.User) error
	DeleteUser(context.Context, domain.UUID) error

	GetUserByUsername(context.Context, string) (*domain.User, error)
}
```

## Тестирование

### Юнит тесты
Для тестирования корректности валидации написаны юнит тесты, для подмены базы данных в тестах используется gomock
```go
// internal/service/service.go

//go:generate mockgen -source=service.go -destination=mocks/repository.go -package=mocks
type Repository interface {
	...
}
```

Для генерации моков есть команда в Makefile
```bash
make gen_mocks
```

### e2e тесты
Для проверки работоспособности были написаны сквозные (e2e) тесты. Для избежания дублирования кода, была написана функция, **DoRequest**, которая скрывает детали взаимодействия. Условно, выполнение функции можно разделить на этапы:

- подготовка параметров запроса
- выполнение запроса
- извлечение тела ответа и проверка статус код ответа

Проверка корректности возвращенных значений выполняется вне этой функции

```go
// tests/e2e_test.go
func DoRequest(t *testing.T, method, url string, body any, header map[string]string, code int, respBody any) {
	// method - http method GET, POST ...
	// body - any struct, if not nil, will be json encoded before request
	// header - set of kv pairs
	// code - expected http status code
	// respBody - any struct, if not nil, response body will be decoded into it
}
```

Запуск тестов из корня проекта
```bash
make test # запустит e2e и юнит тесты
# или
go test ./tests/* -count 1 -v
```

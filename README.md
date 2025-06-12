## Задание
Необходимо написать веб сервис который занимается хранением профилей пользователей и их авторизацией.

Зависимости:
  - Роутер chi, полностью совместим с net/http (github.com/go-chi/chi/v5 v5.2.1)
  - Для генерации uuid (github.com/google/uuid v1.6.0)

Для тестов:
  - gofakeit для генерации тестовых данных
  - testify для assert/require

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

## API
У сервиса должeн быть набор ручек (rest, json):
```
GET /user (выдача листинга пользователей)
```

Выдача профиля по id, изменение и удаление профиля, доступно только пользователям с **Admin:true**
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

В случае, если в запросе не был указан нужный заголовок, в заголовке ответа будет указано следующее:
```
WWW-Authenticate: Basic realm="user service"
```
Что означает, что пользователю нужно ввести данные аутентификации в **user service**

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
Для хранения данных профилей необходимо реализовать примитивную in memory базу данных

### Интерфейс, требуемый для работы сервиса

```go
// internal/service/service.go
type Repository interface {
	ListUsers() []models.User

	// crud
	CreateUser(models.User) (*uuid.UUID, error)
	GetUserById(uuid.UUID) (*models.User, error)
	UpdateUser(models.User) error
	DeleteUser(uuid.UUID) error

	// используется в middleware
	GetUserByUsername(string) (*models.User, error)
}
```

### Стуктура, имплементирующая базу данных

В каждой операции чтения берется блокировка на чтение, при каждой операции записи берется блокировка на запись.

Перед возвратом записи о пользователе, запись копируется, благодаря этому изменения в разных запросах будут атомарны (одно изменение перезапишет другое)
```go
// internal/storage/storage.go
type Storage struct {
	users     map[uuid.UUID]*models.User // uuid -> user
	usernames map[string]uuid.UUID // username -> uuid

	m *sync.RWMutex
}
```

## Тестирование
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
go test ./tests/* -count 1 -v
```

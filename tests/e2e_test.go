package tests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"encoding/base64"
	"encoding/json"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/google/uuid"
	domain "github.com/liriquew/test_task/internal/domain"
	"github.com/liriquew/test_task/internal/lib/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cfg config.AppTestConfig

type Id struct {
	Id uuid.UUID `json:"id"`
}

func TestMain(m *testing.M) {
	cfg = config.MustLoadPathTest("../config/test_config.yaml")

	m.Run()
}

func Copy(u *domain.User) *domain.User {
	return &domain.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
		IsAdmin:  u.IsAdmin,
	}
}

func GetRandomUser() *domain.User {
	return &domain.User{
		Username: domain.NewOptString(gofakeit.Username() + gofakeit.AchAccount()),
		Password: domain.NewOptString(gofakeit.Password(true, true, true, false, false, 12) + "1" + "a" + "A"),
		Email:    domain.NewOptString(gofakeit.Email()),
		IsAdmin: domain.OptBool{
			Value: false,
			Set:   true,
		},
	}
}

func GetDefaultAdmin() *domain.User {
	return &domain.User{
		Username: domain.NewOptString("admin"),
		Email:    domain.NewOptString("admin@admin.ru"),
		Password: domain.NewOptString("admin"),
		IsAdmin: domain.OptBool{
			Value: true,
			Set:   true,
		},
	}
}

func GetAuthHeader(user *domain.User) map[string]string {
	value := base64.StdEncoding.EncodeToString(fmt.Appendf(nil,
		"%s:%s", user.Username.Value, user.Password.Value),
	)

	return map[string]string{
		"Authorization": "Basic " + value,
		"Content-Type":  "application/json",
	}
}

func DoRequest(t *testing.T, method, url string, body any, header map[string]string, code int, respBody any) {
	// method - http method GET, POST ...
	// body - any struct, will be json encoded before request
	// header - set of kv pairs
	// code - expected http status code
	// respBody - any struct, if not nil, response body will be decoded into it
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, fmt.Sprintf(
		"http://%s:%d/%s", cfg.API.Host, cfg.API.Port, url,
	), reqBody)
	for k, v := range header {
		req.Header.Add(k, v)
	}
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	if code != 0 {
		if ok := assert.Equal(t, code, resp.StatusCode); !ok {
			body, err := io.ReadAll(resp.Body)
			fmt.Println(string(body), err)
			t.Fail()
		}
	}

	if respBody != nil {
		err := json.NewDecoder(resp.Body).Decode(respBody)
		require.NoError(t, err)
	}
}

func CreateUser(t *testing.T, user *domain.User) uuid.UUID {
	const url = "users/"
	var id Id
	DoRequest(t, "POST", url, user, GetAuthHeader(GetDefaultAdmin()), 201, &id)
	return id.Id
}

func GetUser(t *testing.T, id uuid.UUID) domain.User {
	url := fmt.Sprintf("users/%s", id.String())
	user := domain.User{}
	DoRequest(t, "GET", url, nil, GetAuthHeader(GetDefaultAdmin()), 200, &user)
	return user
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	url := "users/"
	user := GetRandomUser()

	t.Run("New user", func(t *testing.T) {
		t.Parallel()
		var id Id
		DoRequest(t, "POST", url, user, GetAuthHeader(GetDefaultAdmin()), 201, &id)
		require.NotZero(t, id.Id)
	})

	t.Run("Conflict", func(t *testing.T) {
		t.Parallel()
		var id Id
		user := GetRandomUser()

		DoRequest(t, "POST", url, user, GetAuthHeader(GetDefaultAdmin()), 201, &id)
		DoRequest(t, "POST", url, user, GetAuthHeader(GetDefaultAdmin()), 409, nil)
	})

	t.Run("Forbidden", func(t *testing.T) {
		t.Parallel()
		var id Id
		user := GetRandomUser()

		DoRequest(t, "POST", url, user, GetAuthHeader(GetDefaultAdmin()), 201, &id)
		DoRequest(t, "POST", url, user, GetAuthHeader(user), 403, nil)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		t.Parallel()
		DoRequest(t, "POST", url, user, nil, 401, nil)
	})
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	var id Id
	user := GetRandomUser()
	DoRequest(t, "POST", "users/", user, GetAuthHeader(GetDefaultAdmin()), 201, &id)
	user.ID.SetTo(domain.UUID(id.Id))

	t.Run("New user", func(t *testing.T) {
		t.Parallel()
		url := fmt.Sprintf("users/%s", id.Id.String())
		newUser := domain.User{}
		DoRequest(t, "GET", url, nil, GetAuthHeader(GetDefaultAdmin()), 200, &newUser)
		require.NotZero(t, id.Id)
		newUser.Password.Value = ""
		userToCheck := Copy(user)
		userToCheck.Password.Value = ""

		assert.Equal(t, *userToCheck, newUser)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		t.Parallel()
		url := fmt.Sprintf("users/%s", id.Id.String())
		DoRequest(t, "GET", url, nil, nil, 401, nil)
	})

	t.Run("Not found", func(t *testing.T) {
		t.Parallel()

		id, _ := uuid.NewV7()
		url := fmt.Sprintf("users/%s", id.String())
		DoRequest(t, "GET", url, nil, GetAuthHeader(GetDefaultAdmin()), 404, nil)
	})

	t.Run("Forbidden", func(t *testing.T) {
		t.Parallel()

		id, _ := uuid.NewV7()
		url := fmt.Sprintf("users/%s", id.String())
		DoRequest(t, "GET", url, nil, GetAuthHeader(user), 403, nil)
	})
}

func TestPatchUser(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		user := GetRandomUser()
		id := CreateUser(t, user)
		user.ID.SetTo(domain.UUID(id))

		url := fmt.Sprintf("users/%s", id.String())
		user.Username.SetTo(gofakeit.Username() + gofakeit.AchAccount())
		user.Email.SetTo(gofakeit.Email())
		DoRequest(t, "PATCH", url, &domain.User{
			Username: user.Username,
			Email:    user.Email,
		}, GetAuthHeader(GetDefaultAdmin()), 200, nil)

		patchedUser := GetUser(t, id)
		user.Password.Value = ""
		patchedUser.Password.Value = ""

		assert.Equal(t, *user, patchedUser)
	})
}

func TestPutUser(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		user := GetRandomUser()
		id := CreateUser(t, user)
		user.ID.SetTo(domain.UUID(id))

		newUser := GetRandomUser()
		newUser.ID = user.ID

		url := fmt.Sprintf("users/%s", id.String())
		DoRequest(t, "PUT", url, newUser, GetAuthHeader(GetDefaultAdmin()), 200, nil)

		updatedUser := GetUser(t, id)
		newUser.Password.Value = ""
		updatedUser.Password.Value = ""
		assert.Equal(t, *newUser, updatedUser)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		user := GetRandomUser()
		id := CreateUser(t, user)
		user.ID.SetTo(domain.UUID(id))

		url := fmt.Sprintf("users/%s", id.String())
		DoRequest(t, "DELETE", url, nil, GetAuthHeader(GetDefaultAdmin()), 200, nil)

		DoRequest(t, "GET", url, nil, GetAuthHeader(GetDefaultAdmin()), 404, nil)
	})
}

func TestListUsers(t *testing.T) {
	t.Parallel()

	const cnt = 10

	usersSet := make(map[uuid.UUID]int, 0)
	for range cnt {
		user := GetRandomUser()

		id := CreateUser(t, user)
		user.ID.SetTo(domain.UUID(id))

		user.Password.Value = ""
		usersSet[id] = 0
	}

	var resp []domain.User
	DoRequest(t, "GET", "users/", nil, GetAuthHeader(GetDefaultAdmin()), 200, &resp)

	for _, user := range resp {
		id := uuid.UUID(user.ID.Value)
		if _, ok := usersSet[id]; ok {
			usersSet[id]++
		}
	}

	for _, v := range usersSet {
		require.Equal(t, 1, v)
	}
}

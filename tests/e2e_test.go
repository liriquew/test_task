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
	"github.com/liriquew/test_task/internal/lib/config"
	"github.com/liriquew/test_task/internal/models"
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

func GetRandomUser() *models.User {
	return &models.User{
		Username: gofakeit.Username(),
		Password: gofakeit.AppName(),
		Email:    gofakeit.Email(),
		Admin: models.Bool{
			Ok: true,
		},
	}
}

func GetAuthHeader(user *models.User) map[string]string {
	value := base64.StdEncoding.EncodeToString(fmt.Appendf(nil,
		"%s:%s", user.Username, user.Password),
	)

	return map[string]string{
		"Authorization": "Basic " + value,
	}
}

func CreateUser(t *testing.T, user *models.User) uuid.UUID {
	const url = "users/"
	var id Id
	DoRequest(t, "POST", url, user, GetAuthHeader(models.GetDefaultAdmin()), 201, &id)
	return id.Id
}

func GetUser(t *testing.T, id uuid.UUID) models.User {
	url := fmt.Sprintf("users/%s", id.String())
	user := models.User{}
	DoRequest(t, "GET", url, nil, GetAuthHeader(models.GetDefaultAdmin()), 200, &user)
	return user
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
		require.Equal(t, code, resp.StatusCode)
	}

	if respBody != nil {
		err := json.NewDecoder(resp.Body).Decode(respBody)
		require.NoError(t, err)
	}

	return
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	url := "users/"
	user := GetRandomUser()

	t.Run("New user", func(t *testing.T) {
		t.Parallel()
		var id Id
		DoRequest(t, "POST", url, user, GetAuthHeader(models.GetDefaultAdmin()), 201, &id)
		require.NotZero(t, id.Id)
	})

	t.Run("Conflict", func(t *testing.T) {
		t.Parallel()
		var id Id
		user := GetRandomUser()

		DoRequest(t, "POST", url, user, GetAuthHeader(models.GetDefaultAdmin()), 201, &id)
		DoRequest(t, "POST", url, user, GetAuthHeader(models.GetDefaultAdmin()), 409, nil)
	})

	t.Run("Forbidden", func(t *testing.T) {
		t.Parallel()
		var id Id
		user := GetRandomUser()

		DoRequest(t, "POST", url, user, GetAuthHeader(models.GetDefaultAdmin()), 201, &id)
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
	DoRequest(t, "POST", "users/", user, GetAuthHeader(models.GetDefaultAdmin()), 201, &id)
	user.Id = id.Id

	t.Run("New user", func(t *testing.T) {
		t.Parallel()
		url := fmt.Sprintf("users/%s", id.Id.String())
		newUser := models.User{}
		DoRequest(t, "GET", url, nil, GetAuthHeader(user), 200, &newUser)
		require.NotZero(t, id.Id)

		userCp := user.Copy()
		userCp.Password = ""
		assert.Equal(t, *userCp, newUser)
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
		DoRequest(t, "GET", url, nil, GetAuthHeader(user), 404, nil)
	})
}

func TestPatchUser(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		user := GetRandomUser()
		user.Id = CreateUser(t, user)

		url := fmt.Sprintf("users/%s", user.Id.String())
		user.Username = gofakeit.Username()
		user.Email = gofakeit.Email()
		DoRequest(t, "PATCH", url, models.User{
			Username: user.Username,
			Email:    user.Email,
		}, GetAuthHeader(models.GetDefaultAdmin()), 200, nil)

		patchedUser := GetUser(t, user.Id)
		user.Password = ""
		assert.Equal(t, *user, patchedUser)
	})
}

func TestPutUser(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		user := GetRandomUser()
		user.Id = CreateUser(t, user)

		newUser := GetRandomUser()
		newUser.Id = user.Id

		url := fmt.Sprintf("users/%s", user.Id.String())
		DoRequest(t, "PUT", url, newUser, GetAuthHeader(models.GetDefaultAdmin()), 200, nil)

		updatedUser := GetUser(t, user.Id)
		newUser.Password = ""
		assert.Equal(t, *newUser, updatedUser)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		user := GetRandomUser()
		user.Id = CreateUser(t, user)

		url := fmt.Sprintf("users/%s", user.Id.String())
		DoRequest(t, "DELETE", url, nil, GetAuthHeader(models.GetDefaultAdmin()), 200, nil)

		DoRequest(t, "GET", url, nil, GetAuthHeader(models.GetDefaultAdmin()), 404, nil)
	})
}

func TestListUsers(t *testing.T) {
	t.Parallel()

	const cnt = 10

	usersSet := make(map[models.User]int, 0)
	for range cnt {
		user := GetRandomUser()

		user.Id = CreateUser(t, user)

		user.Password = ""
		usersSet[*user] = 0
	}

	var resp []models.User
	DoRequest(t, "GET", "users/", nil, GetAuthHeader(models.GetDefaultAdmin()), 200, &resp)

	for _, user := range resp {
		usersSet[user]++
	}

	for _, v := range usersSet {
		require.Equal(t, 1, v)
	}
}

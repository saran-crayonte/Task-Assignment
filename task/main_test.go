package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/saran-crayonte/task/database"
	"github.com/saran-crayonte/task/models"
	"github.com/saran-crayonte/task/routes"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	database.ConnectDB()
	app := fiber.New()
	routes.SetupRoutes(app)

	user := models.User{
		Username: "Test User 5",
		Name:     "tester4",
		Email:    "tester4email",
		Password: "test",
	}

	payload, _ := json.Marshal(user)

	userReq := httptest.NewRequest(http.MethodPost, "/api/user", bytes.NewReader(payload))
	taskResp, err := app.Test(userReq)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, taskResp.StatusCode)
}

func TestLogin(t *testing.T) {
	database.ConnectDB()
	app := fiber.New()
	routes.SetupRoutes(app)

	user := models.User{
		Username: "Test User 4",
		Password: "test",
	}
	payload, _ := json.Marshal(user)

	userReq := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(payload))
	userResp, err := app.Test(userReq)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusAccepted, userResp.StatusCode)
	//token := userResp.Header.Get("token")
	//assert.NotEmpty(t, token)
}

//const Tokenn = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImFrZUBnbWFpbC5jb20iLCJVc2VybmFtZSI6ImthcmFuIiwiZXhwIjoxNzA3NDY5MzkzfQ.Ff2ZwcnRSOnr5kPPg5OxJ7CNaOq9n4FaqntJqsDh7QE"

func TestGetTask(t *testing.T) {
	database.ConnectDB()
	app := fiber.New()
	routes.SetupRoutes(app)
	task := models.Task{
		ID: 1,
	}
	payload, _ := json.Marshal(task)

	getTasksReq := httptest.NewRequest(http.MethodGet, "/api/v2/task/id", bytes.NewReader(payload))

	getTasksReq.Header.Set("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImFrZUBnbWFpbC5jb20iLCJVc2VybmFtZSI6ImthcmFuIiwiZXhwIjoxNzA3NDY5MzkzfQ.Ff2ZwcnRSOnr5kPPg5OxJ7CNaOq9n4FaqntJqsDh7QE")

	getTasksResp, err := app.Test(getTasksReq)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, getTasksResp.StatusCode)
}

func TestCreateTask(t *testing.T) {
	database.ConnectDB()
	app := fiber.New()
	routes.SetupRoutes(app)
	task := models.Task{
		Title:          "sample1",
		Status:         "pending",
		EstimatedHours: 50,
	}
	payload, _ := json.Marshal(task)
	TasksReq := httptest.NewRequest(http.MethodPost, "/api/v2/task", bytes.NewReader(payload))
	TasksReq.Header.Set("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImFrZUBnbWFpbC5jb20iLCJVc2VybmFtZSI6ImthcmFuIiwiZXhwIjoxNzA3NDY5MzkzfQ.Ff2ZwcnRSOnr5kPPg5OxJ7CNaOq9n4FaqntJqsDh7QE")
	TasksResp, err := app.Test(TasksReq)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, TasksResp.StatusCode)
}

func TestUpdateTask(t *testing.T) {
	database.ConnectDB()
	app := fiber.New()
	routes.SetupRoutes(app)
	task := models.Task{
		ID:             5,
		Title:          "sample1",
		Status:         "Inprogress",
		EstimatedHours: 50,
	}
	payload, _ := json.Marshal(task)
	TasksReq := httptest.NewRequest(http.MethodPut, "/api/v2/task/id", bytes.NewReader(payload))
	TasksReq.Header.Set("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImFrZUBnbWFpbC5jb20iLCJVc2VybmFtZSI6ImthcmFuIiwiZXhwIjoxNzA3NDY5MzkzfQ.Ff2ZwcnRSOnr5kPPg5OxJ7CNaOq9n4FaqntJqsDh7QE")
	TasksResp, err := app.Test(TasksReq)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, TasksResp.StatusCode)
}

func TestDeleteTask(t *testing.T) {
	database.ConnectDB()
	app := fiber.New()
	routes.SetupRoutes(app)
	task := models.Task{
		ID: 4,
	}
	payload, _ := json.Marshal(task)
	TasksReq := httptest.NewRequest(http.MethodDelete, "/api/v2/task/id", bytes.NewReader(payload))
	TasksReq.Header.Set("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImFrZUBnbWFpbC5jb20iLCJVc2VybmFtZSI6ImthcmFuIiwiZXhwIjoxNzA3NDY5MzkzfQ.Ff2ZwcnRSOnr5kPPg5OxJ7CNaOq9n4FaqntJqsDh7QE")
	TasksResp, err := app.Test(TasksReq)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, TasksResp.StatusCode)
}

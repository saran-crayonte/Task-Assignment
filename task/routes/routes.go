package routes

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/saran-crayonte/task/database"
	"github.com/saran-crayonte/task/models"
	"github.com/saran-crayonte/task/user"
)

func SetupRoutes(app *fiber.App) {

	ap := app.Group("/api")
	ap.Post("/user", user.Register())
	ap.Post("/user/login", user.Login())

	api := ap.Group("/v2", user.Authenticate())

	api.Get("/refreshToken", user.RefreshToken())

	api.Put("/user", user.UpdatePassword())

	api.Post("/task", CreateTasks)
	api.Get("/task/id", GetTasks)
	api.Put("/task/id", UpdateTasks)
	api.Delete("/task/id", DeleteTasks)

	api.Post("/taskAssignment", CreateTaskAssignment)
	api.Get("/taskAssignment/id", GetTaskAssignment)
	api.Put("/taskAssignment/id", UpdateTaskAssignment)
	api.Delete("/taskAssignment/id", DeleteTaskAssignment)

	api.Post("/holiday", CreateHoliday)
	api.Get("/holiday/id", GetHoliday)
	api.Put("/holiday/id", UpdateHoliday)
	api.Delete("/holiday/id", DeleteHoliday)

	// api.Post("/user", CreateUser)
	// api.Post("/user/login", LoginUser)
	// api.Put("/user", UpdateUser)

}

func CreateTasks(c fiber.Ctx) error {
	task := new(models.Task)
	if err := json.Unmarshal(c.Body(), &task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	var existingTask models.Task
	database.DB.Where("title = ?", task.Title).First(&existingTask)
	if existingTask.ID != 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Task with the same title already exists"})
	}

	database.DB.Create(&task)
	return c.Status(fiber.StatusCreated).JSON(task)
}

func GetTasks(c fiber.Ctx) error {
	task := new(models.Task)
	if err := json.Unmarshal(c.Body(), &task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	var newTask models.Task
	database.DB.First(&newTask, task.ID)
	if newTask.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}
	return c.Status(fiber.StatusOK).JSON(newTask)
}
func UpdateTasks(c fiber.Ctx) error {
	task := new(models.Task)
	if err := json.Unmarshal(c.Body(), &task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})

	}

	var existingTask models.Task
	database.DB.First(&existingTask, task.ID)
	if existingTask.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	// newData := new(models.Task)
	// if err := json.Unmarshal(c.Body(), &newData); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	// }

	database.DB.Model(&existingTask).Updates(task)
	return c.Status(fiber.StatusOK).JSON(existingTask)

}
func DeleteTasks(c fiber.Ctx) error {
	task := new(models.Task)
	if err := json.Unmarshal(c.Body(), &task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	var newTask models.Task
	database.DB.First(&newTask, task.ID)
	if newTask.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	database.DB.Delete(&newTask)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task deleted successfully",
	})
}

func CreateTaskAssignment(c fiber.Ctx) error {
	taskAssignment := new(models.TaskAssignment)
	if err := json.Unmarshal(c.Body(), &taskAssignment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	var existingUser models.User
	database.DB.First(&existingUser, "username = ?", taskAssignment.Username)
	if len(existingUser.Username) == 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username doesn't exists"})
	}

	var existingTask models.Task
	database.DB.First(&existingTask, taskAssignment.TaskID)
	if existingTask.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	//to check if task already assigned
	var checkAlreadyAssigned models.TaskAssignment
	database.DB.First(&checkAlreadyAssigned, taskAssignment.TaskID)
	if checkAlreadyAssigned.ID != 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task is already assigned to somebody"})
	}

	estimatedHours := existingTask.EstimatedHours
	layout := "2006-01-02 3:04 PM"
	startDate, err := time.Parse(layout, taskAssignment.Start_Date)
	// startDate, err := time.Parse(time.RFC3339, taskAssignment.Start_Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid date time format"})
	}
	result := calculateEndDate(startDate, estimatedHours)
	/*
		startDate, err := time.Parse("2006-01-02", taskAssignment.Start_Date)
		if err != nil {
			return err
		}

		startDate = startDate.Truncate(24 * time.Hour)

		for hours := estimatedHours; hours > 0; {
			if startDate.Weekday() == time.Saturday ||
				startDate.Weekday() == time.Sunday ||
				isHoliday(startDate) {
				startDate = startDate.AddDate(0, 0, 1)
				continue
			}
			hours -= 8
			startDate = startDate.AddDate(0, 0, 1)
		}
	*/
	taskAssignment.Start_Date = startDate.Format("2006-01-02 3:04 PM")
	taskAssignment.End_Date = result.Format("2006-01-02 3:04 PM")
	database.DB.Create(taskAssignment)
	return c.JSON(taskAssignment)
}
func calculateEndDate(startDate time.Time, estimatedHours int) time.Time {
	//workingHoursPerDay := 8
	endDate := startDate
	remainingHours := estimatedHours

	for remainingHours > 0 {

		if endDate.Weekday() == time.Saturday || endDate.Weekday() == time.Sunday || isHoliday(endDate) {
			endDate = endDate.AddDate(0, 0, 1)
			continue
		}

		if endDate.Hour() == 12 {
			endDate = endDate.Add(time.Hour)
			continue
		}

		if endDate.Hour() >= 18 {
			endDate = endDate.AddDate(0, 0, 1).Truncate(24 * time.Hour).Add(9 * time.Hour)
			continue
		}

		remainingHours--
		endDate = endDate.Add(time.Hour)
	}

	if endDate.Hour() > 18 {
		endDate = endDate.AddDate(0, 0, 1).Truncate(24 * time.Hour).Add(9 * time.Hour)
	}

	return endDate
}
func isHoliday(date time.Time) bool {
	holiday := new(models.Holiday)
	if err := database.DB.Where("holiday_date = ?", date.Format("2006-01-02")).First(holiday).Error; err != nil {
		return false
	}
	return true
}

func GetTaskAssignment(c fiber.Ctx) error {
	taskAssignment := new(models.TaskAssignment)
	if err := json.Unmarshal(c.Body(), &taskAssignment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	var newTaskAssignment models.TaskAssignment
	database.DB.First(&newTaskAssignment, taskAssignment.ID)
	if newTaskAssignment.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task Assignment not found"})
	}
	return c.JSON(newTaskAssignment)
}

func UpdateTaskAssignment(c fiber.Ctx) error {
	taskAssignment := new(models.TaskAssignment)
	if err := json.Unmarshal(c.Body(), &taskAssignment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	var existingUser models.User
	database.DB.First(&existingUser, "username = ?", taskAssignment.Username)
	if len(existingUser.Username) == 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username doesn't exists"})
	}

	var existingTask models.Task
	database.DB.First(&existingTask, taskAssignment.TaskID)
	if existingTask.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	//to check if task already assigned
	var checkAlreadyAssigned models.TaskAssignment
	database.DB.First(&checkAlreadyAssigned, taskAssignment.TaskID)
	if checkAlreadyAssigned.ID != 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task is already assigned to somebody"})
	}

	estimatedHours := existingTask.EstimatedHours
	layout := "2006-01-02 3:04 PM"
	startDate, err := time.Parse(layout, taskAssignment.Start_Date)
	// startDate, err := time.Parse(time.RFC3339, taskAssignment.Start_Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid date time format"})
	}
	result := calculateEndDate(startDate, estimatedHours)
	taskAssignment.Start_Date = startDate.Format("2006-01-02 3:04 PM")
	taskAssignment.End_Date = result.Format("2006-01-02 3:04 PM")

	var existingTaskAssignment models.TaskAssignment
	database.DB.First(&existingTaskAssignment, taskAssignment.ID)
	if existingTaskAssignment.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task Assignment not found"})
	}

	database.DB.Model(&existingTaskAssignment).Updates(taskAssignment)
	return c.JSON(existingTaskAssignment)
}

func DeleteTaskAssignment(c fiber.Ctx) error {
	taskAssignment := new(models.TaskAssignment)
	if err := json.Unmarshal(c.Body(), &taskAssignment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	var existingTaskAssignment models.TaskAssignment
	database.DB.First(&existingTaskAssignment, taskAssignment.ID)
	if existingTaskAssignment.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task Assignment not found"})
	}

	database.DB.Delete(&existingTaskAssignment)
	return c.JSON(fiber.Map{
		"message": "Task Assignment entry deleted successfully",
	})
}

func CreateHoliday(c fiber.Ctx) error {
	holiday := new(models.Holiday)
	if err := json.Unmarshal(c.Body(), &holiday); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	var newHoliday models.Holiday
	database.DB.First(&newHoliday, holiday.HolidayDate)
	if newHoliday.ID != 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Holiday already defined"})
	}
	database.DB.Create(&holiday)
	return c.Status(fiber.StatusCreated).JSON(holiday)
}
func GetHoliday(c fiber.Ctx) error {
	holiday := new(models.Holiday)
	if err := json.Unmarshal(c.Body(), &holiday); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	var newHoliday models.Holiday
	database.DB.First(&newHoliday, holiday.ID)
	if newHoliday.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Holiday not found"})
	}
	return c.JSON(newHoliday)
}
func UpdateHoliday(c fiber.Ctx) error {
	holiday := new(models.Holiday)
	if err := json.Unmarshal(c.Body(), &holiday); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	var newHoliday models.Holiday
	database.DB.First(&newHoliday, holiday.ID)
	if newHoliday.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Holiday not found"})
	}
	database.DB.Model(&newHoliday).Updates(holiday)
	return c.JSON(newHoliday)
}
func DeleteHoliday(c fiber.Ctx) error {
	holiday := new(models.Holiday)
	if err := json.Unmarshal(c.Body(), &holiday); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	var newHoliday models.Holiday
	database.DB.First(&newHoliday, holiday.ID)
	if newHoliday.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Holiday not found"})
	}
	database.DB.Delete(&newHoliday)
	return c.JSON(fiber.Map{
		"message": "Holiday deleted successfully",
	})
}

/*
func CreateUser(c fiber.Ctx) error {
	user := new(models.User)
	if err := json.Unmarshal(c.Body(), &user); err != nil {
		return err
	}
	var existingUser models.User
	database.DB.First(&existingUser, "username = ?", user.Username)
	if len(existingUser.Username) != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "this username already exists"})
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "problem in hashing password"})
	}

	newData := models.User{Username: user.Username, Name: user.Name, Email: user.Email, Password: string(hash)}

	database.DB.Create(newData)
	return c.JSON(newData)
}

func LoginUser(c fiber.Ctx) error {
	user := new(models.User)
	if err := json.Unmarshal(c.Body(), &user); err != nil {
		return err
	}
	var existingUser models.User
	database.DB.First(&existingUser, "username = ?", user.Username)
	if existingUser.Username != user.Username {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid username"})
	}
	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user.Name,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret_key_12345"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token string could not be created"})
	}
	return c.SendString("Token: " + tokenString)
}

func UpdateUser(c fiber.Ctx) error {
	user := new(models.User)
	if err := json.Unmarshal(c.Body(), &user); err != nil {
		return err
	}
	var existingUser models.User
	database.DB.First(&existingUser, "username = ?", user.Username)
	if existingUser.Username != user.Username {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username doesn't exists"})
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "problem in hashing password"})
	}

	newData := models.User{Username: user.Username, Name: user.Name, Email: user.Email, Password: string(hash)}
	database.DB.Model(&existingUser).Updates(newData)
	return c.JSON(existingUser)
}
*/

// func CreateUser(c fiber.Ctx) error {
// 	user := new(models.User)
// 	if err := json.Unmarshal(c.Body(), &user); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
// 	}

// 	var existingUser models.User
// 	if err := database.DB.First(&existingUser, "username = ?", user.Username).Error; err == nil {
// 		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
// 	}

// 	database.DB.Create(&user)
// 	return c.Status(fiber.StatusCreated).JSON(user)
// }

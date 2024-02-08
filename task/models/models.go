package models

type Task struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Title          string `gorm:"not null" json:"title"`
	Status         string `gorm:"not null" json:"status"`
	EstimatedHours int    `gorm:"not null" json:"estimatedHours"`
}

type TaskAssignment struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Username   string `gorm:"not null" json:"username"`
	TaskID     uint   `gorm:"not null" json:"taskid"`
	Start_Date string `gorm:"not null" json:"startDate"`
	End_Date   string `json:"endDate"`
}

type Holiday struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	HolidayName string `gorm:"not null" json:"holidayName"`
	HolidayDate string `gorm:"not null" json:"holidayDate"`
}

type User struct {
	Username string `gorm:"primaryKey;uniqueIndex;not null" json:"username"`
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
}

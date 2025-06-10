package models

import "time"

type User struct {
	ID        string
	Name      string
	Username  string
	Password  string
	Salary    float64
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
	CreatedBy string
	UpdatedBy string
}

type UserSalary struct {
	UserID string
	Salary float64
}

package user

import (
	"database/sql"
)

type SQLUser struct {
	ID        sql.NullString
	Name      sql.NullString
	Username  sql.NullString
	Password  sql.NullString
	Salary    sql.NullFloat64
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
	CreatedBy sql.NullString `db:"created_by"`
	UpdatedBy sql.NullString `db:"updated_by"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

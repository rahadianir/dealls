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
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
	CreatedBy sql.NullString
	UpdatedBy sql.NullString
}

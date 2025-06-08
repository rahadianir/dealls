package migrations

import (
	"log"
	"math/rand/v2"
	"os"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func SetupData() {
	// connect to db
	dsn := os.Getenv("DB_URL")
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to database for data setup: ", err)
	}
	defer db.Close()

	generateMockData(db)

}

func generateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.IntN(len(letters))]
	}
	return string(result)
}

func generateRandomName(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, rand.IntN(length)+1)
	for i := range result {
		result[i] = letters[rand.IntN(len(letters))]
	}
	return string(result)
}

func generateMockData(db *sqlx.DB) {
	tx, err := db.Begin()
	if err != nil {
	}
	defer tx.Rollback()

	// insert admin data
	// setup admin role
	roleID := uuid.New()
	q := `INSERT INTO hr.roles (id, name, created_at) VALUES ($1, $2, now())`
	_, err = tx.Exec(q, roleID, "admin")
	if err != nil {
		log.Fatal("failed to insert admin role: ", err)
	}

	// insert admin user
	adminPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	userID := uuid.New()
	salary := rand.Int64N(100000000)

	pwBytes, err := bcrypt.GenerateFromPassword([]byte(adminPassword), 12)
	if err != nil {
		log.Fatal("failed to hash admin password: ", err)
	}

	q = `INSERT INTO hr.users (id, name, username, password, salary, created_at) VALUES ($1, $2, $2, $3, $4, now())`
	_, err = tx.Exec(q, userID, "admin", string(pwBytes), salary)
	if err != nil {
		log.Fatal("failed to insert admin data: ", err)
	}

	// assign admin role to admin user
	mapID := uuid.New()
	q = `INSERT INTO hr.user_role_map (id, user_id, role_id, created_at) VALUES ($1, $2, $3, now())`
	_, err = tx.Exec(q, mapID, userID, roleID)
	if err != nil {
		log.Fatal("failed to assign admin role to admin user: ", err)
	}

	for i := 1; i <= 100; i++ {
		id := uuid.New()
		name := generateRandomName(20)
		password := generateRandomString(8)
		salary := rand.Int64N(100000000)

		pwBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			log.Println("failed to hash password")
			continue
		}

		q := `INSERT INTO hr.users (id, name, username, password, salary, created_at) VALUES ($1, $2, $2, $3, $4, now())`
		_, err = tx.Exec(q, id, name, string(pwBytes), salary)
		if err != nil {
			log.Println("failed to insert data")
			continue
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal("failed to commit mock data insertions: ", err)
	}
}

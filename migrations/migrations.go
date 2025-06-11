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
		log.Fatal("failed to start db transaction: ", err)
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
	log.Println("inserting admin data")
	adminPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	adminID := "46bf7246-6186-45cc-aa7f-2c7fd8b32c81"
	salary := rand.Int64N(100000000)

	pwBytes, err := bcrypt.GenerateFromPassword([]byte(adminPassword), 12)
	if err != nil {
		log.Fatal("failed to hash admin password: ", err)
	}

	q = `INSERT INTO hr.users (id, name, username, password, salary, created_at) VALUES ($1, $2, $2, $3, $4, now())`
	_, err = tx.Exec(q, adminID, "admin", string(pwBytes), salary)
	if err != nil {
		log.Fatal("failed to insert admin data: ", err)
	}

	// assign admin role to admin user
	mapID := uuid.New()
	q = `INSERT INTO hr.user_role_map (id, user_id, role_id, created_at) VALUES ($1, $2, $3, now())`
	_, err = tx.Exec(q, mapID, adminID, roleID)
	if err != nil {
		log.Fatal("failed to assign admin role to admin user: ", err)
	}

	// prepare 3 static employee data for testing
	password := "password" // static password for easier use to test
	pwBytes, err = bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatal("failed to hash test user password: ", err)
	}

	// static UUID for testing
	id1 := "81d1bcd4-d5b3-4495-92ce-ef2c9b0f5e54"
	id2 := "8f29acd8-c18a-4e1c-9662-f102562bc893"
	id3 := "cc3a57a3-79cf-438e-9dc3-3a18bd86480b"

	// inserting 3 static employee data for testing
	q = `INSERT INTO hr.users (id, name, username, password, salary, created_at) VALUES 
	($1, 'ani', 'ani', $4, 10000000, now()),
	($2, 'budi', 'budi', $4, 13000000, now()),
	($3, 'coki', 'coki', $4, 17000000, now())`
	_, err = tx.Exec(q, id1, id2, id3, string(pwBytes))
	if err != nil {
		log.Fatal("failed to insert static employee data: ", err)
	}

	for i := 1; i <= 100; i++ {
		log.Println("inserting fake employee data: ", i)
		id := uuid.New()
		name := generateRandomName(20)
		salary := rand.Int64N(100000000)

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

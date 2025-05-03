package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"titles.run/site_api/services/interfaces"

	strava "titles.run/strava/models"
)

type DBServiceImpl struct {
	db *sql.DB
}

func NewDBService() interfaces.DBService {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	maxRetries := 20
	retryInterval := 5 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Println("unable to open connection to database")
			return nil
		}

		if err := db.Ping(); err != nil {
			log.Printf("unable to reach the database: %v\n", err)
		} else {
			return &DBServiceImpl{db}
		}

		log.Printf("Database not ready (attempt %d/%d): %v\n", attempt, maxRetries, err)
		time.Sleep(retryInterval)
	}

	log.Fatalln("Database not ready after max retries")
	return nil
}

func (d *DBServiceImpl) Close() {
	d.db.Close()
}

func (d *DBServiceImpl) AcceptTerms(userID int64) error {
	_, err := d.db.Exec("UPDATE users SET plan = 'free', ai = true, terms_accepted = true WHERE id = $1;", userID)
	if err != nil {
		log.Println("error updating terms accepted:", err)
		return errors.New("could not update terms accepted")
	}
	return nil
}

func (d *DBServiceImpl) GetUser(userID int64) (strava.User, error) {
	user := strava.User{}

	err := d.db.QueryRow("SELECT id, name, pic, ai, plan, terms_accepted FROM users WHERE id = $1;", userID).Scan(&user.ID, &user.Name, &user.Pic, &user.AI, &user.Plan, &user.TermsAccepted)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no user found for given ID:", userID)
			return user, errors.New("no user found for GetUser")
		}
		log.Println("error querying database:", err)
		return user, errors.New("could not query database for GetUser")
	}

	return user, nil
}

func (d *DBServiceImpl) GetUserInternal(userID int64) (strava.UserInternal, error) {
	user := strava.UserInternal{}

	err := d.db.QueryRow("SELECT id, name, pic, access_token, refresh_token, expires_at, ai, plan, terms_accepted FROM users WHERE id = $1;", userID).Scan(&user.ID, &user.Name, &user.Pic, &user.AccessToken, &user.RefreshToken, &user.ExpiresAt, &user.AI, &user.Plan, &user.TermsAccepted)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no user found for given ID:", userID)
			return user, errors.New("no user found for GetUser")
		}
		log.Println("error querying database:", err)
		return user, errors.New("could not query database for GetUser")
	}

	return user, nil
}

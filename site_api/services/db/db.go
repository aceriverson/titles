package db

import (
	"database/sql"
	"encoding/json"
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

	defaultSettings := strava.Settings{
		AutomaticTitle: true,
		Tone:           50,
		Attribution:    true,
		Description:    false,
	}
	data, err := json.Marshal(defaultSettings)
	if err != nil {
		return err
	}

	_, err = d.db.Exec("UPDATE user_settings SET settings = $1 WHERE user_id = $2;", data, userID)
	if err != nil {
		log.Println("error updating user settings:", err)
		return errors.New("could not update user settings")
	}

	return nil
}

func (d *DBServiceImpl) GetUser(userID int64) (strava.User, error) {
	user := strava.User{}
	var settingsJSON []byte

	err := d.db.QueryRow(`
		SELECT u.id, u.name, u.pic, u.ai, u.plan, u.terms_accepted, COALESCE(s.settings, '{}'::jsonb)
		FROM users u
		LEFT JOIN user_settings s ON u.id = s.user_id
		WHERE u.id = $1;
	`, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Pic,
		&user.AI,
		&user.Plan,
		&user.TermsAccepted,
		&settingsJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no user found for given ID:", userID)
			return user, errors.New("no user found for GetUser")
		}
		log.Println("error querying database:", err)
		return user, errors.New("could not query database for GetUser")
	}

	if err := json.Unmarshal(settingsJSON, &user.Settings); err != nil {
		log.Println("failed to unmarshal settings JSON:", err)
		return user, errors.New("invalid user settings data")
	}

	return user, nil
}

func (d *DBServiceImpl) UpdateSettings(userID int64, settings strava.Settings) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(`
		INSERT INTO user_settings (user_id, settings)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET settings = EXCLUDED.settings
	`, userID, data)

	return err
}

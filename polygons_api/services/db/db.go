package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"titles.run/services/interfaces"
	strava "titles.run/strava/models"
	"titles.run/titles/models"

	_ "github.com/lib/pq"
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

func (d *DBServiceImpl) DeletePolygon(userID int64, polygon models.Polygon) error {
	_, err := d.db.Exec("DELETE FROM polygons WHERE user_id = $1 AND id = $2;", userID, polygon.ID)
	if err != nil {
		log.Println("error deleting polygon:", err)
		return err
	}

	return nil
}

func (d *DBServiceImpl) GetUser(userID int64) (strava.User, error) {
	user := strava.User{}

	err := d.db.QueryRow("SELECT id, name, pic, plan, terms_accepted FROM users WHERE id = $1;", userID).Scan(&user.ID, &user.Name, &user.Pic, &user.Plan, &user.TermsAccepted)
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

func (d *DBServiceImpl) GetPolygons(userID int64) ([]models.Polygon, error) {
	rows, err := d.db.Query("SELECT id, name, ST_AsText(geom) AS geom_text FROM polygons WHERE user_id = $1;", userID)
	if err != nil {
		log.Println("error querying polygons:", err)
		return nil, err
	}
	defer rows.Close()

	polygons := []models.Polygon{}
	for rows.Next() {
		var id, name, geomText string
		if err := rows.Scan(&id, &name, &geomText); err != nil {
			log.Println("error scanning row:", err)
			return nil, err
		}

		polygon := models.Polygon{
			ID:   id,
			Name: name,
		}

		if err := polygon.ParseWKT(geomText); err != nil {
			log.Println("error parsing WKT:", err)
			return nil, err
		}

		polygons = append(polygons, polygon)
	}

	if err := rows.Err(); err != nil {
		log.Println("error in row iteration:", err)
		return nil, err
	}

	return polygons, nil
}

func (d *DBServiceImpl) PostPolygon(userID int64, polygon models.Polygon) error {
	_, err := d.db.Exec("INSERT INTO polygons (user_id, name, geom) VALUES ($1, $2, $3);", userID, polygon.Name, polygon.ToWKT())
	if err != nil {
		log.Println("error inserting polygon:", err)
		return err
	}

	return nil
}

func (d *DBServiceImpl) PutPolygon(userID int64, polygon models.Polygon) error {
	_, err := d.db.Exec("UPDATE polygons SET geom = $1 WHERE user_id = $2 AND id = $3;", polygon.ToWKT(), userID, polygon.ID)
	if err != nil {
		log.Println("error updating polygon:", err)
		return err
	}

	return nil
}

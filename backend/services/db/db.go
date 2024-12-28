package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"titles.run/services/interfaces"
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

func (d *DBServiceImpl) GetIntersectingPolygons(userID int64, points [][]float64) ([]models.Polygon, error) {
	polystr := "LINESTRING ("
	for _, p := range points {
		polystr += fmt.Sprintf("%f %f, ", p[1], p[0])
	}
	polystr = polystr[:len(polystr)-2] + ")"

	rows, err := d.db.Query("WITH input_linestring AS (SELECT ST_GeomFromText($1, 4326) AS geom) SELECT polygons.name as name FROM polygons JOIN input_linestring ON ST_Intersects(input_linestring.geom, polygons.geom) WHERE user_id=$2 ORDER BY ST_LineLocatePoint(input_linestring.geom, ST_StartPoint(polygons.geom)) DESC;", polystr, userID)
	if err != nil {
		log.Println("error querying polygons:", err)
		return nil, err
	}

	polygons := []models.Polygon{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Println("error scanning row:", err)
			return nil, err
		}

		polygon := models.Polygon{
			Name: name,
		}

		polygons = append(polygons, polygon)
	}

	return polygons, nil
}

func (d *DBServiceImpl) GetUser(userID int64) (models.User, error) {
	user := models.User{}

	err := d.db.QueryRow("SELECT id, name, pic FROM users WHERE id = $1;", userID).Scan(&user.ID, &user.Name, &user.Pic)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no user found for given ID")
			return user, errors.New("no user found for GetUser")
		}
		log.Println("error querying database:", err)
		return user, errors.New("could not query database for GetUser")
	}

	return user, nil
}

func (d *DBServiceImpl) GetUserInternal(userID int64) (models.UserInternal, error) {
	user := models.UserInternal{}

	err := d.db.QueryRow("SELECT id, name, pic, access_token, refresh_token, expires_at, ai FROM users WHERE id = $1;", userID).Scan(&user.ID, &user.Name, &user.Pic)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no user found for given ID")
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

func (d *DBServiceImpl) NewUser(user models.UserInternal) error {
	_, err := d.db.Exec(
		"INSERT INTO users (id, name, pic, access_token, refresh_token, expires_at) VALUES ($1, $2, $3, $4, $5, $6);",
		user.ID, user.Name, user.Pic, user.AccessToken, user.RefreshToken, user.ExpiresAt,
	)
	if err != nil {
		log.Println("error inserting user:", err)
		return err
	}

	return nil
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

func (d *DBServiceImpl) UnauthorizeUser(userID int64) error {
	_, err := d.db.Exec("DELETE FROM users WHERE user_id = $1;", userID)
	if err != nil {
		log.Println("error updating polygon:", err)
		return err
	}

	return nil
}

func (d *DBServiceImpl) UpdateUser(user models.UserInternal) error {
	_, err := d.db.Exec(
		"UPDATE users SET name = $1, pic = $2, access_token = $3, refresh_token = $4, expires_at = $5 WHERE id = $6;",
		user.Name, user.Pic, user.AccessToken, user.RefreshToken, user.ExpiresAt, user.ID,
	)
	if err != nil {
		log.Println("error updating user:", err)
		return err
	}

	return nil
}

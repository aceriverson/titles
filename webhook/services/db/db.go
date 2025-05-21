package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"titles.run/webhook/models"
	"titles.run/webhook/services/interfaces"

	strava "titles.run/strava/models"

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

func (d *DBServiceImpl) CreateSubscription(userID int64, customer, subscription, customerEmail string) error {
	_, err := d.db.Exec(
		"INSERT INTO subscriptions (user_id, customer, subscription, email) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id) DO UPDATE SET customer = $2, subscription = $3, email = $4;",
		userID, customer, subscription, customerEmail,
	)
	if err != nil {
		log.Println("error inserting subscription:", err)
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

	rows, err := d.db.Query("WITH input_linestring AS (SELECT ST_GeomFromText($1, 4326) AS geom) SELECT polygons.name as name FROM polygons JOIN input_linestring ON ST_Intersects(input_linestring.geom, polygons.geom) WHERE user_id=$2 ORDER BY ST_LineLocatePoint(input_linestring.geom, ST_StartPoint(polygons.geom));", polystr, userID)
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

func (d *DBServiceImpl) GetPOI(points [][]float64) (models.POIs, error) {
	polystr := "LINESTRING ("
	for _, p := range points {
		polystr += fmt.Sprintf("%f %f, ", p[1], p[0])
	}
	polystr = polystr[:len(polystr)-2] + ")"

	rows, err := d.db.Query("SELECT title, lat, lng FROM poi WHERE active = true AND ST_DWithin(poi.geom, ST_SetSRID(ST_GeomFromText($1), 4326)::geography, 25);", polystr)
	if err != nil {
		log.Println("error querying polygons:", err)
		return models.POIs{}, err
	}

	pois := []models.POI{}
	for rows.Next() {
		var name string
		var lat, lon float64
		if err := rows.Scan(&name, &lat, &lon); err != nil {
			log.Println("error scanning row:", err)
			return models.POIs{}, err
		}

		poi := models.POI{
			Title: name,
			Position: models.Coordinates{
				Latitude:  lat,
				Longitude: lon,
			},
		}

		pois = append(pois, poi)
	}

	return models.POIs{Items: pois}, nil
}

func (d *DBServiceImpl) GetUserInternal(userID int64) (strava.UserInternal, error) {
	user := strava.UserInternal{}
	var settingsJSON []byte

	err := d.db.QueryRow(`
		SELECT 
			u.id, u.name, u.pic, u.access_token, u.refresh_token, u.expires_at, 
			u.ai, u.plan, u.terms_accepted, COALESCE(s.settings, '{}'::jsonb) 
		FROM users u 
		LEFT JOIN user_settings s ON u.id = s.user_id
		WHERE u.id = $1;
	`, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Pic,
		&user.AccessToken,
		&user.RefreshToken,
		&user.ExpiresAt,
		&user.AI,
		&user.Plan,
		&user.TermsAccepted,
		&settingsJSON)
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

func (d *DBServiceImpl) SetPOI(pois models.POIs) error {
	for _, poi := range pois.Items {
		_, err := d.db.Exec(
			"INSERT INTO poi (id, title, lat, lng) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING;",
			poi.ID, poi.Title, poi.Position.Latitude, poi.Position.Longitude,
		)
		if err != nil {
			log.Println("error inserting POI:", err)
			return err
		}
	}

	return nil
}

func (d *DBServiceImpl) UnauthorizeUser(userID int64) error {
	_, err := d.db.Exec("DELETE FROM users WHERE user_id = $1;", userID)
	if err != nil {
		log.Println("error deleting user:", err)
		return err
	}

	return nil
}

func (d *DBServiceImpl) UpdateSubscription(customer, subscription, plan string) error {
	_, err := d.db.Exec(`
		UPDATE users
		SET plan = $1
		FROM subscriptions
		WHERE subscriptions.customer = $2
  			AND subscriptions.subscription = $3
  			AND subscriptions.user_id = users.id;
	`, plan, customer, subscription)
	if err != nil {
		log.Println("error updating subscription:", err)
		return err
	}

	return nil
}

func (d *DBServiceImpl) UpdateUser(user strava.UserInternal) error {
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

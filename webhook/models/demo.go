package models

type Point struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type Gpx struct {
	Points []Point `json:"points"`
}

type TitleRequest struct {
	Title      string `json:"title"`
	ActivityID int64  `json:"activity"`
}

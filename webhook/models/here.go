package models

type Coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type POI struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Position struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lng"`
	} `json:"position"`
}

type POIs struct {
	Items []POI `json:"items"`
}

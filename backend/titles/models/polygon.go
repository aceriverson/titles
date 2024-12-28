package models

import (
	"errors"
	"fmt"
	"strings"
)

type Polygon struct {
	ID     string       `json:"id"`
	Name   string       `json:"name"`
	Points [][2]float64 `json:"points"`
}

func (p *Polygon) ParseWKT(wkt string) error {
	// Remove "POLYGON((" and "))"
	coords := strings.TrimPrefix(wkt, "POLYGON((")
	coords = strings.TrimSuffix(coords, "))")

	// Split by commas to get each "lng lat" pair
	pairs := strings.Split(coords, ",")
	points := make([][2]float64, len(pairs))

	for i, pair := range pairs {
		var lng, lat float64
		_, err := fmt.Sscanf(pair, "%f %f", &lng, &lat)
		if err != nil {
			return errors.New("failed to parse point")
		}
		points[i] = [2]float64{lng, lat}
	}

	p.Points = points
	return nil
}

func (p *Polygon) ToWKT() string {
	var sb strings.Builder
	sb.WriteString("POLYGON((")

	for i, point := range p.Points {
		sb.WriteString(fmt.Sprintf("%f %f", point[0], point[1]))
		if i < len(p.Points)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString("))")
	return sb.String()
}

package utils

import (
	"errors"

	"github.com/heremaps/flexible-polyline/golang/flexpolyline"
	"github.com/twpayne/go-polyline"
)

type Polyline struct {
	Points [][]float64
	Google string
	Flex   string
}

func PolylineFromGoogle(line string) *Polyline {
	p := &Polyline{}

	p.Google = line

	buf := []byte(line)
	coords, _, err := polyline.DecodeCoords(buf)
	if err != nil {
		coords = [][]float64{}
	}
	p.Points = coords

	points := make([]flexpolyline.Point, len(coords))
	for i, c := range coords {
		points[i] = flexpolyline.Point{Lat: c[0], Lng: c[1]}
	}
	t, _ := flexpolyline.CreatePolyline(3, points)
	flex, err := t.Encode()
	if err != nil {
		flex = ""
	}
	p.Flex = flex

	return p
}

func (p *Polyline) Validate() error {
	if len(p.Points) == 0 {
		return errors.New("invalid polyline")
	}

	if p.Google == "" {
		return errors.New("empty polyline")
	}

	if p.Flex == "" {
		return errors.New("invalid flexible polyline")
	}

	return nil
}

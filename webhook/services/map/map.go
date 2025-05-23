package map_service

import (
	"bytes"
	"encoding/base64"
	"image/color"
	"image/png"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"

	"titles.run/webhook/services/interfaces"
)

type MapServiceImpl struct {
}

func NewMapService() interfaces.MapService {
	return &MapServiceImpl{}
}

func (m *MapServiceImpl) GenerateMap(coords [][]float64) (string, error) {
	ctx := sm.NewContext()
	ctx.SetTileProvider(sm.NewTileProviderOpenStreetMaps())
	ctx.SetSize(720, 720)

	latlngs := make([]s2.LatLng, 0)
	for _, coord := range coords {
		latlngs = append(latlngs, s2.LatLngFromDegrees(coord[0], coord[1]))
	}

	ctx.AddObject(
		sm.NewPath(
			latlngs,
			color.RGBA{252, 82, 0, 255},
			3,
		),
	)

	ctx.AddObject(
		sm.NewCircle(
			latlngs[len(latlngs)-1],
			color.RGBA{252, 82, 0, 255},
			color.RGBA{255, 0, 0, 255},
			80,
			3,
		),
	)

	ctx.AddObject(
		sm.NewCircle(
			latlngs[0],
			color.RGBA{252, 82, 0, 255},
			color.RGBA{0, 255, 0, 255},
			80,
			3,
		),
	)

	img, err := ctx.Render()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	err = png.Encode(&buf, img)
	if err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64String, nil
}

module titles.run/webhook

go 1.23.1

replace titles.run/strava => ../shared/strava

require (
	github.com/golang/geo v0.0.0-20250417192230-a483f6ae7110
	github.com/jarcoal/httpmock v1.4.0
	github.com/lib/pq v1.10.9
	github.com/redis/go-redis/v9 v9.7.3
	titles.run/strava v0.0.0-00010101000000-000000000000
)

require (
	github.com/flopp/go-coordsparser v0.0.0-20250311184423-61a7ff62d17c // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/mazznoer/csscolorparser v0.1.5 // indirect
	github.com/tkrajina/gpxgo v1.4.0 // indirect
	golang.org/x/image v0.26.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/flopp/go-staticmaps v0.0.0-20250419123608-f8593a4338f2
	github.com/heremaps/flexible-polyline v0.1.0
	github.com/twpayne/go-polyline v1.1.1
)

module titles.run/site_api

go 1.23.1

replace titles.run/jwt => ../shared/jwt

replace titles.run/strava => ../shared/strava

replace titles.run/turnstile => ../shared/turnstile

require (
	github.com/lib/pq v1.10.9
	titles.run/jwt v0.0.0-00010101000000-000000000000
	titles.run/strava v0.0.0-00010101000000-000000000000
	titles.run/turnstile v0.0.0-00010101000000-000000000000
)

require github.com/golang-jwt/jwt/v5 v5.2.2 // indirect

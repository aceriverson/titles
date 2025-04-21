module titles.run

go 1.23.1

replace titles.run/auth => ../shared/auth

require (
	github.com/lib/pq v1.10.9
	titles.run/auth v0.0.0-00010101000000-000000000000
)

require github.com/golang-jwt/jwt/v5 v5.2.2 // indirect

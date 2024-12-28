CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR PRIMARY KEY,
    name VARCHAR,
    pic VARCHAR,
    access_token VARCHAR,
	refresh_token VARCHAR,
	expires_at BIGINT,
    ai BOOLEAN
);

CREATE TABLE IF NOT EXISTS polygons (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR REFERENCES users(id),
    name VARCHAR,
    geom geometry(Polygon, 4326)
);
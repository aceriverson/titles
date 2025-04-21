CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE plan AS ENUM ('none', 'free', 'pro');

CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    name VARCHAR,
    pic VARCHAR,
    access_token VARCHAR,
	refresh_token VARCHAR,
	expires_at BIGINT,
    ai BOOLEAN DEFAULT FALSE,
    plan plan DEFAULT 'none',
    terms_accepted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS polygons (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    name VARCHAR,
    geom geometry(Polygon, 4326)
);

CREATE TABLE IF NOT EXISTS poi (
    id VARCHAR PRIMARY KEY,
    title VARCHAR,
    lat DOUBLE PRECISION NOT NULL,
    lng DOUBLE PRECISION NOT NULL,
    geom GEOGRAPHY(POINT, 4326) GENERATED ALWAYS AS (
        ST_SetSRID(ST_MakePoint(lng, lat), 4326)::GEOGRAPHY
    ) STORED,
    active BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_poi_geom ON poi USING GIST (geom);
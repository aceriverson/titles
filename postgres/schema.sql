CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE plan AS ENUM ('none', 'free', 'pro');

CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,                 -- Strava ID
    name VARCHAR,                          -- Concatenated Strava name
    pic VARCHAR,                           -- URL to Strava profile picture
    access_token VARCHAR,                  -- Strava OAuth access token
	refresh_token VARCHAR,                 -- Strava OAuth refresh token
	expires_at BIGINT,                     -- Strava OAuth token expiration time
    ai BOOLEAN DEFAULT FALSE,              -- Where user has enabled AI, should be set to true when terms are accepted
    plan plan DEFAULT 'none',              -- User plan, should be set to 'free' when terms are accepted
    terms_accepted BOOLEAN DEFAULT FALSE,  -- Whether user has accepted terms
    created_at TIMESTAMP DEFAULT NOW()     -- User creation time
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

CREATE TABLE IF NOT EXISTS subscriptions (
    user_id BIGINT PRIMARY KEY REFERENCES users(id),
    customer VARCHAR,         -- Stripe customer ID
    subscription VARCHAR     -- Stripe subscription ID
);

CREATE TABLE IF NOT EXISTS user_settings (
    user_id BIGINT PRIMARY KEY REFERENCES users(id),
    settings JSONB DEFAULT '{}'::JSONB
);
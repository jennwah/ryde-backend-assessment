CREATE SCHEMA users;
CREATE EXTENSION IF NOT EXISTS postgis; -- Enable PostGIS extension

CREATE TABLE IF NOT EXISTS users.users
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(255)           NOT NULL,
    date_of_birth DATE                   NOT NULL,
    address       TEXT,
    description   TEXT,
    location      GEOGRAPHY(POINT, 4326) NOT NULL, -- Stores latitude and longitude in WGS 84
    created_at    TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (name)
);

CREATE INDEX IF NOT EXISTS idx_users_location ON users.users USING GIST (location);

CREATE TABLE IF NOT EXISTS users.friends
(
    user_id    UUID REFERENCES users.users (id) ON DELETE CASCADE NOT NULL,
    friend_id  UUID REFERENCES users.users (id) ON DELETE CASCADE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, friend_id),
    CONSTRAINT unique_friendship CHECK (user_id < friend_id) -- Enforces consistent ordering
);

CREATE TABLE IF NOT EXISTS weather (
    id SERIAL PRIMARY KEY,
    temperature FLOAT NOT NULL,
    humidity FLOAT NOT NULL,
    recorded_at DATE NOT NULL UNIQUE,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_weather_recorded_at ON weather (recorded_at);
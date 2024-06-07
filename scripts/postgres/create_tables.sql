BEGIN;

-- Create Users table
CREATE TABLE IF NOT EXISTS users (
    discordid TEXT NOT NULL PRIMARY KEY,
    world TEXT NOT NULL,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL
);

COMMIT;
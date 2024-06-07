BEGIN;

-- Example of how to seed some user data
INSERT INTO users (discordid, world, firstname, lastname)
VALUES
('TestDiscordUser', 'Leviathan', 'Test', 'User');

COMMIT;
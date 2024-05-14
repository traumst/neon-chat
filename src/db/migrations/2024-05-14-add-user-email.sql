PRAGMA foreign_keys=off;
BEGIN TRANSACTION;

ALTER TABLE users RENAME TO _users_old;

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    name TEXT, 
    email TEXT,
    type TEXT,
    salt INTEGER
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name ON users(name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);

INSERT INTO users 
SELECT id, 
    name, 
    concat(name, '@email.ok'), 
    type, 
    salt
FROM _users_old;

DROP TABLE _users_old;

COMMIT;
PRAGMA foreign_keys=on;
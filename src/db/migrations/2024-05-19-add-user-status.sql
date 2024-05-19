PRAGMA foreign_keys=off;
BEGIN TRANSACTION;

ALTER TABLE users RENAME TO _users_old;

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    name TEXT UNIQUE, 
    email TEXT UNIQUE,
    type TEXT,
    status TEXT,
    salt INTEGER
);

INSERT INTO users 
    SELECT id, 
        name, 
        email,
        type, 
        'active',
        salt
    FROM _users_old;

DROP TABLE _users_old;

COMMIT;
PRAGMA foreign_keys=on;
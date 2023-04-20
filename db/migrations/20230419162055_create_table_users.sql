-- migrate:up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    lastname TEXT NOT NULL
);


-- migrate:down
DROP TABLE users;

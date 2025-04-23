CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name TEXT NOT NULL,
                       username TEXT UNIQUE NOT NULL,
                       password TEXT NOT NULL,
                       role TEXT NOT NULL
);

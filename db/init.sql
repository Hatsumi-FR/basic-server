CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     email TEXT UNIQUE NOT NULL,
                                     password TEXT NOT NULL
);

INSERT INTO users (email, password)
VALUES ('alice@example.com', '1234')
    ON CONFLICT (email) DO NOTHING;

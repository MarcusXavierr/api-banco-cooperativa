CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR (256) NOT NULL,
    email VARCHAR (256) UNIQUE NOT NULL,
    password VARCHAR (128) NOT NULL,
    credit_limit INTEGER NOT NULL,
    balance INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    value INTEGER NOT NULL,
    type CHAR(1) NOT NULL,
    description VARCHAR(10),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
INSERT INTO users (name, credit_limit, email, password)
  VALUES
    ('Paulo 🇧🇷', 1000 * 100, 'paulo@gmail.com', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8'),
    ('Sujyro 🇯🇵', 800 * 100, 'sujyro@gmail.com', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8'),
    ('Giuseppe 🇮🇹', 10000 * 100, 'giuseppe@gmail.com', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8'),
    ('Jalan 🇮🇳', 100000 * 100, 'jalan@gmail.com', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8'),
    ('Jallim 🇸🇦', 5000 * 100, 'jallim@gmail.com', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8');
END; $$

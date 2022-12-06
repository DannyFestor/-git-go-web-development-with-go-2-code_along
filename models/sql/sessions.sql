CREATE TABLE sessions (
	id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users (id) ON DELETE CASCADE, -- references part alternative below
    token_hash TEXT UNIQUE NOT NULL--,
    --FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE -- alternative above
);
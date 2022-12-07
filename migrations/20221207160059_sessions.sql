-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions (
	id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users (id) ON DELETE CASCADE, -- references part alternative below
    token_hash TEXT UNIQUE NOT NULL--,
    --FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE -- alternative above
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username TEXT NOT NULL,
                                     first_name TEXT NOT NULL,
                                     middle_name TEXT DEFAULT '',
                                     last_name TEXT NOT NULL,
                                     email TEXT NOT NULL,
                                     gender CHAR(1) CHECK (gender IN ('M', 'F', 'O')),
                                     age SMALLINT CHECK (age >= 0 AND age <= 150)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

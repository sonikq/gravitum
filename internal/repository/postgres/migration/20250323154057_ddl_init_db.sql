-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username TEXT NOT NULL UNIQUE,
                                     first_name TEXT NOT NULL,
                                     middle_name TEXT DEFAULT '',
                                     last_name TEXT NOT NULL,
                                     email TEXT NOT NULL,
                                     gender CHAR(1) CHECK (gender IN ('M', 'F', 'O')),
    age SMALLINT CHECK (age >= 1 AND age <= 150),
    beg_date TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE
                           );

CREATE INDEX IF NOT EXISTS user_id_idx ON users (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

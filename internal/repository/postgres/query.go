package postgres

//CREATE TABLE IF NOT EXISTS users (
//id SERIAL PRIMARY KEY,
//username TEXT NOT NULL UNIQUE,
//first_name TEXT NOT NULL,
//middle_name TEXT DEFAULT '',
//last_name TEXT NOT NULL,
//email TEXT NOT NULL,
//gender CHAR(1) CHECK (gender IN ('M', 'F', 'O')),
//age SMALLINT CHECK (age >= 0 AND age <= 150)
//);

const (
	createUser = `insert into users(username, first_name, middle_name, last_name, email, gender, age, beg_date) values ($1, $2, $3, $4, $5, $6, $7, now()) returning id`
	getUser    = `select id, username, first_name, middle_name, last_name, email, gender, age, end_date from users where id = $1`
	updateUser = `update users set username = $1, first_name = $2,
middle_name = $3, last_name = $4, email = $5, gender = $6, age = $7, updated_at = now() where id = $8;`
	deleteUser = `update users set end_date = now() where id = $1`
)

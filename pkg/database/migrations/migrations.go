package migrations

import "fmt"

const usersTable = `
CREATE TABLE users(
id serial NOT NULL,
created_at timestamp NOT NULL DEFAULT current_timestamp,
updated_at timestamp DEFAULT current_timestamp,
email text NOT NULL,
password text NOT NULL,
role text[] NOT NULL DEFAULT '{"user"}',
verification_token text,
user_status	integer NOT NULL DEFAULT 1,
PRIMARY KEY (id)
);

INSERT INTO users (email, password) VALUES ('test@test.com', 'password1234');

SELECT * FROM users;

`

func init() {
	up := []string{
		usersTable,
	}

	down := []string{
		`DROP TABLE users`,
	}

	fmt.Println(up, down)
}

-- +goose Up

CREATE TABLE Users(
	id serial primary key,
	email varchar(50) NOT NULL UNIQUE,
	password varchar(100) NOT NULL,
    admin boolean DEFAULT FALSE
);

-- +goose Down

DROP TABLE Users;
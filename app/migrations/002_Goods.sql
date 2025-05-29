-- +goose Up

CREATE TABLE Goods(
	id serial primary key,
	description text,
    name varchar(30),
    active boolean NOT NULL,
    cost smallint NOT NULL,
    image_url varchar NOT NULL
);

-- +goose Down

DROP TABLE Goods;
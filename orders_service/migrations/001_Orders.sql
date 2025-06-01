-- +goose Up
CREATE TABLE Orders (
    order_uuid uuid NOT NULL,
    user_id smallint NOT NULL,
    good_id smallint NOT NULL,
    amount smallint NOT NULL,
    PRIMARY KEY (order_uuid, good_id)
);

-- +goose Down
DROP TABLE Orders;
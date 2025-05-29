-- +goose Up
CREATE TABLE Orders (
    order_uuid uuid NOT NULL,
    user_id smallint NOT NULL REFERENCES Users(id) ON DELETE CASCADE,
    good_id smallint NOT NULL REFERENCES Goods(id) ON DELETE CASCADE,
    amount smallint NOT NULL,
    PRIMARY KEY (order_uuid, good_id)
);

-- +goose Down
DROP TABLE Orders;
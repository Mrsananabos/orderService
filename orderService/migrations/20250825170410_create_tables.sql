-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS delivery (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    zip VARCHAR(10) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region VARCHAR(100),
    email VARCHAR(255)
    );

CREATE TABLE IF NOT EXISTS payment (
    id SERIAL PRIMARY KEY,
    transaction VARCHAR(255) NOT NULL,
    request_id VARCHAR(255) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(100) NOT NULL,
    amount FLOAT8 NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(100),
    delivery_cost FLOAT8 NOT NULL,
    goods_total FLOAT8 NOT NULL,
    custom_fee FLOAT8 NOT NULL
    );

CREATE TABLE IF NOT EXISTS "order" (
   uid UUID PRIMARY KEY,
   track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(255),
    locale VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(255),
    shard_key VARCHAR(255) NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMPTZ NOT NULL,
    oof_shard VARCHAR(255),
    delivery_id int not null,
    payment_id int not null,
    FOREIGN KEY (delivery_id) REFERENCES delivery(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (payment_id) REFERENCES payment(id) ON UPDATE CASCADE ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS item (
    id SERIAL PRIMARY KEY,
    order_uid UUID NOT NULL,
    chrt_id INT NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price FLOAT8 NOT NULL,
    rid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale FLOAT8 NOT NULL,
    size VARCHAR(50),
    total_price FLOAT8 NOT NULL,
    nm_id INT NOT NULL,
    brand VARCHAR(100),
    status INT NOT NULL,
    FOREIGN KEY (order_uid) REFERENCES "order"(uid) ON UPDATE CASCADE ON DELETE CASCADE
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "order";
DROP TABLE IF EXISTS delivery;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS item;
-- +goose StatementEnd

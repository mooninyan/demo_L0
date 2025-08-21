-- Таблица для MainOrder
CREATE TABLE main_orders
(
    order_uid          VARCHAR(255) PRIMARY KEY,
    track_number       VARCHAR(255),
    entry              VARCHAR(255),
    delivery_id        INTEGER,
    payment_id         INTEGER,
    locale             VARCHAR(50),
    internal_signature VARCHAR(255),
    customer_id        VARCHAR(255),
    delivery_service   VARCHAR(255),
    shardkey           VARCHAR(50),
    sm_id              INTEGER,
    date_created       TIMESTAMP,
    oof_shard          VARCHAR(50)
);

-- Таблица для Delivery
CREATE TABLE deliveries
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(255),
    phone   VARCHAR(50),
    zip     VARCHAR(20),
    city    VARCHAR(100),
    address VARCHAR(255),
    region  VARCHAR(100),
    email   VARCHAR(255)
);

-- Таблица для Payment
CREATE TABLE payments
(
    id            SERIAL PRIMARY KEY,
    transaction   VARCHAR(255),
    request_id    VARCHAR(255),
    currency      VARCHAR(10),
    provider      VARCHAR(50),
    amount        INTEGER,
    payment_dt    BIGINT,
    bank          VARCHAR(50),
    delivery_cost INTEGER,
    goods_total   INTEGER,
    custom_fee    INTEGER
);

-- Таблица для Item
CREATE TABLE items
(
    id            SERIAL PRIMARY KEY,
    order_uid     VARCHAR(255) NOT NULL,
    chrt_id       INTEGER,
    track_number  VARCHAR(255),
    price         INTEGER,
    rid           VARCHAR(255),
    name          VARCHAR(255),
    sale          INTEGER,
    size          VARCHAR(50),
    total_price   INTEGER,
    nm_id         INTEGER,
    brand         VARCHAR(100),
    status        INTEGER
);
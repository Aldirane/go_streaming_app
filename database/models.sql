CREATE TABLE orders (
    order_uid varchar(100) UNIQUE PRIMARY KEY,
    track_number varchar(100) UNIQUE NOT NULL,
    entry varchar(20) NOT NULL,
    locale varchar(20) NOT NULL,
    internal_signature varchar(100),
    customer_id varchar(100) NOT NULL,
    delivery_service varchar(100) NOT NULL,
    shardkey varchar(100) NOT NULL,
    sm_id integer NOT NULL,
    date_created timestamp NOT NULL,
    oof_shard varchar(10) NOT NULL
);


CREATE TABLE delivery (
    track_number varchar(100) UNIQUE PRIMARY KEY,
    name varchar(50) NOT NULL,
    phone varchar(20) NOT NULL,
    zip varchar(100) NOT NULL,
    city varchar(100) NOT NULL,
    address varchar(100) NOT NULL,
    region varchar(100) NOT NULL,
    email varchar(100) NOT NULL,
    FOREIGN KEY (track_number) REFERENCES orders (track_number)
    ON DELETE CASCADE
);


CREATE TABLE payment (
    transaction varchar(100) UNIQUE PRIMARY KEY,
    request_id varchar(100),
    currency varchar(20) NOT NULL,
    provider varchar(20) NOT NULL,
    amount integer NOT NULL,
    payment_dt integer NOT NULL,
    bank varchar(20) NOT NULL,
    delivery_cost integer NOT NULL,
    goods_total integer NOT NULL,
    custom_fee integer NOT NULL,
    FOREIGN KEY (transaction) REFERENCES orders (order_uid)
    ON DELETE CASCADE
);



CREATE TABLE item (
    chrt_id integer UNIQUE PRIMARY KEY,
    track_number varchar(100) NOT NULL,
    price integer NOT NULL,
    rid varchar(100) NOT NULL,
    name varchar(40) NOT NULL,
    sale integer NOT NULL,
    size varchar(20) NOT NULL,
    total_price integer NOT NULL,
    nm_id integer NOT NULL,
    brand varchar(20) NOT NULL,
    status integer NOT NULL,
    FOREIGN KEY (track_number) REFERENCES orders (track_number)
    ON DELETE CASCADE
);
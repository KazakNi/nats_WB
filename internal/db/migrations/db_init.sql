--drop if not exists
DROP TABLE IF EXISTS ORDERS, DELIVERIES, PAYMENTS, ITEMS;
-- create
CREATE TABLE ORDERS (
  order_uid VARCHAR (50) PRIMARY KEY NOT NULL,
  track_number VARCHAR (50),
  entry VARCHAR (50),
  locale VARCHAR (50),
  internal_signature VARCHAR (50),
  customer_id INT,
  delivery_service VARCHAR (50),
  shardkey INT,
  sm_id INT,
  date_created TIME,
  oof_shard INT
);

CREATE TABLE DELIVERIES (
  order_id VARCHAR (50) REFERENCES orders,
  name VARCHAR (50),
  phone VARCHAR (50),
  zip INT,
  city VARCHAR (50),
  address VARCHAR (50),
  region VARCHAR (50),
  email VARCHAR (50)
);

CREATE TABLE PAYMENTS (
  order_id VARCHAR (50) REFERENCES orders,
  transaction VARCHAR (50),
  request_id INT,
  currency VARCHAR (50),
  provider VARCHAR (50),
  amount INT,
  payment_dt INT,
  bank VARCHAR (50),
  delivery_cost INT,
  goods_total INT,
  custom_fee INT
);

CREATE TABLE ITEMS (
  order_id VARCHAR (50) REFERENCES orders,
  chrt_id INT,
  track_number VARCHAR (50),
  price INT,
  rid VARCHAR (50),
  name VARCHAR (50),
  sale INT,
  size INT,
  total_price INT,
  nm_id INT,
  brand VARCHAR (50),
  status INT
);


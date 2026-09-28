CREATE DATABASE IF NOT EXISTS agentgo_test CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE agentgo_test;

DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS Orders_Log;

CREATE TABLE Orders_Log (
  amount FLOAT NULL,
  created_at TIMESTAMP NULL,
  note VARCHAR(32) CHARACTER SET latin1 NULL,
  user_id BIGINT NULL,
  code VARCHAR(32) NULL,
  a INT, b INT, c INT, d INT, e INT,
  INDEX idx_code (code),
  INDEX idx_code_duplicate (code),
  INDEX idx_user_id (user_id),
  INDEX idx_wide (a, b, c, d, e)
) COMMENT='';

CREATE TABLE order_items (
  id BIGINT PRIMARY KEY,
  order_id BIGINT NOT NULL,
  CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES Orders_Log(user_id)
) COMMENT='order items';

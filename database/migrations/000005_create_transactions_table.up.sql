CREATE TABLE transactions (
  id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  invoice_code VARCHAR(20) NOT NULL UNIQUE,
  outlet_id CHAR(36) NOT NULL,
  total BIGINT NOT NULL,
  is_paid BOOLEAN NOT NULL DEFAULT FALSE,
  note TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (outlet_id) REFERENCES outlets(id) ON DELETE CASCADE
);
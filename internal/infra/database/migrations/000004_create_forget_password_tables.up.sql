CREATE TABLE reset_password_tokens(
  token varchar(21) PRIMARY KEY DEFAULT generate_nanoid(),
  user_id varchar(21) NOT NULL,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  expires_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP + interval '20 minutes',
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


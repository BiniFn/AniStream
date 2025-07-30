CREATE TABLE users(
  id varchar(21) PRIMARY KEY DEFAULT generate_nanoid(),
  username varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  password_hash varchar(255) NOT NULL,
  profile_picture text,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (username),
  UNIQUE (email)
);

CREATE TYPE Provider AS ENUM(
  'myanimelist',
  'anilist'
);

CREATE TABLE user_tokens(
  id varchar(21) PRIMARY KEY DEFAULT generate_nanoid(),
  user_id varchar(21) NOT NULL,
  token text NOT NULL,
  refresh_token text NOT NULL,
  provider Provider NOT NULL,
  expires_at timestamp NOT NULL,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE sessions(
  id varchar(21) PRIMARY KEY DEFAULT generate_nanoid(),
  user_id varchar(21) NOT NULL,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  expires_at timestamp NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


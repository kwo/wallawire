CREATE TABLE users (
  id            UUID        NOT NULL PRIMARY KEY,
  disabled      BOOL        NOT NULL,
  username      VARCHAR(64) NOT NULL CHECK (LENGTH(BTRIM(username)) > 0) UNIQUE,
  name          VARCHAR(64) NOT NULL CHECK (LENGTH(BTRIM(name)) > 0),
  password_hash CHAR(120)   NOT NULL CHECK (LENGTH(BTRIM(password_hash)) > 0),
  created       INTEGER     NOT NULL,
  updated       INTEGER     NOT NULL
);

CREATE VIEW usernames (username)
AS
SELECT LOWER(username)
FROM users;

CREATE TABLE roles (
  id   UUID        NOT NULL PRIMARY KEY,
  name VARCHAR(64) NOT NULL CHECK (LENGTH(BTRIM(name)) > 0) UNIQUE
);

CREATE TABLE user_role (
  user_id    UUID NOT NULL REFERENCES users (id),
  role_id    UUID NOT NULL REFERENCES roles (id),
  valid_from INTEGER,
  valid_to   INTEGER,
  PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idxUserRoleRole ON user_role (role_id, user_id);

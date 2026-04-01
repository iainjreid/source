CREATE TABLE IF NOT EXISTS "objects" (
  type SMALLINT,
  hash CHAR(40),
  parent_hash CHAR(40),
  cont BYTEA,
  length INTEGER
);

CREATE TABLE IF NOT EXISTS "refs" (
  type SMALLINT,
  hash CHAR(40),
  name VARCHAR,
  target VARCHAR
);

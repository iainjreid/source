-- This query creates the tables associated with a given repository.
CREATE TABLE IF NOT EXISTS "%s_objects" (
  hash CHAR(40) NOT NULL,
  parent_hash CHAR(40),
  type SMALLINT NOT NULL,
  length INTEGER NOT NULL,
  cont BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS "%s_refs" (
  hash CHAR(40) NOT NULL,
  type SMALLINT NOT NULL,
  target VARCHAR NOT NULL,
  name VARCHAR NOT NULL
);

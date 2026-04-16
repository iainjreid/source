CREATE TABLE IF NOT EXISTS "source_objects" (
  cont BYTEA NOT NULL,
  hash CHAR(40) NOT NULL,
  parent_hash CHAR(40),
  length INTEGER NOT NULL,
  type SMALLINT NOT NULL
);

CREATE TABLE IF NOT EXISTS "source_refs" (
  hash CHAR(40) NOT NULL,
  name VARCHAR NOT NULL,
  target VARCHAR NOT NULL,
  type SMALLINT NOT NULL
);

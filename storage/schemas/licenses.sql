-- The "licenses" table stores software licenses that are used to match against
-- licenses within a repository and offer metadata (TBD) about them.
CREATE TABLE licenses (
    id              CHAR(36) PRIMARY KEY,

    name            TEXT NOT NULL,
    desription      TEXT NOT NULL DEFAULT '',

    contents        TEXT NOT NULL,
);

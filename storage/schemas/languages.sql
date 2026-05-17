-- The "languages" table stores languages with their name and an associated
-- color code in HEX format.
CREATE TABLE languages (
    id              CHAR(36) PRIMARY KEY,

    name            TEXT NOT NULL UNIQUE,

    color           CHAR(7) NOT NULL
);

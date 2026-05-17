-- The "extensions" table stores file extensions used by different languages,
-- along with a foreign key reference to the respective language.
CREATE TABLE extensions (
    id              CHAR(36) PRIMARY KEY,
    language_id     CHAR(36) NOT NULL,

    extension       TEXT NOT NULL,

    CONSTRAINT fk_extension_language
        FOREIGN KEY (language_id)
        REFERENCES languages(id)
        ON DELETE RESTRICT,

    -- The mapping between the language and the extension must be unique.
    CONSTRAINT uq_language_extension
        UNIQUE (language_id, extension)
);

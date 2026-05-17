-- The "objects" table stores Git objects.
CREATE TABLE objects (
    repo_id         CHAR(36) NOT NULL,

    type            SMALLINT NOT NULL,
    hash            CHAR(40) NOT NULL,

    contents        BYTEA NOT NULL,
    length          INTEGER NOT NULL,

    PRIMARY KEY (repo_id, hash),

    CONSTRAINT fk_objects_repo
        FOREIGN KEY (repo_id)
        REFERENCES repos(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_objects_hash
    ON objects(repo_id, hash);

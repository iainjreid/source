-- The "refs" table stores Git references.
CREATE TABLE refs (
    repo_id         CHAR(36) NOT NULL,
    
    type            SMALLINT NOT NULL,
    hash            CHAR(40) NOT NULL,
    
    name            TEXT NOT NULL,
    target          TEXT,
    
    PRIMARY KEY (repo_id, name),

    CONSTRAINT fk_refs_repo
        FOREIGN KEY (repo_id)
        REFERENCES repos(id)
        ON DELETE CASCADE

    -- Add foreign key contraint for object hash
);

CREATE INDEX idx_refs_hash
    ON refs(hash);

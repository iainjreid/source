-- The "topics" table stores descriptive terms that can be applied to Git
-- repositories to act as a grouping and discoverability mechanic.
CREATE TABLE topics (
    id              CHAR(36) PRIMARY KEY,
    
    text            TEXT NOT NULL UNIQUE
);

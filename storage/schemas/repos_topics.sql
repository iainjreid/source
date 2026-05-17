-- The "repos_tags" table is a join table that associates tags with Git repositories.
CREATE TABLE repos_topics (
    repo_id         CHAR(36) NOT NULL,
    topic_id        CHAR(36) NOT NULL,

    PRIMARY KEY (repo_id, topic_id),

    FOREIGN KEY (repo_id)
        REFERENCES repos(id)
        ON DELETE CASCADE,

    FOREIGN KEY (topic_id)
        REFERENCES topics(id)
        ON DELETE RESTRICT
);

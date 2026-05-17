-- The "repos" table stores data about the Git repositories held within the
-- system.
--
-- It acts as the top-level entity for repository-related data. Some columns
-- contain normalised data, and some columns are updated by background tasks.
--
-- For the "license_id" and "language_id" columns to behave as one might expect,
-- the system must have the respective tables populated with data that can be
-- matched against.
CREATE TABLE repos (
    id              CHAR(36) PRIMARY KEY,

    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    license_id      CHAR(36),
    language_id     CHAR(36),

    default_branch  CHAR(36)
);

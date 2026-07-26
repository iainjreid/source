/*
 * This query updates an existing reference the "refs" table belonging to the
 * repository with the provided ID.
 *
 * $1 - The ID of the repository
 * $2 - The new object hash
 * $3 - The old object hash
 */
UPDATE refs
SET
    hash = $2

WHERE refs.repo_id = $1 AND refs.hash = $3;

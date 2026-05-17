/*
 * This query inserts a new reference into the "refs" table under the repository
 * identified by the provided name.
 *
 * $1 - The name of the repository
 * $2 - The reference type
 * $3 - The object hash to point to
 * $4 - The full name of the reference
 * $5 - The target reference (if symbolic)
 */
INSERT INTO refs (
    repo_id,
    type,
    hash,
    name,
    target
)

SELECT
    repos.id,
    $2, $3, $4, $5

FROM repos
WHERE repos.name = $1;

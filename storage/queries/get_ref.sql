/*
 * This query selects a reference within a repository given their respective
 * names.
 *
 * $1 - The ID of the repository to which the reference belongs
 * $2 - The name of the reference to find
 */
SELECT
    refs.type,
    refs.hash,
    refs.name,
    refs.target
FROM refs

WHERE refs.repo_id = $1
  AND refs.name = $2;

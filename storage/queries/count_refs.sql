/*
 * This query counts the number of references within a given repository.
 *
 * $1 - The ID of the repository
 */
SELECT
    COUNT(*)
FROM refs

WHERE refs.repo_id = $1;

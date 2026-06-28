/*
 * This query selects references within a repository given the respository's
 * name.
 *
 * $1 - The ID of the repository to search within
 */
 SELECT
    refs.type,
    refs.hash,
    refs.name,
    refs.target
FROM refs

WHERE refs.repo_id = $1;

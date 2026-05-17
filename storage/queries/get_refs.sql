/*
 * This query selects references within a repository given the respository's
 * name.
 *
 * $1 - The name of the repository to search within
 */
 SELECT
    refs.type,
    refs.hash,
    refs.name,
    refs.target
FROM refs

JOIN repos
  ON refs.repo_id = repos.id

WHERE repos.name = $1;

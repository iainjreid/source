/*
 * This query counts the number of references within a given repository.
 *
 * $1 - The name of the repository
 */
SELECT
    COUNT(*)
FROM refs

JOIN repos
  ON refs.repo_id = repos.id

WHERE repos.name = $1;

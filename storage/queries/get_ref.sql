/*
 * This query selects a reference within a repository given their respective
 * names.
 *
 * $1 - The name of the repository to which the reference belongs
 * $2 - The name of the reference to find
 */
SELECT
    refs.type,
    refs.hash,
    refs.name,
    refs.target
FROM refs

JOIN repos
  ON refs.repo_id = repos.id

WHERE repos.name = $1
  AND refs.name = $2;

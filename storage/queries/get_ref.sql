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
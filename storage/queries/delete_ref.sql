/*
 * This query deletes a reference from a repository given their respective
 * names.
 *
 * $1 - The ID of the repository to which the reference belongs
 * $2 - The name of the reference to delete
 */
DELETE FROM refs

WHERE refs.repo_id = $1
  AND refs.name = $2;

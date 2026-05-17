/*
 * This query deletes a reference from a repository given their respective
 * names. 
 *
 * $1 - The name of the repository to which the reference belongs
 * $2 - The name of the reference to delete
 */
DELETE FROM refs
WHERE EXISTS (
    SELECT 1
    FROM refs r

    JOIN repos
        ON r.repo_id = repos.id

    WHERE repos.name = $1
)
AND refs.name = $2;

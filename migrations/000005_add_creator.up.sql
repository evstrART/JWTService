WITH ranked_users AS (
    SELECT DISTINCT ON (ou_inner.organization_id)
    ou_inner.id
FROM organization_users ou_inner
    JOIN users_migrate u ON ou_inner.user_id = u.id
ORDER BY
    ou_inner.organization_id,
    (u.status = 'ACTIVE') DESC,
    ou_inner.date_created ASC
    )
UPDATE organization_users ou
SET type = 'creator'
    FROM ranked_users
WHERE ou.id = ranked_users.id;

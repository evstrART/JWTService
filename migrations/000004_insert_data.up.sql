INSERT INTO users_migrate (name, status) VALUES
    ('User1', 'ACTIVE'),
    ('User2', 'INACTIVE'),
    ('User3', 'ACTIVE'),
    ('User4', 'INACTIVE'),
    ('User5', 'ACTIVE'),
    ('User6', 'ACTIVE'),
    ('User7', 'INACTIVE');

INSERT INTO organization_users (user_id, organization_id, date_created) VALUES
    ((SELECT id FROM users_migrate WHERE name = 'User1'), 1, '2025-01-10 10:00:00'), -- ACTIVE
    ((SELECT id FROM users_migrate WHERE name = 'User2'), 1, '2025-01-11 10:00:00'), -- INACTIVE
    ((SELECT id FROM users_migrate WHERE name = 'User3'), 1, '2025-01-12 10:00:00'); -- ACTIVE

INSERT INTO organization_users (user_id, organization_id, date_created) VALUES
    ((SELECT id FROM users_migrate WHERE name = 'User4'), 2, '2025-02-01 09:00:00'), -- INACTIVE
    ((SELECT id FROM users_migrate WHERE name = 'User5'), 2, '2025-02-02 09:00:00'), -- ACTIVE
    ((SELECT id FROM users_migrate WHERE name = 'User6'), 2, '2025-02-03 09:00:00'); -- ACTIVE

INSERT INTO organization_users (user_id, organization_id, date_created) VALUES
    ((SELECT id FROM users_migrate WHERE name = 'User7'), 3, '2025-03-01 11:00:00'); -- INACTIVE


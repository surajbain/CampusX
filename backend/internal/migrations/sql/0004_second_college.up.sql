-- =====================================================================
-- 0004_second_college
-- Second college (IIT Delhi) + its admin, for tenant isolation tests.
--
-- Known password: Admin@123 (same as IITB admin)
-- =====================================================================

INSERT INTO colleges (id, name, slug, city, state, contact_email)
VALUES
    ('44444444-4444-4444-4444-444444444401',
     'IIT Delhi', 'iitd', 'New Delhi', 'Delhi', 'events@iitd.ac.in')
ON CONFLICT (slug) DO NOTHING;

-- Use the same Argon2id hash as admin@iitb.edu (password: Admin@123)
-- so tests can log in with a known password.
INSERT INTO users (id, email, password_hash, full_name, role, college_id, email_verified)
SELECT
    '44444444-4444-4444-4444-444444444402',
    'admin@iitd.ac.in',
    password_hash,
    'IITD Admin',
    'COLLEGE_ADMIN',
    '44444444-4444-4444-4444-444444444401',
    TRUE
FROM users
WHERE email = 'admin@iitb.edu'
ON CONFLICT (email) DO NOTHING;
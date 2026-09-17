-- =====================================================================
-- 0002_seed — Local development seed data
-- WARNING: Only safe for development. Guarded by APP_ENV check in Go.
-- =====================================================================

-- ---------- College ----------
INSERT INTO colleges (id, name, slug, city, state, contact_email)
VALUES
    ('11111111-1111-1111-1111-111111111111',
     'IIT Bombay', 'iitb', 'Mumbai', 'Maharashtra', 'events@iitb.ac.in')
ON CONFLICT (slug) DO NOTHING;

-- ---------- Users ----------
-- Passwords (Argon2id hash of the plain text below):
--   super@campusx.dev       -> SuperAdmin@123
--   admin@iitb.edu          -> Admin@123
--   organizer@iitb.edu      -> Organizer@123
--   volunteer@iitb.edu      -> Volunteer@123
--   student@iitb.edu        -> Student@123
--   advertiser@brandco.in   -> Advertiser@123
--
-- NOTE: These are placeholder bcrypt-format hashes for Phase 2 verification only.
--       Phase 3 replaces them with proper Argon2id hashes.
INSERT INTO users (id, email, password_hash, full_name, role, college_id, email_verified)
VALUES
    ('22222222-2222-2222-2222-222222222201',
     'super@campusx.dev',
     '$argon2id$v=19$m=65536,t=3,p=4$SEEDPLACEHOLDER',
     'Super Admin', 'SUPER_ADMIN', NULL, TRUE),

    ('22222222-2222-2222-2222-222222222202',
     'admin@iitb.edu',
     '$argon2id$v=19$m=65536,t=3,p=4$SEEDPLACEHOLDER',
     'IITB Admin', 'COLLEGE_ADMIN', '11111111-1111-1111-1111-111111111111', TRUE),

    ('22222222-2222-2222-2222-222222222203',
     'organizer@iitb.edu',
     '$argon2id$v=19$m=65536,t=3,p=4$SEEDPLACEHOLDER',
     'Aarav Organizer', 'ORGANIZER', '11111111-1111-1111-1111-111111111111', TRUE),

    ('22222222-2222-2222-2222-222222222204',
     'volunteer@iitb.edu',
     '$argon2id$v=19$m=65536,t=3,p=4$SEEDPLACEHOLDER',
     'Diya Volunteer', 'VOLUNTEER', '11111111-1111-1111-1111-111111111111', TRUE),

    ('22222222-2222-2222-2222-222222222205',
     'student@iitb.edu',
     '$argon2id$v=19$m=65536,t=3,p=4$SEEDPLACEHOLDER',
     'Rohan Student', 'STUDENT', '11111111-1111-1111-1111-111111111111', TRUE),

    ('22222222-2222-2222-2222-222222222206',
     'advertiser@brandco.in',
     '$argon2id$v=19$m=65536,t=3,p=4$SEEDPLACEHOLDER',
     'Brand Co', 'ADVERTISER', NULL, TRUE)
ON CONFLICT (email) DO NOTHING;

-- ---------- Events ----------
INSERT INTO events (
    id, college_id, created_by, slug, title, description, category, venue, city,
    status, starts_at, ends_at, registration_opens, registration_closes,
    price_paise, capacity, allow_teams, team_size_min, team_size_max,
    prize_pool_paise, is_featured
)
VALUES
    ('33333333-3333-3333-3333-333333333301',
     '11111111-1111-1111-1111-111111111111',
     '22222222-2222-2222-2222-222222222203',
     'techfest-2026', 'TechFest 2026',
     'Annual technology festival with hackathons, workshops, and tech talks.',
     'TECH_FEST', 'Main Auditorium', 'Mumbai', 'PUBLISHED',
     NOW() + INTERVAL '30 days', NOW() + INTERVAL '32 days',
     NOW() - INTERVAL '10 days', NOW() + INTERVAL '25 days',
     0, 1000, FALSE, NULL, NULL, 50000000, TRUE),

    ('33333333-3333-3333-3333-333333333302',
     '11111111-1111-1111-1111-111111111111',
     '22222222-2222-2222-2222-222222222203',
     'codefest-2026', 'CodeFest 2026 — 24h Hackathon',
     'Solve real-world problems in 24 hours. Prizes worth ₹5L.',
     'HACKATHON', 'Innovation Lab', 'Mumbai', 'PUBLISHED',
     NOW() + INTERVAL '15 days', NOW() + INTERVAL '16 days',
     NOW() - INTERVAL '5 days', NOW() + INTERVAL '10 days',
     9900, 200, TRUE, 2, 4, 500000, TRUE)
ON CONFLICT (college_id, slug) DO NOTHING;

-- ---------- Ad slots ----------
INSERT INTO ad_slots (code, description, max_active)
VALUES
    ('HOME_HERO',       'Hero banner on student home',       2),
    ('HOME_MID',        'Mid-feed sponsored card',           3),
    ('EVENT_SIDEBAR',   'Sidebar on event detail page',      2)
ON CONFLICT (code) DO NOTHING;
-- =====================================================================
-- 0006_registration_seed
-- Sample registrations for dev testing.
-- =====================================================================

-- Student registered (confirmed, free event)
INSERT INTO registrations (id, event_id, user_id, college_id, type, status)
VALUES
    ('77777777-7777-7777-7777-777777777701',
     '33333333-3333-3333-3333-333333333301',  -- TechFest 2026
     '22222222-2222-2222-2222-222222222205',  -- student@iitb.edu
     '11111111-1111-1111-1111-111111111111',
     'ATTEND',
     'CONFIRMED')
ON CONFLICT (event_id, user_id, type) DO NOTHING;

-- Volunteer registered for TechFest (participate)
INSERT INTO registrations (id, event_id, user_id, college_id, type, status)
VALUES
    ('77777777-7777-7777-7777-777777777702',
     '33333333-3333-3333-3333-333333333301',  -- TechFest 2026
     '22222222-2222-2222-2222-222222222204',  -- volunteer@iitb.edu
     '11111111-1111-1111-1111-111111111111',
     'PARTICIPATE',
     'CONFIRMED')
ON CONFLICT (event_id, user_id, type) DO NOTHING;
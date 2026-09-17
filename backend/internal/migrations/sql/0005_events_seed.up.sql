-- =====================================================================
-- 0005_events_seed
-- Add events to IITD + a draft to IITB for lifecycle testing.
-- =====================================================================

-- IITD PUBLISHED events
INSERT INTO events (
    id, college_id, created_by, slug, title, description, category,
    venue, city, status, starts_at, ends_at,
    registration_opens, registration_closes,
    price_paise, capacity, allow_teams, team_size_min, team_size_max,
    prize_pool_paise, is_featured
)
VALUES
    ('55555555-5555-5555-5555-555555555501',
     '44444444-4444-4444-4444-444444444401',
     '44444444-4444-4444-4444-444444444402',
     'rendezvous-2026', 'Rendezvous 2026 — IITD Cultural Fest',
     'Four days of music, dance, drama, and food.',
     'CULTURAL', 'Open Air Theatre', 'New Delhi', 'PUBLISHED',
     NOW() + INTERVAL '45 days', NOW() + INTERVAL '49 days',
     NOW() - INTERVAL '5 days', NOW() + INTERVAL '40 days',
     0, 5000, FALSE, NULL, NULL, 30000000, TRUE),

    ('55555555-5555-5555-5555-555555555502',
     '44444444-4444-4444-4444-444444444401',
     '44444444-4444-4444-4444-444444444402',
     'sports-meet-2026', 'IITD Annual Sports Meet',
     'Inter-hostel sports championship — cricket, football, athletics.',
     'SPORTS', 'Main Ground', 'New Delhi', 'PUBLISHED',
     NOW() + INTERVAL '20 days', NOW() + INTERVAL '25 days',
     NOW() - INTERVAL '10 days', NOW() + INTERVAL '15 days',
     0, 2000, FALSE, NULL, NULL, 500000, FALSE)
ON CONFLICT (college_id, slug) DO NOTHING;

-- IITB DRAFT event (for lifecycle testing)
INSERT INTO events (
    id, college_id, created_by, slug, title, description, category,
    venue, city, status, starts_at, ends_at,
    price_paise, capacity, prize_pool_paise, is_featured
)
VALUES
    ('55555555-5555-5555-5555-555555555503',
     '11111111-1111-1111-1111-111111111111',
     '22222222-2222-2222-2222-222222222203',
     'workshop-ai-2026', 'AI/ML Bootcamp — IITB',
     'Two-day intensive workshop on modern ML. Draft — not yet published.',
     'WORKSHOP', 'CS Department', 'Mumbai', 'DRAFT',
     NOW() + INTERVAL '60 days', NOW() + INTERVAL '62 days',
     4900, 100, 0, FALSE)
ON CONFLICT (college_id, slug) DO NOTHING;

-- Add `created_by` organizer for IITB second event if missing
UPDATE events
SET created_by = '22222222-2222-2222-2222-222222222203'
WHERE created_by IS NULL
  AND college_id = '11111111-1111-1111-1111-111111111111';
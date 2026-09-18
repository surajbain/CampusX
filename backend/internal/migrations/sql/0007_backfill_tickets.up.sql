-- =====================================================================
-- 0007_backfill_tickets
-- Backfill tickets for existing CONFIRMED registrations.
-- Note: This is a SQL migration, so QR payloads will be placeholder.
--       The app will re-sign these via IssueForRegistration when accessed.
-- =====================================================================

-- We insert a placeholder QR payload/signature for existing registrations.
-- The ticket service will re-sign them on first access via IssueForRegistration.
INSERT INTO tickets (id, registration_id, user_id, event_id, ticket_code, qr_payload, qr_signature, status, issued_at)
SELECT
    gen_random_uuid(),
    r.id,
    r.user_id,
    r.event_id,
    'CX-' || SUBSTRING(MD5(r.id::text), 1, 4) || '-' || SUBSTRING(MD5(r.id::text), 5, 4),
    'PENDING',
    'PENDING',
    'ACTIVE',
    NOW()
FROM registrations r
WHERE r.status = 'CONFIRMED'
  AND r.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM tickets t WHERE t.registration_id = r.id
  )
ON CONFLICT (registration_id) DO NOTHING;
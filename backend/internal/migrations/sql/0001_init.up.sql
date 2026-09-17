-- =====================================================================
-- 0001_init — CampusX core schema
-- =====================================================================

-- Enable UUID generation (pgcrypto for gen_random_uuid).
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ---------------------------------------------------------------------
-- colleges
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS colleges (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(200) NOT NULL,
    slug            VARCHAR(120) NOT NULL UNIQUE,
    city            VARCHAR(120) NOT NULL,
    state           VARCHAR(120),
    logo_url        TEXT,
    website         TEXT,
    contact_email   VARCHAR(200),
    contact_phone   VARCHAR(30),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_colleges_deleted_at ON colleges (deleted_at);
CREATE INDEX IF NOT EXISTS idx_colleges_city ON colleges (city);

-- ---------------------------------------------------------------------
-- users
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    full_name       VARCHAR(200) NOT NULL,
    phone           VARCHAR(30),
    role            VARCHAR(30) NOT NULL,   -- STUDENT | COLLEGE_ADMIN | ORGANIZER | VOLUNTEER | SUPER_ADMIN | ADVERTISER
    college_id      UUID REFERENCES colleges(id) ON DELETE SET NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    email_verified  BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);
CREATE INDEX IF NOT EXISTS idx_users_college_id ON users (college_id);

-- ---------------------------------------------------------------------
-- events
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS events (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    college_id          UUID NOT NULL REFERENCES colleges(id) ON DELETE CASCADE,
    created_by          UUID REFERENCES users(id) ON DELETE SET NULL,
    slug                VARCHAR(160) NOT NULL,
    title               VARCHAR(250) NOT NULL,
    description         TEXT,
    category            VARCHAR(40) NOT NULL,   -- HACKATHON | CULTURAL | SPORTS | WORKSHOP | TECH_FEST | OTHER
    poster_url          TEXT,
    venue               VARCHAR(250),
    city                VARCHAR(120),

    -- Lifecycle
    status              VARCHAR(30) NOT NULL DEFAULT 'DRAFT', -- DRAFT | PUBLISHED | ONGOING | COMPLETED | CANCELLED

    -- Dates
    starts_at           TIMESTAMPTZ,
    ends_at             TIMESTAMPTZ,
    registration_opens  TIMESTAMPTZ,
    registration_closes TIMESTAMPTZ,

    -- Participation model
    price_paise         BIGINT NOT NULL DEFAULT 0,      -- 0 = free. Store in paise to avoid float.
    currency            VARCHAR(8) NOT NULL DEFAULT 'INR',
    capacity            INTEGER,                        -- NULL = unlimited
    allow_teams         BOOLEAN NOT NULL DEFAULT FALSE,
    team_size_min       INTEGER,
    team_size_max       INTEGER,

    -- Metadata
    prize_pool_paise    BIGINT NOT NULL DEFAULT 0,
    contact_email       VARCHAR(200),
    whatsapp_link       TEXT,
    rules               TEXT,
    schedule_json       JSONB,
    faq_json            JSONB,

    is_featured         BOOLEAN NOT NULL DEFAULT FALSE,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    UNIQUE (college_id, slug)
);
CREATE INDEX IF NOT EXISTS idx_events_deleted_at ON events (deleted_at);
CREATE INDEX IF NOT EXISTS idx_events_college_id ON events (college_id);
CREATE INDEX IF NOT EXISTS idx_events_status ON events (status);
CREATE INDEX IF NOT EXISTS idx_events_category ON events (category);
CREATE INDEX IF NOT EXISTS idx_events_starts_at ON events (starts_at);
CREATE INDEX IF NOT EXISTS idx_events_city ON events (city);

-- ---------------------------------------------------------------------
-- registrations
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS registrations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    college_id      UUID NOT NULL REFERENCES colleges(id) ON DELETE CASCADE,
    type            VARCHAR(20) NOT NULL,       -- ATTEND | PARTICIPATE
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING | CONFIRMED | CANCELLED | WAITLISTED
    team_id         UUID,                       -- FK added after teams table
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    -- Anti-fraud: one user can register once per event per type
    UNIQUE (event_id, user_id, type)
);
CREATE INDEX IF NOT EXISTS idx_registrations_deleted_at ON registrations (deleted_at);
CREATE INDEX IF NOT EXISTS idx_registrations_event_id ON registrations (event_id);
CREATE INDEX IF NOT EXISTS idx_registrations_user_id ON registrations (user_id);
CREATE INDEX IF NOT EXISTS idx_registrations_status ON registrations (status);

-- ---------------------------------------------------------------------
-- payments
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS payments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_id     UUID NOT NULL REFERENCES registrations(id) ON DELETE CASCADE,
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id            UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    provider            VARCHAR(30) NOT NULL DEFAULT 'MOCK', -- RAZORPAY | MOCK

    -- Anti-fraud
    gateway_payment_id  VARCHAR(120) UNIQUE,   -- from Razorpay; NULL for pending
    gateway_order_id    VARCHAR(120),
    idempotency_key     VARCHAR(120) NOT NULL UNIQUE,

    amount_paise        BIGINT NOT NULL,
    currency            VARCHAR(8) NOT NULL DEFAULT 'INR',
    status              VARCHAR(20) NOT NULL DEFAULT 'CREATED', -- CREATED | AUTHORIZED | CAPTURED | FAILED | REFUNDED

    signature           TEXT,
    raw_payload         JSONB,

    verified_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_payments_deleted_at ON payments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_payments_registration_id ON payments (registration_id);
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments (user_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments (status);
CREATE INDEX IF NOT EXISTS idx_payments_gateway_order_id ON payments (gateway_order_id);

-- ---------------------------------------------------------------------
-- tickets
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS tickets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_id UUID NOT NULL UNIQUE REFERENCES registrations(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id        UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,

    ticket_code     VARCHAR(40) NOT NULL UNIQUE,   -- human readable
    qr_payload      TEXT NOT NULL,                 -- HMAC-signed payload
    qr_signature    TEXT NOT NULL,
    valid_from      TIMESTAMPTZ,
    valid_until     TIMESTAMPTZ,

    status          VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE | USED | REVOKED | EXPIRED
    issued_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    used_at         TIMESTAMPTZ,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_tickets_deleted_at ON tickets (deleted_at);
CREATE INDEX IF NOT EXISTS idx_tickets_user_id ON tickets (user_id);
CREATE INDEX IF NOT EXISTS idx_tickets_event_id ON tickets (event_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets (status);

-- ---------------------------------------------------------------------
-- checkins
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS checkins (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id       UUID NOT NULL UNIQUE REFERENCES tickets(id) ON DELETE CASCADE,
    scanner_id      UUID REFERENCES users(id) ON DELETE SET NULL,   -- volunteer who scanned
    event_id        UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    college_id      UUID NOT NULL REFERENCES colleges(id) ON DELETE CASCADE,
    result          VARCHAR(20) NOT NULL,   -- SUCCESS | DUPLICATE | INVALID | EXPIRED
    reason          TEXT,
    scanned_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_checkins_event_id ON checkins (event_id);
CREATE INDEX IF NOT EXISTS idx_checkins_scanner_id ON checkins (scanner_id);
CREATE INDEX IF NOT EXISTS idx_checkins_scanned_at ON checkins (scanned_at);

-- ---------------------------------------------------------------------
-- teams
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS teams (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    college_id      UUID NOT NULL REFERENCES colleges(id) ON DELETE CASCADE,
    name            VARCHAR(150) NOT NULL,
    code            VARCHAR(20) NOT NULL,       -- join code
    leader_id       UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    is_locked       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    UNIQUE (event_id, name),
    UNIQUE (event_id, code)
);
CREATE INDEX IF NOT EXISTS idx_teams_deleted_at ON teams (deleted_at);
CREATE INDEX IF NOT EXISTS idx_teams_event_id ON teams (event_id);
CREATE INDEX IF NOT EXISTS idx_teams_leader_id ON teams (leader_id);

-- Add FK for registrations.team_id
ALTER TABLE registrations
    DROP CONSTRAINT IF EXISTS fk_registrations_team;
ALTER TABLE registrations
    ADD CONSTRAINT fk_registrations_team
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE SET NULL;

-- ---------------------------------------------------------------------
-- team_members
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS team_members (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id         UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role            VARCHAR(20) NOT NULL DEFAULT 'MEMBER',   -- LEADER | MEMBER
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (team_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_team_members_team_id ON team_members (team_id);
CREATE INDEX IF NOT EXISTS idx_team_members_user_id ON team_members (user_id);

-- ---------------------------------------------------------------------
-- results
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS results (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    college_id      UUID NOT NULL REFERENCES colleges(id) ON DELETE CASCADE,
    team_id         UUID REFERENCES teams(id) ON DELETE SET NULL,
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    position        INTEGER NOT NULL,       -- 1, 2, 3 ...
    title           VARCHAR(200),
    score           NUMERIC(10, 2),
    remarks         TEXT,
    published_by    UUID REFERENCES users(id) ON DELETE SET NULL,
    published_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (event_id, position)
);
CREATE INDEX IF NOT EXISTS idx_results_event_id ON results (event_id);
CREATE INDEX IF NOT EXISTS idx_results_team_id ON results (team_id);
CREATE INDEX IF NOT EXISTS idx_results_user_id ON results (user_id);

-- ---------------------------------------------------------------------
-- winners  (materialized view of results with certificate flag)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS winners (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    result_id       UUID NOT NULL REFERENCES results(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    position        INTEGER NOT NULL,
    certificate_url TEXT,
    issued_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (event_id, user_id, position)
);
CREATE INDEX IF NOT EXISTS idx_winners_event_id ON winners (event_id);
CREATE INDEX IF NOT EXISTS idx_winners_user_id ON winners (user_id);

-- ---------------------------------------------------------------------
-- certificates
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS certificates (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id        UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    college_id      UUID NOT NULL REFERENCES colleges(id) ON DELETE CASCADE,
    kind            VARCHAR(30) NOT NULL,   -- PARTICIPATION | WINNER | RUNNER_UP | SPECIAL
    file_url        TEXT NOT NULL,
    serial          VARCHAR(60) NOT NULL UNIQUE,
    issued_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, event_id, kind)
);
CREATE INDEX IF NOT EXISTS idx_certificates_user_id ON certificates (user_id);
CREATE INDEX IF NOT EXISTS idx_certificates_event_id ON certificates (event_id);

-- ---------------------------------------------------------------------
-- ads
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ads (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    advertiser_id   UUID REFERENCES users(id) ON DELETE SET NULL,
    college_id      UUID REFERENCES colleges(id) ON DELETE CASCADE,   -- NULL = global
    title           VARCHAR(200) NOT NULL,
    body            TEXT,
    image_url       TEXT,
    target_url      TEXT NOT NULL,
    slot_type       VARCHAR(30) NOT NULL,   -- HOME_BANNER | SIDEBAR | EVENT_FOOTER
    starts_at       TIMESTAMPTZ,
    ends_at         TIMESTAMPTZ,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    daily_budget_paise BIGINT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_ads_deleted_at ON ads (deleted_at);
CREATE INDEX IF NOT EXISTS idx_ads_slot_type ON ads (slot_type);
CREATE INDEX IF NOT EXISTS idx_ads_is_active ON ads (is_active);

-- ---------------------------------------------------------------------
-- ad_slots
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ad_slots (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(40) NOT NULL UNIQUE,   -- HOME_HERO | HOME_MID | EVENT_SIDEBAR
    description     VARCHAR(200),
    max_active      INTEGER NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ---------------------------------------------------------------------
-- ad_impressions
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ad_impressions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ad_id           UUID NOT NULL REFERENCES ads(id) ON DELETE CASCADE,
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    slot_code       VARCHAR(40) NOT NULL,
    clicked         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_ad_impressions_ad_id ON ad_impressions (ad_id);
CREATE INDEX IF NOT EXISTS idx_ad_impressions_created_at ON ad_impressions (created_at);

-- ---------------------------------------------------------------------
-- audit_logs  (fraud detection + compliance)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id        UUID REFERENCES users(id) ON DELETE SET NULL,
    college_id      UUID REFERENCES colleges(id) ON DELETE SET NULL,
    action          VARCHAR(80) NOT NULL,       -- EVENT_PUBLISHED | PAYMENT_VERIFIED | CHECKIN_REJECTED ...
    entity_type     VARCHAR(60),
    entity_id       UUID,
    severity        VARCHAR(20) NOT NULL DEFAULT 'INFO',   -- INFO | WARN | CRITICAL
    metadata        JSONB,
    ip_address      VARCHAR(45),
    user_agent      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_logs_actor_id ON audit_logs (actor_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs (action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_severity ON audit_logs (severity);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs (created_at);
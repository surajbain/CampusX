-- =====================================================================
-- 0003_real_password_hashes
-- Replace placeholder hashes with real Argon2id hashes.
--
-- Known passwords (DEV ONLY — rotate before production):
--   super@campusx.dev       -> SuperAdmin@123
--   admin@iitb.edu          -> Admin@123
--   organizer@iitb.edu      -> Organizer@123
--   volunteer@iitb.edu      -> Volunteer@123
--   student@iitb.edu        -> Student@123
--   advertiser@brandco.in   -> Advertiser@123
-- =====================================================================

UPDATE users SET password_hash =
  '$argon2id$v=19$m=65536,t=3,p=4$W1bp9bdxzfIQZPpwngnOtA$2OnheF41fxbdHbv+xdspJodk7BV3fUB7frMPJKC6aWM'
WHERE email = 'super@campusx.dev';

UPDATE users SET password_hash =
  '$argon2id$v=19$m=65536,t=3,p=4$dHovjjOLPx3UdfW/7cL2qA$f0Si+P/DnA35+fxrurN7y0ednVdrDFXTo0PQ71PhzEw'
WHERE email = 'admin@iitb.edu';

UPDATE users SET password_hash =
  '$argon2id$v=19$m=65536,t=3,p=4$DOqXLk+S2MEw4+FSKkyASw$PyIyJGdxz69OlLd/k0g4okNbR7vHf/w6RXS03i6Jlbs'
WHERE email = 'organizer@iitb.edu';

UPDATE users SET password_hash =
  '$argon2id$v=19$m=65536,t=3,p=4$lh3piQvCIi1eC09OtRIPvw$v4tK5sBY8Pk3r1qZhBcQwgYpzi66eYa3WR86CnlaVPc'
WHERE email = 'volunteer@iitb.edu';

UPDATE users SET password_hash =
  '$argon2id$v=19$m=65536,t=3,p=4$/lHf86fRtYrCRnVzwNn2uA$IKMIJ9cUZ3fHjZPPlKbToSSJzn76YXfkYG6JsUbPiwU'
WHERE email = 'student@iitb.edu';

UPDATE users SET password_hash =
  '$argon2id$v=19$m=65536,t=3,p=4$J1MSCzNraZiK2ahvWA20Yw$qAKxdL4iPJge2wr6McRk/dlCnh+iV3ZkswVXSIkYARs'
WHERE email = 'advertiser@brandco.in';
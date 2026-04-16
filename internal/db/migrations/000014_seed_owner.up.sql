-- Seed owner user: email=owner@binhvuong.vn password=BinhVuong2026!
INSERT INTO users (email, password_hash, full_name, role, status)
VALUES (
    'owner@binhvuong.vn',
    '$2b$12$jL4eHmMtHDs.Q/HaIQN9rO22yrCuDf5CKSuaDyI9Gba10I0e1Bn7G',
    'Bình Vương',
    'owner',
    'active'
) ON CONFLICT (email) DO NOTHING;

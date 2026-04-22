DROP INDEX IF EXISTS idx_notifications_dedup_per_day;
ALTER TABLE notifications
    DROP COLUMN IF EXISTS ref_type,
    DROP COLUMN IF EXISTS ref_id,
    DROP COLUMN IF EXISTS notif_date;

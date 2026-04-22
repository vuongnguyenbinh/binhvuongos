-- Enable (ref_type, ref_id, per-day) dedup so background notifiers can
-- safely run every restart + every 24h without creating duplicate rows.
ALTER TABLE notifications
    ADD COLUMN ref_type   VARCHAR(30),
    ADD COLUMN ref_id     UUID,
    ADD COLUMN notif_date DATE NOT NULL DEFAULT CURRENT_DATE;

CREATE UNIQUE INDEX idx_notifications_dedup_per_day
    ON notifications(user_id, ref_type, ref_id, notif_date)
    WHERE ref_type IS NOT NULL;

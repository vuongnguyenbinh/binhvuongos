package handler

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// deadlineWarningDays is the window during which a company end_date triggers a warning.
// Kept as a package-level constant so the notifier goroutine can reuse the same threshold.
const deadlineWarningDays = 10

// deadlineBadge returns (label, tailwind class) for a company's end_date.
// Returns empty strings when the date is nil or far in the future.
func deadlineBadge(endDate pgtype.Date) (label string, class string) {
	if !endDate.Valid {
		return "", ""
	}
	today := time.Now().Truncate(24 * time.Hour)
	deadline := endDate.Time.Truncate(24 * time.Hour)
	days := int(deadline.Sub(today).Hours() / 24)

	switch {
	case days < 0:
		return fmt.Sprintf("⚠️ Quá hạn %d ngày", -days), "bg-rust/15 text-rust"
	case days == 0:
		return "⚠️ Hết hạn hôm nay", "bg-rust/15 text-rust"
	case days <= deadlineWarningDays:
		return fmt.Sprintf("⏰ Còn %d ngày", days), "bg-ember/15 text-ember"
	default:
		return "", ""
	}
}

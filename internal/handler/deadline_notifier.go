package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
)

// deadlineCheckInterval controls how often the notifier scans for companies nearing their end_date.
// Paired with the per-day unique index on `notifications`, so running every 24h is idempotent.
const deadlineCheckInterval = 24 * time.Hour

// StartDeadlineNotifier spawns a background goroutine that creates in-app
// notifications when a company's end_date is within deadlineWarningDays.
//
// Runs once on startup, then every 24h. Panics in the loop are recovered so a single
// bad cycle cannot crash the server.
func StartDeadlineNotifier(q *generated.Queries) {
	log.Printf("deadline notifier: scheduled (initial run in 30s, then every %v)", deadlineCheckInterval)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("deadline notifier panic: %v", r)
			}
		}()
		// Slight delay so the notifier doesn't slam DB during boot while migrations run.
		time.Sleep(30 * time.Second)
		log.Printf("deadline notifier: first run starting")
		runDeadlineCheck(q)
		ticker := time.NewTicker(deadlineCheckInterval)
		defer ticker.Stop()
		for range ticker.C {
			runDeadlineCheck(q)
		}
	}()
}

// runDeadlineCheck finds companies due within deadlineWarningDays and inserts
// de-duplicated notifications for every assigned user + every owner.
func runDeadlineCheck(q *generated.Queries) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	companies, err := q.ListCompaniesDueSoon(ctx, deadlineWarningDays)
	if err != nil {
		log.Printf("deadline notifier list: %v", err)
		return
	}
	if len(companies) == 0 {
		return
	}
	owners, err := q.ListUsersByRole(ctx, "owner")
	if err != nil {
		log.Printf("deadline notifier owners: %v", err)
	}

	total := 0
	for _, co := range companies {
		label, _ := deadlineBadge(co.EndDate)
		title := "Công ty sắp hết hạn: " + co.Name
		bodyText := label
		link := "/companies/" + middleware.UUIDToString(co.ID)

		// Owners always get notified even if not directly assigned.
		for _, u := range owners {
			if err := q.CreateNotificationRef(ctx, generated.CreateNotificationRefParams{
				UserID:  u.ID,
				Title:   title,
				Body:    &bodyText,
				Link:    &link,
				RefType: "company_deadline",
				RefID:   co.ID,
			}); err == nil {
				total++
			}
		}
		// Every assigned user for this company.
		assignees, err := q.ListAssignmentsByCompany(ctx, co.ID)
		if err != nil {
			continue
		}
		for _, a := range assignees {
			if err := q.CreateNotificationRef(ctx, generated.CreateNotificationRefParams{
				UserID:  a.UserID,
				Title:   title,
				Body:    &bodyText,
				Link:    &link,
				RefType: "company_deadline",
				RefID:   co.ID,
			}); err == nil {
				total++
			}
		}
	}
	log.Printf("deadline notifier: %d companies scanned, %d notifications attempted", len(companies), total)
	_ = fmt.Sprintf // quiet static analyzers if unused
}

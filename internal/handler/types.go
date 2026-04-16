package handler

import (
	"fmt"
	"time"
)

// DashboardData — data cho trang tổng quan
type DashboardData struct {
	UserName         string
	PendingReviews   int64
	ContentReview    int64
	OverdueTasks     int64
	RawInbox         int64
	OpenTasks        int64
	DoneTasks        int64
	RunningCampaigns int64
	Companies        []CompanyView
	TodayTasks       []TaskView
}

// CompanyView — data hiển thị 1 company trong list/card
type CompanyView struct {
	ID        string
	Name      string
	ShortCode string
	Industry  string
	MyRole    string
	Status    string
	Health    string
	Scope     string // joined string from scope array
	StartDate string
}

// CompanyListData — data cho trang danh sách công ty
type CompanyListData struct {
	Companies []CompanyView
	Total     int64
}

// InboxItemView — data hiển thị 1 inbox item
type InboxItemView struct {
	ID        string
	Content   string
	URL       string
	Source    string
	ItemType  string
	Status    string
	CreatedAt string // formatted time
	TimeAgo   string // "2 giờ trước"
}

// InboxListData — data cho trang hộp thư đến
type InboxListData struct {
	Items     []InboxItemView
	RawCount  int64
	Total     int64
	Companies []CompanyView // for triage panel dropdown
}

// TaskView — data hiển thị 1 task
type TaskView struct {
	ID          string
	Title       string
	Description string
	Category    string
	GroupName   string
	CompanyCode string
	CompanyName string
	AssigneeName string
	AssigneeInit string
	Status      string
	Priority    string
	DueDate     string
}

// TaskListData — data cho trang kanban
type TaskListData struct {
	Todo       []TaskView
	InProgress []TaskView
	Waiting    []TaskView
	Review     []TaskView
	Done       []TaskView
	StatusCounts map[string]int64
}

// BookmarkView — data hiển thị 1 bookmark
type BookmarkView struct {
	ID          string
	Title       string
	URL         string
	Description string
	Tags        []string
	Notes       string
	CreatedAt   string
}

// BookmarkListData — data cho trang bookmarks
type BookmarkListData struct {
	Bookmarks []BookmarkView
	Total     int64
}

// ContentView — data hiển thị 1 content item
type ContentView struct {
	ID          string
	Title       string
	ContentType string
	Platforms   string
	CompanyCode string
	CompanyName string
	AuthorName  string
	Status      string
	PublishDate string
	Engagement  string
}

// WorkLogView — data hiển thị 1 work log
type WorkLogView struct {
	ID           string
	WorkDate     string
	UserName     string
	CompanyCode  string
	CompanyName  string
	CampaignName string
	WorkTypeName string
	WorkTypeIcon string
	Quantity     string
	Unit         string
	Status       string
	Notes        string
	SheetURL     string
	EvidenceURL  string
}

// CampaignView — data hiển thị 1 campaign
type CampaignView struct {
	ID          string
	Name        string
	CompanyCode string
	CompanyName string
	Status      string
	StartDate   string
	EndDate     string
	ProgressPct int
}

// KnowledgeView — data hiển thị 1 knowledge item
type KnowledgeView struct {
	ID          string
	Title       string
	Description string
	Category    string
	Topics      []string
	Quality     int
	Scope       string
	Format      string
	SourceURL   string
	CreatedAt   string
}

// Helper: format time to Vietnamese-friendly string
func formatTime(t time.Time) string {
	return t.Format("02/01/2006")
}

// Helper: format time ago in Vietnamese
func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "vừa xong"
	case d < time.Hour:
		return formatDuration(int(d.Minutes()), "phút")
	case d < 24*time.Hour:
		return formatDuration(int(d.Hours()), "giờ")
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "hôm qua"
		}
		return formatDuration(days, "ngày")
	}
}

func formatDuration(n int, unit string) string {
	return fmt.Sprintf("%d %s trước", n, unit)
}

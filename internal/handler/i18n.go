package handler

// Vietnamese label maps for all status/priority/category values

var taskStatusVi = map[string]string{
	"todo":        "Cần làm",
	"in_progress": "Đang làm",
	"waiting":     "Chờ",
	"review":      "Cần duyệt",
	"done":        "Hoàn thành",
	"cancelled":   "Đã huỷ",
}

var priorityVi = map[string]string{
	"urgent": "Gấp",
	"high":   "Cao",
	"normal": "Trung bình",
	"low":    "Thấp",
}

var contentStatusVi = map[string]string{
	"idea":      "Ý tưởng",
	"drafting":  "Đang viết",
	"review":    "Cần duyệt",
	"revise":    "Sửa lại",
	"approved":  "Đã duyệt",
	"published": "Đã đăng",
	"killed":    "Đã huỷ",
}

var workLogStatusVi = map[string]string{
	"submitted": "Chờ duyệt",
	"approved":  "Đã duyệt",
	"rejected":  "Từ chối",
	"needs_fix": "Cần sửa",
}

var campaignStatusVi = map[string]string{
	"planning":  "Lên kế hoạch",
	"running":   "Đang chạy",
	"paused":    "Tạm dừng",
	"ended":     "Kết thúc",
	"cancelled": "Đã huỷ",
}

var companyStatusVi = map[string]string{
	"active": "Hoạt động",
	"paused": "Tạm dừng",
	"ended":  "Kết thúc",
}

var companyHealthVi = map[string]string{
	"ok":        "Ổn",
	"attention": "Chú ý",
	"urgent":    "Gấp",
}

var roleVi = map[string]string{
	"owner":        "Chủ sở hữu",
	"core_staff":   "Nhân viên",
	"ctv":          "Cộng tác viên",
	"client_staff": "NV khách hàng",
	"intern":       "Thực tập",
}

var inboxSourceVi = map[string]string{
	"telegram": "Telegram",
	"zalo":     "Zalo",
	"web":      "Web",
	"email":    "Email",
	"manual":   "Thủ công",
	"tiktok":   "TikTok",
	"facebook": "Facebook",
}

// LabelVi returns Vietnamese label for a value from any category
func LabelVi(category, value string) string {
	var m map[string]string
	switch category {
	case "task_status":
		m = taskStatusVi
	case "priority":
		m = priorityVi
	case "content_status":
		m = contentStatusVi
	case "worklog_status":
		m = workLogStatusVi
	case "campaign_status":
		m = campaignStatusVi
	case "company_status":
		m = companyStatusVi
	case "company_health":
		m = companyHealthVi
	case "role":
		m = roleVi
	case "inbox_source":
		m = inboxSourceVi
	}
	if m != nil {
		if v, ok := m[value]; ok {
			return v
		}
	}
	return value
}

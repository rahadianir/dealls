package attendance

type AttendanceRequest struct {
	UserID string `json:"user_id"`
}

type OvertimeRequest struct {
	Hours int
}

type ReimbursementRequest struct {
	Amount      float64
	Description string
}

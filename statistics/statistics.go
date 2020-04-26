package statistics

// All counts for user
type Statistics struct {
	CategoryCount int `json:"categoryCount"`
	DataCount     int `json:"dataCount"`
	MessageCount  int `json:"messageCount"`
	ReminderCount int `json:"reminderCount"`
}

func NewStatistics() *Statistics {
	return &Statistics{}
}

// system statistics
type System struct {
	UserCount int     `json:"userCount"`
	DbSize    float64 `json:"dbSize"`
}

package main

type TimeBlock struct {
	Day 	 string `json:"day"` // Day of the week, e.g., "Monday", "Tuesday", etc.
	StartTime string `json:"start_time"` // Time in 24-hour format, e.g., "09:00"
	EndTime   string `json:"end_time"` // Time in 24-hour format, e.g., "17:00"
}

type Schedule struct {
	Blocks []TimeBlock `json:"blocks"` // List of time blocks for the week
	DailyTotal map[string]float64 `json:"daily_total"` // Map of day to total hours for that day
	WeeklyTotal float64 `json:"weekly_total"` // Total hours for the entire week
	Status string `json:"status"` // Status can be "draft", "pending", "approved", or "rejected"
	AdminComment string `json:"admin_comment,omitempty"` // Optional comment from admin when rejecting
	UserName string `json:"user_name"` // Name of the user who submitted the schedule
}

type User struct {
	UserName string `json:"user_name"` // Name of the user who submitted the schedule
	PasswordHash string `json:"password_hash"` // Hashed password for authentication
	Role string `json:"role"` // Role can be "employee" or "admin"
	ProfilePicture string `json:"profile_picture,omitempty"` // Optional URL or path to profile picture
}

type AppState struct {
	Schedule map[string]*Schedule `json:"schedule"` // Map of user name to their schedule\
	Users map[string]*User `json:"users"` // Map of user name to user details
}
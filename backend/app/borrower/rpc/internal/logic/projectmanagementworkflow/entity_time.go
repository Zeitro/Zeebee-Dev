package projectmanagementworkflowlogic

import "time"

func formatEntityTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339)
}

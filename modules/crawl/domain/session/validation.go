package session

func ValidateStatus(status string) bool {
	switch status {
	case "running", "ok", "failed":
		return true
	default:
		return false
	}
}

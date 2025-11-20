package utils

import (
	"fmt"
	"regexp"
)

// ValidateJID validates a WhatsApp JID format (phone@s.whatsapp.net)
func ValidateJID(jid string) bool {
	// Simple JID validation: numeric@s.whatsapp.net
	re := regexp.MustCompile(`^\d+@s\.whatsapp\.net$`)
	return re.MatchString(jid)
}

// TruncateString truncates a string to max length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// SafeStringJoin joins strings safely, returning a default if the result is empty
func SafeStringJoin(strs []string, sep, defaultVal string) string {
	if len(strs) == 0 {
		return defaultVal
	}
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	if result == "" {
		return defaultVal
	}
	return result
}

// FormatFileSize formats bytes to human-readable file size
func FormatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	if bytes < KB {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < MB {
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	} else if bytes < GB {
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	}
	return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
}

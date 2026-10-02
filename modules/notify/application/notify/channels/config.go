package channels

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func cfgString(cfg map[string]any, key, envKey string) string {
	if cfg != nil {
		if v, ok := cfg[key].(string); ok && strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	if envKey != "" {
		return strings.TrimSpace(os.Getenv(envKey))
	}
	return ""
}

func cfgInt(cfg map[string]any, key string) int {
	if cfg == nil {
		return 0
	}
	switch v := cfg[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		n, _ := strconv.Atoi(v)
		return n
	}
	return 0
}

func cfgBool(cfg map[string]any, key string) (bool, bool) {
	if cfg == nil {
		return false, false
	}
	v, ok := cfg[key]
	if !ok {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

func envBool(key string, fallback bool) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch v {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func parseHHMM(s string) (hour, minute int, ok bool) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	h, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, 0, false
	}
	return h, m, true
}

func InPushWindow(now time.Time, cfg map[string]any) bool {
	enabled := envBool("NOTIFY_PUSH_WINDOW_ENABLED", false)
	if v, ok := cfgBool(cfg, "push_window_enabled"); ok {
		enabled = v
	}
	if !enabled {
		return true
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	now = now.In(loc)
	start := cfgString(cfg, "push_window_start", "NOTIFY_PUSH_WINDOW_START")
	end := cfgString(cfg, "push_window_end", "NOTIFY_PUSH_WINDOW_END")
	if start == "" {
		start = "20:00"
	}
	if end == "" {
		end = "22:00"
	}
	sh, sm, ok1 := parseHHMM(start)
	eh, em, ok2 := parseHHMM(end)
	if !ok1 || !ok2 {
		return true
	}
	cur := now.Hour()*60 + now.Minute()
	from := sh*60 + sm
	to := eh*60 + em
	if from <= to {
		return cur >= from && cur <= to
	}
	return cur >= from || cur <= to
}

func OncePerDay(cfg map[string]any) bool {
	if v, ok := cfgBool(cfg, "once_per_day"); ok {
		return v
	}
	return envBool("NOTIFY_PUSH_WINDOW_ONCE_PER_DAY", true) && envBool("NOTIFY_PUSH_WINDOW_ENABLED", false)
}

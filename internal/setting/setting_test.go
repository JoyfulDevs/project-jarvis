package setting_test

import (
	"encoding/json"
	"testing"

	"github.com/joyfuldevs/project-jarvis/internal/setting"
)

func TestUnmarshal(t *testing.T) {
	raw := `[{
        "id": "C087PGPV6DN",
        "created_at": "2026-01-08T12:00:51+00:00",
        "updated_at": "2026-01-08T12:00:53+00:00",
        "scrum_notification_enabled": false,
        "scrum_notification_time": "09:00:00"
    }]`

	var settings []setting.JarvisSetting
	if err := json.Unmarshal([]byte(raw), &settings); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(settings) != 1 {
		t.Fatalf("expected 1 setting, got %d", len(settings))
	}
}

package supabase_test

import (
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joyfuldevs/project-jarvis/internal/setting"
	"github.com/joyfuldevs/project-jarvis/internal/supabase"
)

func TestSupabase(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "on" {
		t.Skip("skipping integration test")
	}

	domain := os.Getenv("SUPABASE_DOMAIN")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	authKey := os.Getenv("SUPABASE_AUTH_KEY")

	id := "TESTID"
	c := supabase.Client{
		Domain:     domain,
		APIKey:     apiKey,
		AuthKey:    authKey,
		HTTPClient: http.DefaultClient,
	}

	localSetting := &setting.JarvisSetting{
		ID:                       id,
		CreatedAt:                time.Now().In(time.FixedZone("", 0)),
		UpdatedAt:                time.Now().In(time.FixedZone("", 0)),
		ScrumNotificationEnabled: true,
		ScrumNotificationTime:    "09:30:00",
	}

	err := c.SetJarvisSetting(t.Context(), localSetting)
	if err != nil {
		t.Fatalf("failed to set jarvis setting: %v", err)
	}

	remoteSetting, err := c.GetJarvisSetting(t.Context(), id)
	if err != nil {
		t.Fatalf("failed to get jarvis setting: %v", err)
	}

	if !reflect.DeepEqual(localSetting, remoteSetting) {
		t.Fatalf("expected %+v, got %+v", localSetting, remoteSetting)
	}
}

package config

import (
	"os"
	"reflect"
	"testing"
)

func TestSplitAndTrim(t *testing.T) {
	in := "a, b, ,c"
	want := []string{"a", "b", "c"}
	if got := splitAndTrim(in); !reflect.DeepEqual(got, want) {
		t.Fatalf("splitAndTrim(%q) = %v, want %v", in, got, want)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	key := "TEST_GETENV"
	_ = os.Unsetenv(key)
	if got := getEnvOrDefault(key, "fallback"); got != "fallback" {
		t.Fatalf("expected fallback, got %q", got)
	}
	_ = os.Setenv(key, "value")
	if got := getEnvOrDefault(key, "fallback"); got != "value" {
		t.Fatalf("expected value, got %q", got)
	}
}

func TestLoadValidation(t *testing.T) {

	os.Clearenv()
	_, err := Load()
	if err == nil {
		t.Fatalf("expected error when no providers configured")
	}
}

func TestLoadDockerEventFilters(t *testing.T) {
	os.Clearenv()
	// Set a provider so validation passes
	os.Setenv("SLACK_BOT_TOKEN", "token")
	os.Setenv("SLACK_CHANNEL_IDS", "channel")

	expectedFilters := "event=start,type=container"
	os.Setenv("DOCKER_EVENT_FILTERS", expectedFilters)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"event=start", "type=container"}
	if !reflect.DeepEqual(cfg.DockerFilters, want) {
		t.Errorf("got filters %v, want %v", cfg.DockerFilters, want)
	}
}

func TestLoadTeamsWebhooks(t *testing.T) {
	os.Clearenv()
	os.Setenv("TEAMS_WEBHOOK_URLS", "https://teams.webhook.url/1, https://teams.webhook.url/2")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cfg.Teams.Enabled {
		t.Errorf("expected Teams config to be enabled")
	}

	want := []string{"https://teams.webhook.url/1", "https://teams.webhook.url/2"}
	if !reflect.DeepEqual(cfg.Teams.WebhookURLs, want) {
		t.Errorf("got webhook URLs %v, want %v", cfg.Teams.WebhookURLs, want)
	}
}

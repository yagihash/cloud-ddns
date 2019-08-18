package config

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	keyLogPath        = "LOG_PATH"
	keyZoneConfig     = "ZONE_CONFIG"
	keySlackWebHook   = "SLACK_WEBHOOK"
	keySlackBotName   = "SLACK_BOTNAME"
	keySlackIconEmoji = "SLACK_ICONEMOJI"
	keySlackChannel   = "SLACK_CHANNEL"

	sampleLogPath        = "/tmp/ddns.log"
	sampleZoneConfig     = "testdata/zones.json"
	sampleSlackWebHook   = "https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/xxxxxxxxxxxxxxxxxxxxxxxx"
	sampleSlackBotName   = "Cloud-DDNS"
	sampleSlackIconEmoji = ":robot_face:"
	sampleSlackChannel   = ""
)

const (
	sampleZoneConfigValid = `
{
  "zones": [
    {
      "dns_name": "example.com.",
      "records": [
        "example.com.",
        "*.example.com."
      ]
    },
    {
      "dns_name": "example.jp.",
      "records": [
        "sub.example.jp."
      ]
    }
  ]
}
`
)

var (
	sampleEnv = map[string]string{
		keyLogPath:        sampleLogPath,
		keyZoneConfig:     sampleZoneConfig,
		keySlackWebHook:   sampleSlackWebHook,
		keySlackBotName:   sampleSlackBotName,
		keySlackIconEmoji: sampleSlackIconEmoji,
		keySlackChannel:   sampleSlackChannel,
	}
)

func SetupTestConfigAndEnv(t *testing.T) (*env, func()) {
	t.Helper()

	e := &env{}

	for k, v := range sampleEnv {
		if err := os.Setenv(k, v); err != nil {
			t.Error(err)
		}
		switch k {
		case keyLogPath:
			e.LogPath = v
		case keyZoneConfig:
			e.ZoneConfig = v
		case keySlackWebHook:
			e.SlackWebHook = v
		case keySlackBotName:
			e.SlackBotName = v
		case keySlackIconEmoji:
			e.SlackIconEmoji = v
		case keySlackChannel:
			e.SlackChannel = v
		default:
			t.Errorf("unrecognized env: %s=%s", k, v)
		}
	}

	e.replaceZone(t, sampleZoneConfigValid)

	return e, func() {
		for k := range sampleEnv {
			if err := os.Unsetenv(k); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestLoad(t *testing.T) {
	e, teardown := SetupTestConfigAndEnv(t)
	defer teardown()

	t.Run("Default", func(t *testing.T) {
		want := e.clone(t)
		got, err := Load()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("ReplaceLogPath", func(t *testing.T) {
		path := "/tmp/replaced.log"

		want := e.clone(t)
		want.LogPath = path

		if err := os.Setenv(keyLogPath, path); err != nil {
			assert.NoError(t, err)
		}

		got, err := Load()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("MissingRequiredEnv", func(t *testing.T) {
		if err := os.Unsetenv(keyZoneConfig); err != nil {
			assert.NoError(t, err)
		}

		_, err := Load()
		assert.Error(t, err)
	})

	t.Run("NoZoneConfig", func(t *testing.T) {
		path := "/tmp/nosuchfile.json"

		want := e.clone(t)
		want.ZoneConfig = path

		if err := os.Setenv(keyZoneConfig, path); err != nil {
			assert.NoError(t, err)
		}

		_, err := Load()
		assert.Error(t, err)
	})

	t.Run("NoRecordZone", func(t *testing.T) {
		path := "testdata/zonesNoRecord.json"

		want := e.clone(t)
		want.ZoneConfig = path

		if err := os.Setenv(keyZoneConfig, path); err != nil {
			assert.NoError(t, err)
		}

		_, err := Load()
		assert.Error(t, err)
	})

	t.Run("ZoneWithSyntaxError", func(t *testing.T) {
		path := "testdata/zonesSyntaxError.json"

		want := e.clone(t)
		want.ZoneConfig = path

		if err := os.Setenv(keyZoneConfig, path); err != nil {
			assert.NoError(t, err)
		}

		_, err := Load()
		assert.Error(t, err)
	})
}

func (e *env) clone(t *testing.T) *env {
	t.Helper()

	c := *e
	return &c
}

func (e *env) replaceZone(t *testing.T, jsonStr string) {
	t.Helper()

	if err := json.Unmarshal([]byte(jsonStr), &e.ManagedZones); err != nil {
		t.Errorf("invalid json string\n%s", jsonStr)
	}
}

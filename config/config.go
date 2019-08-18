package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type env struct {
	LogPath        string `envconfig:"LOG_PATH" default:"/tmp/ddns.log"`
	ZoneConfig     string `required:"true" envconfig:"ZONE_CONFIG"`
	SlackWebHook   string `required:"true" envconfig:"SLACK_WEBHOOK"`
	SlackBotName   string `envconfig:"SLACK_BOTNAME" default:"Cloud-DDNS"`
	SlackIconEmoji string `envconfig:"SLACK_ICONEMOJI" default:":robot_face:"`
	SlackChannel   string `envconfig:"SLACK_CHANNEL"`
	ManagedZones   *ManagedZones
}

func Load() (*env, error) {
	var e env
	if err := envconfig.Process("", &e); err != nil {
		return nil, err
	}

	f, err := os.Open(e.ZoneConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read zone config")
	}
	defer func() {
		_ = f.Close()
	}()

	b, err := ioutil.ReadAll(f)

	var managed *ManagedZones
	if err := json.Unmarshal(b, &managed); err != nil {
		return nil, errors.Wrap(err, "failed to parse zone config")
	}

	e.ManagedZones = managed

	for _, zone := range e.ManagedZones.Zones {
		if err := zone.validate(); err != nil {
			return nil, errors.Wrap(err, "invalid zone config")
		}
	}

	return &e, nil
}

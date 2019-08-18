package config

import "fmt"

type ManagedZones struct {
	Zones []Zone `json:"zones"`
}

type Zone struct {
	Name    string   `json:"name"`
	Records []string `json:"records"`
}

func (zone *Zone) validate() error {
	if len(zone.Records) == 0 {
		return fmt.Errorf("%s has no record", zone.Name)
	}
	return nil
}

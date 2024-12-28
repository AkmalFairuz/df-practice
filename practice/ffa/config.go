package ffa

import (
	"fmt"
	"github.com/akmalfairuz/df-practice/practice/helper"
	"gopkg.in/yaml.v3"
	"os"
)

type config struct {
	VoidY          int          `yaml:"void_y"`
	Spawns         [][5]float64 `yaml:"spawns"`
	AttackCooldown int          `yaml:"attack_cooldown"`
	Icon           string       `yaml:"icon"`
}

func (c config) SpawnLocations() []helper.Location {
	if len(c.Spawns) == 0 {
		panic("no spawn location")
	}
	return helper.ParseSliceOfLocation(c.Spawns)
}

var configs map[string]config

func init() {
	bytes, err := os.ReadFile("assets/ffa.yml")
	if err != nil {
		panic(fmt.Errorf("error read config: %v", err))
	}
	if err := yaml.Unmarshal(bytes, &configs); err != nil {
		panic(fmt.Errorf("error decode config: %v", err))
	}
}

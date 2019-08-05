package trigger

import (
)
type Trigger struct {
        Account string `yaml:"account"`
        Enabled bool `yaml:"enabled"`
        Organization string `yaml:"organization"`
        Registry string `yaml:"registry"`
        Repository string `yml:"repository"`
        Type string `yaml:"type"`
}

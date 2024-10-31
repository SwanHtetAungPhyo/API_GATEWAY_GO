package models

import "sync"

type Instance struct {
	Url         string     `yaml:"url"`
	Connections int32      `yaml:"connections"`
	Method      string     `yaml:"methods,omitempty"`
	Mu          sync.Mutex `yaml:"mu,omitempty"`
}

type Services struct {
	Name      string     `yaml:"name"`
	BasePath  string     `yaml:"basePath"`
	RateLimit int        `yaml:"rateLimit"`
	Instances []Instance `yaml:"instances"`
}

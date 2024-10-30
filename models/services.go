package models

import "sync"

type Instance struct {
	Url         string   `json:"url"`
	Connections int      `json:"connections"`
	Method    string `json:"methods,omitempty"`
	Mu  sync.Mutex
}

type Services struct {
	Name      string     `json:"name"`
	BasePath  string     `json:"basePath"`
	Instances []Instance `json:"instances"`
}

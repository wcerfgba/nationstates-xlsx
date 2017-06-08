package util

type Configurable interface {
	Configure(Configuration)
}

// Configuration is string map to any type, a convenience for... configuring
// stuff.
type Configuration map[string]interface{}

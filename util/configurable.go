package util

type Configurable interface {
	Configure(Configuration)
}

type Configuration map[string]interface{}
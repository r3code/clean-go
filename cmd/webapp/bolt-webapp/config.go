// +build !appengine

package main

type (
	// Config is an example of provider-specific configuration
	// in this case it's for the standalone version only to set
	// the boltdb filename
	Config struct {
		BoltDBFile string
	}
)

var (
	config *Config
)

func init() {
	// this would likely be loaded from flags or a conf file
	config = &Config{
		BoltDBFile: "app.boltdb",
	}
}

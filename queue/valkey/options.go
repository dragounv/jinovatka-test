package valkeyq

import (
	"os"
)

type ValkeyOptions struct {
	Addr string
	Port string
}

const (
	defaultValkeyAddr = "localhost"
	defaultValkeyPort = "6379"
)

// Create ValkeyOptions from enviroment
func NewValkeyOptionsFromEnv() *ValkeyOptions {
	valkeyAddr, ok := os.LookupEnv("VALKEY_ADDR")
	if !ok {
		valkeyAddr = defaultValkeyAddr
	}
	valkeyPort, ok := os.LookupEnv("VALKEY_PORT")
	if !ok {
		valkeyPort = defaultValkeyPort
	}
	return &ValkeyOptions{
		Addr: valkeyAddr,
		Port: valkeyPort,
	}
}

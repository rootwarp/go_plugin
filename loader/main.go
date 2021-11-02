package main

import (
	"errors"
	"reflect"
)

// Constants
const (
	PluginSymbolName = "Plugin"
)

// Errors
var (
	ErrNotSupportParameterType = errors.New("Retrieved parameter is not fit to plugin function")
)

type methodSpec struct {
	Name       string
	ParamTypes []reflect.Type
	Method     reflect.Value
}

// PluginSpec provides call specification for plugin.
type PluginSpec struct {
	Name        string
	Description string
	CallSpecs   map[string]CallSpec
}

// CallSpec is function spec which is provided by plugin.
type CallSpec struct {
	Name       string
	ParamTypes []reflect.Type
	Func       reflect.Value
}

// Config is glocal config values.
type Config struct {
	Plugins []PluginConfig `yaml:"plugins"`
}

// PluginConfig is config loaded from config.yaml.
type PluginConfig struct {
	Name string `yaml:"name"`
	Repo string `yaml:"repo"`
	Desc string `yaml:"desc"`
}

func main() {
}

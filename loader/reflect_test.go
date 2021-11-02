package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"plugin"
	"reflect"
	"strconv"
	"strings"
	"testing"

	git "github.com/go-git/go-git/v5"
	yaml "gopkg.in/yaml.v2"
)

func TestReflect(t *testing.T) {
	// Get input & parse
	userInput := "add 10 2"
	tokens := strings.Split(userInput, " ")

	// TODO: Load plugin from config.yml
	pluginSpec := loadPlugin("../plugin/test-plugin.so")

	for k, v := range pluginSpec.CallSpecs {
		fmt.Printf("%v - %+v\n", k, v)
	}

	ret, err := invokePlugin(pluginSpec, tokens[0], tokens[1:])
	if err != nil {
		log.Panic(err)
	}

	// Result.
	fmt.Println("Result ", ret[0])
}

// func loadPlugin(name string) map[string]methodSpec {
func loadPlugin(name string) PluginSpec {
	pluginSpec := PluginSpec{
		Name:        name,
		Description: "dummy",
	}

	// Load plugin
	// For real env. hashing pluging repository and use this as a key of map.
	p, err := plugin.Open(name)
	if err != nil {
		log.Panic(err)
	}

	sym, err := p.Lookup(PluginSymbolName)
	if err != nil {
		log.Println(err)
	}

	// Reflection process
	symType := reflect.TypeOf(sym)
	symVal := reflect.ValueOf(sym)

	pluginSpec.CallSpecs = map[string]CallSpec{}

	for i := 0; i < symType.NumMethod(); i++ {
		method := symType.Method(i)

		callSpec := CallSpec{
			Name:       method.Name,
			ParamTypes: make([]reflect.Type, method.Type.NumIn()-1),
			Func:       symVal.Method(i),
		}

		// First element of parameter is self instance. ignore it.
		for j := 0; j < method.Type.NumIn()-1; j++ {
			callSpec.ParamTypes[j] = method.Type.In(j + 1)
		}

		pluginSpec.CallSpecs[strings.ToLower(method.Name)] = callSpec

	}

	return pluginSpec
}

func invokePlugin(spec PluginSpec, method string, params []string) ([]reflect.Value, error) {
	s, ok := spec.CallSpecs[method]
	if !ok {
		log.Panicf("No function matched")
	}

	paramValues := make([]reflect.Value, len(s.ParamTypes))
	for i, paramType := range s.ParamTypes {
		val, err := convert(params[i], paramType)
		if err != nil {
			return nil, err
		}

		paramValues[i] = val
	}

	ret := s.Func.Call(paramValues)

	return ret, nil
}

func convert(in string, expectType reflect.Type) (reflect.Value, error) {
	// TODO: Add all types.
	switch expectType.Kind() {
	case reflect.Bool:
		v, err := strconv.ParseBool(in)
		return reflect.ValueOf(v), err
	case reflect.Int:
		v, err := strconv.ParseInt(in, 10, 64)
		return reflect.ValueOf(int(v)), err
	case reflect.Int8:
		v, err := strconv.ParseInt(in, 10, 8)
		return reflect.ValueOf(int8(v)), err
	default:
		return reflect.ValueOf(1), nil
	}
}

func TestYml(t *testing.T) {
	f, err := os.Open("./config.yml")
	if err != nil {
		log.Panic(err)
	}

	rawData, err := io.ReadAll(f)
	if err != nil {
		log.Panic(err)
	}

	parsed := Config{}
	err = yaml.Unmarshal(rawData, &parsed)
	if err != nil {
		log.Panic(err)
	}

	for _, p := range parsed.Plugins {
		fmt.Printf("%+v\n", p)
	}

	// TODO: Load repo.
	// TODO: How to store load info?
}

func TestGitClone(t *testing.T) {
	repo, err := git.PlainClone("./env", false, &git.CloneOptions{
		URL: "https://github.com/rootwarp/env",
	})

	fmt.Printf("Result %+v, %+v", repo, err)
}

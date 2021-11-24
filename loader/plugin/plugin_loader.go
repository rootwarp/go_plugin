package plugin

import (
	"errors"
	"log"
	"plugin"
	"reflect"
	"strconv"
	"strings"

	"go-plugin/config"
)

// Constants
const (
	PluginSymbolName = "Plugin"
)

// Errors
var (
	ErrSymbolNotExist = errors.New("Not exist symbol")
)

// Loader provides plugin control interfaces.
type Loader interface {
	Load(config config.Plugin) ([]Spec, error)

	// TODO: Draft
	// List() ([]Spec, error)
	GetPlugin(name string) (*Spec, error)
}

// Spec provides call specification for plugin.
type Spec struct {
	Name        string
	Description string
	CallSpecs   map[string]CallSpec
}

// Invoke calls selected function on plugin.
func (s *Spec) Invoke(funcName string, args ...string) error {
	spec, ok := s.CallSpecs[funcName]
	if !ok {
		log.Panicf("No function matched")
	}

	paramValues := make([]reflect.Value, len(spec.ParamTypes))
	for i, paramType := range spec.ParamTypes {
		val, err := convert(args[i], paramType)
		if err != nil {
			return err
		}

		paramValues[i] = val
	}

	_ = spec.Func.Call(paramValues)

	return nil
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
	case reflect.Int16:
		v, err := strconv.ParseInt(in, 10, 16)
		return reflect.ValueOf(int16(v)), err
	case reflect.Int32:
		v, err := strconv.ParseInt(in, 10, 32)
		return reflect.ValueOf(int32(v)), err
	case reflect.Int64:
		v, err := strconv.ParseInt(in, 10, 64)
		return reflect.ValueOf(int64(v)), err

	case reflect.String:
		return reflect.ValueOf(in), nil

	default:
		return reflect.ValueOf(1), nil
	}
}

// CallSpec is function spec which is provided by plugin.
type CallSpec struct {
	Name       string
	ParamTypes []reflect.Type
	Func       reflect.Value
}

type loader struct {
	plugins []Spec
}

func (l *loader) Load(config config.Plugin) ([]Spec, error) {
	// TODO: Read and load config.
	// TODO: Install not installed.
	return nil, nil
}

func (l *loader) install(repo string) error {
	// TODO:
	return nil
}

func (l *loader) loadSymbol(symbolName string) (map[string]CallSpec, error) {
	p, err := plugin.Open(symbolName)
	if err != nil {
		return nil, ErrSymbolNotExist
	}

	sym, err := p.Lookup(PluginSymbolName)
	if err != nil {
		return nil, err
	}

	funcCallSpecs := map[string]CallSpec{}

	symType := reflect.TypeOf(sym)
	symVal := reflect.ValueOf(sym)

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

		funcCallSpecs[strings.ToLower(method.Name)] = callSpec
	}

	// TODO: Check mandatory functions.

	return funcCallSpecs, nil
}

//func (l *loader) List() ([]Spec, error) {
//	return l.plugins, nil
//}

func (l *loader) GetPlugin(name string) (*Spec, error) {
	return nil, nil
}

// NewLoader creates plugin loader.
func NewLoader() Loader {
	return &loader{}
}

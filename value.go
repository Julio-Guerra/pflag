package pflag

import (
	"fmt"
	"reflect"
)

type DefaultValue interface {
	// The default value when the flag is provided on the command line but without
	// any value.
	// The empty string as a special meaning for nothing/none/undefined, i.e. the
	// option argument is mandatory.
	DefaultArg() string

	// The default value when the flag is not provided.
	// The empty string as a special meaning for nothing/none/undefined.
	DefaultValue() string
}

type DefaultValues struct {
	defaultValue interface{}
	defaultArg   interface{}
}

// Creates a new DefaultValues structure according to defaultValue and
// defaultArg.
// `t` allows to compute a zero value of v when defaultValue is nil. Otherwise,
// the type is inferred from defaultValue's reflected type.
func NewDefaultValues(t reflect.Type, defaultValue, defaultArg interface{}) (v interface{}, defaults *DefaultValues) {
	ov := &DefaultValues{
		defaultValue: shallowCopy(defaultValue),
		defaultArg:   shallowCopy(defaultArg),
	}
	if ov.defaultValue != nil {
		v = shallowCopy(defaultValue)
	} else if t != nil {
		v = reflect.Zero(t).Interface()
	}
	return v, ov
}

func (v *DefaultValues) DefaultArg() string {
	if v.defaultArg == nil {
		return ""
	}
	return fmt.Sprint(v.defaultArg)
}

func (v *DefaultValues) DefaultValue() string {
	if v.defaultValue == nil {
		return ""
	}
	return fmt.Sprint(v.defaultValue)
}

func shallowCopy(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	return reflect.ValueOf(value).Interface()
}

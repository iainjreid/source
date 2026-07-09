package config

import (
	"flag"
	"fmt"
	"os"
)

type Option struct {
	key string
	val flag.Value
	set bool
}

func NewOption[T flag.Value](key string, val T) (*Option, T) {
	return &Option{
		key: key,
		val: val,
	}, val
}

func (o *Option) String() string {
	return o.val.String()
}

func (o *Option) Set(str string) error {
	o.set = true
	return o.val.Set(str)
}

func (o *Option) LoadFromEnv() error {
	if !o.set {
		if str, ok := os.LookupEnv(o.key); ok {
			if err := o.Set(str); err != nil {
				return fmt.Errorf("invalid value \"%s\" for env %s: %v", str, o.key, err)
			}
		}
	}
	return nil
}

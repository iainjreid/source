// Copyright 2026 Iain J. Reid
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

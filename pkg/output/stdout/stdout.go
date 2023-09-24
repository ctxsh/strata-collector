// Copyright 2023 Rob Lyon <rob@ctxswitch.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stdout

import (
	"fmt"

	"ctx.sh/strata-collector/pkg/output"
)

type Stdout struct {
}

func New() *Stdout {
	return &Stdout{}
}

func (s *Stdout) Connect() error {
	return nil
}

func (s *Stdout) Send(data []byte) error {
	fmt.Println(string(data))
	return nil
}

func (s *Stdout) Close() {
}

var _ output.Output = &Stdout{}

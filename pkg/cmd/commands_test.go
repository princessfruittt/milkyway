// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	qt "github.com/frankban/quicktest"
	"testing"
)

func TestExecute(t *testing.T) {
	c := qt.New(t)

	c.Run("nginx from path", func(c *qt.C) {
		resp := Execute([]string{"generate", "-p", "assets/ansible-role-nginx-master"})
		c.Assert(resp.Err, qt.IsNil)
	})
	//c.Run("nginx from path", func(c *qt.C) {
	//	resp := Execute([]string{"generate", "-u", "https://github.com/geerlingguy/ansible-role-nginx"})
	//	c.Assert(resp.Err, qt.IsNil)
	//})
	c.Run("hosts from path", func(c *qt.C) {
		resp := Execute([]string{"generate", "-p", "assets/ansible-role-hosts-master"})
		c.Assert(resp.Err, qt.IsNil)
	})

	c.Run("users from path", func(c *qt.C) {
		resp := Execute([]string{"generate", "-p", "/home/princessfruitt/GolandProjects/milkyway/assets/ansible-users-master"})
		c.Assert(resp.Err, qt.IsNil)
	})
	//
	//c.Run("datadog from path", func(c *qt.C) {
	//	resp := Execute([]string{"generate", "-p", "/home/princessfruitt/GolandProjects/milkyway/assets/ansible-datadog-master"})
	//	c.Assert(resp.Err, qt.IsNil)
	//})
}

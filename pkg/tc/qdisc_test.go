// Copyright 2019 Yunion
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

package tc

import (
	"testing"
)

type tcCase struct {
	ifname      string
	line        string
	lineDelete  string
	lineReplace string
	isRoot      bool
}

func TestQdiscFqCodel(t *testing.T) {
	cases := []tcCase{
		{
			ifname:      "wp1-136",
			line:        "qdisc fq_codel 10: parent 1: limit 10240p flows 1024 quantum 1514 target 5.0ms interval 100.0ms ecn",
			lineDelete:  "qdisc delete dev wp1-136 parent 1: handle 10:",
			lineReplace: "qdisc replace dev wp1-136 parent 1: handle 10: fq_codel",
		},
		{
			ifname:      "vnet10-1",
			line:        "qdisc fq_codel 8003: root refcnt 2 limit 10240p flows 1024 quantum 1514 target 5.0ms interval 100.0ms ecn",
			lineDelete:  "qdisc delete dev vnet10-1 root handle 8003:",
			lineReplace: "qdisc replace dev vnet10-1 root handle 8003: fq_codel",
			isRoot:      true,
		},
	}
	for _, c := range cases {
		q, err := NewQdiscFromString(c.line)
		if err != nil {
			t.Errorf("%s", err)
			continue
		}
		if c.isRoot && !q.IsRoot() {
			t.Errorf("isroot want: %v, got %v", c.isRoot, q.IsRoot())
		}
		if lineDelete := q.DeleteLine(c.ifname); lineDelete != c.lineDelete {
			t.Errorf("delete line want: %s, got: %s", c.lineDelete, lineDelete)
			continue
		}
		if lineReplace := q.ReplaceLine(c.ifname); lineReplace != c.lineReplace {
			t.Errorf("replace line want: %s, got: %s", c.lineReplace, lineReplace)
			continue
		}
	}
}

func TestQdiscTbf(t *testing.T) {
	cases := []tcCase{
		{
			ifname:      "wp1-136",
			line:        "qdisc tbf 1: root refcnt 2 rate 500000Kbit burst 64000b/1 mpu 0b lat 100.0ms",
			lineDelete:  "qdisc delete dev wp1-136 root handle 1:",
			lineReplace: "qdisc replace dev wp1-136 root handle 1: tbf rate 500Mbit burst 64000b latency 100ms",
			isRoot:      true,
		},
		{
			ifname:      "wp1-136",
			line:        "qdisc tbf 1: root refcnt 2 rate 500000Kbit burst 64000b/4 mpu 0b lat 100.0ms",
			lineDelete:  "qdisc delete dev wp1-136 root handle 1:",
			lineReplace: "qdisc replace dev wp1-136 root handle 1: tbf rate 500Mbit burst 64000b/4 latency 100ms",
			isRoot:      true,
		},
	}

	for _, c := range cases {
		q, err := NewQdiscFromString(c.line)
		if err != nil {
			t.Errorf("%s", err)
			continue
		}
		if c.isRoot && !q.IsRoot() {
			t.Errorf("isroot want: %v, got %v", c.isRoot, q.IsRoot())
		}
		if lineDelete := q.DeleteLine(c.ifname); lineDelete != c.lineDelete {
			t.Errorf("delete line want: %s, got: %s", c.lineDelete, lineDelete)
			continue
		}
		if lineReplace := q.ReplaceLine(c.ifname); lineReplace != c.lineReplace {
			t.Errorf("replace line want: %s, got: %s", c.lineReplace, lineReplace)
			continue
		}
	}
}

func TestQdiscIngress(t *testing.T) {
	cases := []tcCase{
		{
			ifname:      "adm0-99",
			line:        "qdisc ingress ffff: dev adm0-99 parent ffff:fff1 ----------------",
			lineDelete:  "qdisc delete dev adm0-99 parent ffff:fff1 handle ffff:",
			lineReplace: "qdisc replace dev adm0-99 parent ffff:fff1 handle ffff: ingress",
			isRoot:      false,
		},
	}

	for _, c := range cases {
		q, err := NewQdiscFromString(c.line)
		if err != nil {
			t.Errorf("NewQdiscFromString: %s", err)
			continue
		}
		if c.isRoot && !q.IsRoot() {
			t.Errorf("isroot want: %v, got %v", c.isRoot, q.IsRoot())
			continue
		}
		if lineDelete := q.DeleteLine(c.ifname); lineDelete != c.lineDelete {
			t.Errorf("delete line want: %s, got: %s", c.lineDelete, lineDelete)
			continue
		}
		if lineReplace := q.ReplaceLine(c.ifname); lineReplace != c.lineReplace {
			t.Errorf("replace line want: %s, got: %s", c.lineReplace, lineReplace)
			continue
		}
	}
}

// Copyright 2025 The Score Authors
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

package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Ref[k any](in k) *k {
	return &in
}

func DerefOr[k any](in *k, def k) k {
	if in == nil {
		return def
	}
	return *in
}

func parseAndFormatResourceLimits(rl ResourcesLimits) string {
	c, m, err := ParseResourceLimits(rl)
	return fmt.Sprintf("%d %d %v", DerefOr(c, -1), DerefOr(m, -1), err)
}

func TestParseResourceLimits(t *testing.T) {
	assert.Equal(t, "-1 -1 <nil>", parseAndFormatResourceLimits(ResourcesLimits{}))
	assert.Equal(t, "1000 1000000 <nil>", parseAndFormatResourceLimits(ResourcesLimits{Cpu: Ref("1"), Memory: Ref("1M")}))
	assert.Equal(t, "-1 -1 failed to parse cpus 'banana' as a number", parseAndFormatResourceLimits(ResourcesLimits{Cpu: Ref("banana"), Memory: nil}))
	assert.Equal(t, "-1 -1 failed to parse memory 'banana' as a number", parseAndFormatResourceLimits(ResourcesLimits{Cpu: nil, Memory: Ref("banana")}))
	assert.Equal(t, "200 128974848 <nil>", parseAndFormatResourceLimits(ResourcesLimits{Cpu: Ref("200m"), Memory: Ref("123Mi")}))

	// Decimal memory values (the case from score-spec/spec#195).
	// 1.5Gi == 1.5 * 1024^3 == 1610612736 bytes.
	assert.Equal(t, "-1 1610612736 <nil>", parseAndFormatResourceLimits(ResourcesLimits{Memory: Ref("1.5Gi")}))
	// 0.5G == 0.5 * 1000^3 == 500000000 bytes.
	assert.Equal(t, "-1 500000000 <nil>", parseAndFormatResourceLimits(ResourcesLimits{Memory: Ref("0.5G")}))
	// 2.5M == 2.5 * 1000^2 == 2500000 bytes.
	assert.Equal(t, "-1 2500000 <nil>", parseAndFormatResourceLimits(ResourcesLimits{Memory: Ref("2.5M")}))

	// Peta- and exa-scale suffixes are accepted by the schema now; parser should handle them.
	// 1P == 1000^5 == 1e15.
	assert.Equal(t, "-1 1000000000000000 <nil>", parseAndFormatResourceLimits(ResourcesLimits{Memory: Ref("1P")}))
	// 1Pi == 1024^5 == 1125899906842624.
	assert.Equal(t, "-1 1125899906842624 <nil>", parseAndFormatResourceLimits(ResourcesLimits{Memory: Ref("1Pi")}))
}

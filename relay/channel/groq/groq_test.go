/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
package groq_test

import (
	"testing"

	"github.com/QuantumNous/new-api/relay/channel/groq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroqChannelName(t *testing.T) {
	assert.Equal(t, "groq", groq.ChannelName)
}

func TestGroqModelListIsUsable(t *testing.T) {
	require.NotEmpty(t, groq.ModelList)
	seen := make(map[string]struct{}, len(groq.ModelList))
	for _, m := range groq.ModelList {
		assert.NotEmpty(t, m)
		_, dup := seen[m]
		assert.False(t, dup, "duplicate model %s", m)
		seen[m] = struct{}{}
	}
	// Spot-check a couple of models that the default ratio table prices.
	assert.Contains(t, seen, "llama-3.1-8b-instant")
	assert.Contains(t, seen, "openai/gpt-oss-120b")
}

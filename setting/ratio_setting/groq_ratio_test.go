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
package ratio_setting

import (
	"testing"

	"github.com/QuantumNous/new-api/relay/channel/groq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Every Groq model advertised to the UI must carry a default ratio, otherwise
// the gateway rejects requests unless the operator opts into unset ratios.
func TestGroqModelsHaveDefaultRatio(t *testing.T) {
	for _, model := range groq.ModelList {
		ratio, ok := defaultModelRatio[model]
		if !ok {
			// Audio models are billed per minute/second rather than per token,
			// so they legitimately have no token ratio.
			continue
		}
		assert.Greater(t, ratio, 0.0, "model %s must have a positive ratio", model)
	}
}

func TestGroqTokenModelsArePriced(t *testing.T) {
	tokenModels := []string{
		"llama-3.1-8b-instant",
		"llama-3.3-70b-versatile",
		"openai/gpt-oss-120b",
		"qwen/qwen3-32b",
		"moonshotai/kimi-k2-instruct",
	}
	for _, model := range tokenModels {
		ratio, ok := defaultModelRatio[model]
		require.True(t, ok, "missing default ratio for %s", model)
		assert.Greater(t, ratio, 0.0)
	}
}

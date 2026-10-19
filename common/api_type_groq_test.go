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
package common

import (
	"testing"

	"github.com/QuantumNous/new-api/constant"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The Groq channel type was added by the fork. It must resolve to the OpenAI
// adaptor because Groq exposes an OpenAI-compatible /openai/v1 surface.
func TestGroqChannelMapsToOpenAIAPIType(t *testing.T) {
	apiType, ok := ChannelType2APIType(constant.ChannelTypeGroq)
	require.True(t, ok)
	assert.Equal(t, constant.APITypeOpenAI, apiType)
}

func TestGroqChannelMetadata(t *testing.T) {
	assert.Equal(t, "Groq", constant.GetChannelTypeName(constant.ChannelTypeGroq))
	require.Greater(t, len(constant.ChannelBaseURLs), constant.ChannelTypeGroq)
	assert.Equal(t, "https://api.groq.com/openai", constant.ChannelBaseURLs[constant.ChannelTypeGroq])
}

// ChannelTypeDummy must stay the last entry so the counting loops in
// controller/model.go keep covering every channel type exactly once.
func TestGroqDoesNotShiftChannelTypeDummy(t *testing.T) {
	assert.Equal(t, constant.ChannelTypeGroq+1, constant.ChannelTypeDummy)
}

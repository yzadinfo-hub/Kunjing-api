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
package groq

// ChannelName is the stable identifier used in logs and channel records.
var ChannelName = "groq"

// ModelList mirrors the production models published at
// https://console.groq.com/docs/models. Groq serves an OpenAI-compatible
// /openai/v1 surface, so requests are relayed by the OpenAI adaptor; this
// list only powers the "fetch models" helper and model suggestions.
var ModelList = []string{
	// Production LLMs
	"llama-3.1-8b-instant",
	"llama-3.3-70b-versatile",
	"meta-llama/llama-4-scout-17b-16e-instruct",
	"meta-llama/llama-4-maverick-17b-128e-instruct",
	"openai/gpt-oss-120b",
	"openai/gpt-oss-20b",
	"qwen/qwen3-32b",
	"moonshotai/kimi-k2-instruct",
	"deepseek-r1-distill-llama-70b",
	// Audio
	"whisper-large-v3",
	"whisper-large-v3-turbo",
}

# VERIFY.md — 修改验证清单

> **重要说明**：本次改造在**没有 Go 工具链**的环境中完成，所有 Go 代码均经过
> 逐行人工审查与静态检查（括号配平、import 存在性、符号引用、枚举值唯一性、
> map key 去重），但**未经编译器验证**。请在本地按本清单跑一遍再合并/上线。

---

## A. 环境准备

```bash
# Go 1.26+（见 go.mod）
go version

# Bun（前端）
bun --version
```

## B. 编译与静态检查

```bash
# B1. 全量编译（最关键，验证新增代码无语法/引用错误）
go build ./...

# B2. go vet
go vet ./...

# B3. gofmt（本 Fork 新增/修改的文件应全部通过）
gofmt -l common/api_type.go \
        constant/channel.go \
        relay/channel/openai/adaptor.go \
        relay/channel/groq/ \
        setting/ratio_setting/

# B4. 若 B3 有输出，执行 gofmt -w 后重新提交
```

**预期**：B1/B2 无输出；B3 无输出（若有，说明缩进/对齐需 gofmt 修正）。

## C. 单元测试

```bash
# C1. 本 Fork 新增的测试
go test ./common/ -run 'TestGroq' -v
go test ./relay/channel/groq/ -v
go test ./setting/ratio_setting/ -run 'TestGroq' -v

# C2. 回归：确认没有破坏上游既有测试
go test ./common/... ./relay/... ./setting/... ./controller/... 2>&1 | tail -40
```

**预期**：C1 全部 PASS；C2 与基线一致。
> 注意：上游 CI **没有**配置 `go test` 门禁，所以「上游全绿」不是既定事实。
> 若 C2 出现失败，先用 `git stash` 回到基线跑一遍同样命令对比：
> 失败项在基线也存在 → 上游既有问题，记录后继续；
> 失败项仅在本 Fork 出现 → 是你的改动引入的，必须修。

## D. 前端类型检查与测试

上游前端使用 **rsbuild**（非 vite）、`tsgo` 做类型检查、`oxlint` 做 lint、`vitest` 做测试。

```bash
cd web
bun install

# D1. TypeScript 类型检查（验证 constants.ts 改动）
bun run typecheck

# D2. Lint（验证新增代码风格）
bun run lint

# D3. 单元测试（含主题与渠道相关的既有用例）
bun run test

# D4. 生产构建（验证 theme-presets.css 与 i18n JSON 无语法错误）
bun run build

# D5. 版权头检查（新增/修改的 src 下文件必须带头）
bun run copyright:check

# D6. i18n 一致性（应报告 0 missing / 0 extras；Groq 不在未翻译列表）
bun run i18n:sync
```

**预期**：D1–D5 无错误；D6 报告中所有 locale `missingCount: 0`、`extrasCount: 0`。

## E. 功能验证（起服务后）

```bash
# E1. 启动（本地 SQLite 即可）
SQL_DSN=local MEMORY_CACHE_ENABLED=true \
SESSION_SECRET=$(openssl rand -base64 48) \
CRYPTO_SECRET=$(openssl rand -base64 48) \
go run . --log-dir ./logs
```

| # | 验证项 | 操作 | 预期 |
|---|---|---|---|
| E2 | Groq 渠道可选 | 渠道管理 → 新增渠道 → 类型下拉 | 出现 **Groq**，默认地址为 `https://api.groq.com/openai` |
| E3 | 模型列表 | 该渠道点「获取模型」 | 返回 `llama-3.1-8b-instant`、`openai/gpt-oss-120b` 等 |
| E4 | 计费倍率 | 用该渠道发一次 `/v1/chat/completions` | 日志里有额度消耗，且与默认倍率吻合 |
| E5 | 主题 | 右上角主题切换 | 出现 **Zumeng Spotlight**，且为新实例默认主题 |
| E6 | 暗色模式 | 切到暗色 | 主色变为暗色琥珀，图表/侧边栏同步 |
| E7 | 渠道回归 | 用 OpenAI / Claude / Gemini 渠道各发一次 | 全部正常（验证 import groq 未影响 openai adaptor） |
| E8 | 多语言 | 界面语言切到英文/日文 | 「Groq」正常显示，无 raw key 泄漏 |

## F. 每次合并上游后必做

```bash
# F1. 确认上游没加同编号渠道（会与本 Fork 的 64 冲突）
grep -n "ChannelType.*= 6[0-9]" constant/channel.go

# F2. 确认倍率表无重复 key（Go map 重复 key = 编译失败）
grep -oE '^\s*"[^"]+":' setting/ratio_setting/model_ratio.go \
  | tr -d ' \t' | sort | uniq -d
#   ↑ 只应输出跨 map 的既有项（gpt-4-all / tts-1 等属正常，它们分属不同 map）；
#     若输出里出现本 Fork 新增的 Groq 模型名，说明上游也加了同一个 key，需手工合并。

# F3. 确认 ChannelTypeDummy 仍是最后一个
grep -n "ChannelTypeDummy" constant/channel.go

# F4. 重跑 B/C/D 全部检查
```

## G. 已知未验证项（如实披露）

| 项 | 状态 |
|---|---|
| Go 代码编译 | ❌ 未验证（环境无 Go） |
| 单元测试执行 | ❌ 未验证（环境无 Go） |
| 前端 tsc / vitest / build | ❌ 未验证（环境无 bun） |
| i18n JSON 语法 | ✅ 已验证（7 个文件全部 `JSON.parse` 通过，且条目位于 `translation` 对象内） |
| 倍率表 key 唯一性 | ✅ 已验证（全表扫描，Groq 新增 8 个 key 无重复） |
| 枚举编号唯一性 | ✅ 已验证（64 未占用，且位于 Dummy 之前） |
| import 路径与符号存在性 | ✅ 已验证（逐个核对包路径与被引用符号） |

> 结论：**合并或上线前，请务必完整执行 B–E 四节。** 如任一步失败，
> 先修再提交；不要把未通过编译的代码推到 main。

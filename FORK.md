# FORK.md — 本 Fork 说明与 AGPLv3 合规声明

> **这份文件是本 Fork 的合规核心。** 上游 New API 采用 **AGPLv3 + 商业双许可**。
> 在仅使用 AGPLv3 的前提下，你有两项硬性义务：
> ① 保留全部上游品牌标识、LOGO 与版权声明（不得修改/移除/遮盖）；
> ② 若你通过网络提供服务（SaaS）或分发修改版，必须向你的用户提供**完整对应源码**。
>
> 本 Fork 的改造**全部是加法式的**（新增渠道、新增主题、新增部署套件），
> 没有修改或移除任何上游品牌标识。如需去除品牌标识或闭源商用，请联系
> support@quantumnous.com 获取商业授权。

---

## 1. Fork 身份

| 项 | 值 |
|---|---|
| Fork 来源 | <https://github.com/QuantumNous/new-api> |
| 基线提交 | `058c8231b68c51ee02abee15befcd720fee9f4a2`（2026-10-18） |
| 定制分支 | `feat/brand-customization` |
| Go module 路径 | `github.com/QuantumNous/new-api`（**刻意未改**，见 §4） |
| 许可证 | AGPLv3（保留上游 LICENSE 全文） |

## 2. 修改清单（相对上游基线）

全部改动共 **15 个文件修改 + 5 个新增文件**，约 **+260 行**：

### 2.1 新增自研渠道：Groq（channel type 64）

| 文件 | 改动 |
|---|---|
| `constant/channel.go` | 新增 `ChannelTypeGroq = 64`（插在 `ChannelTypeDummy` **之前**，保证 `controller/model.go` 的计数循环 `for i := 1; i <= ChannelTypeDummy` 覆盖新类型）；新增 `ChannelTypeNames` 条目与 `ChannelBaseURLs` 默认地址 `https://api.groq.com/openai` |
| `common/api_type.go` | `ChannelType2APIType` 新增 `ChannelTypeGroq → APITypeOpenAI` 分支 |
| `relay/channel/groq/constant.go` | **新增文件**：`ChannelName` 与 `ModelList`（10 个生产模型） |
| `relay/channel/openai/adaptor.go` | `GetModelList()` / `GetChannelName()` 增加 `ChannelTypeGroq` 分支；import `relay/channel/groq` |
| `setting/ratio_setting/model_ratio.go` | `defaultModelRatio` 新增 8 个 Groq 模型默认倍率（按 groq.com/pricing，$0.002/1K = 1） |
| `common/api_type_groq_test.go` | **新增**：映射、元数据、`ChannelTypeDummy` 未被推移的回归测试 |
| `relay/channel/groq/groq_test.go` | **新增**：模型列表非空/去重/定价抽检 |
| `setting/ratio_setting/groq_ratio_test.go` | **新增**：所有 token 类 Groq 模型必须有默认倍率 |

> **为什么选 Groq**：它是上游尚未覆盖的主流 OpenAI 兼容服务（LPU 推理，延迟极低、价格低），
> 且接入路径完全复用 OpenAI adaptor，风险最小。它是「如何向上游加一个新渠道」的参考实现。

### 2.2 品牌主题预设：`zumeng-spotlight`

| 文件 | 改动 |
|---|---|
| `web/src/styles/theme-presets.css` | **追加** `[data-theme-preset='zumeng-spotlight']` 亮/暗两套 CSS 变量（聚光灯琥珀主色 + 梅子紫辅色），上游 9 个预设**原样保留** |
| `web/src/lib/theme-customization.ts` | `THEME_PRESETS` 追加条目；`DEFAULT_THEME_CUSTOMIZATION.preset` 由 `default` 改为 `zumeng-spotlight` |

> 这是纯加法：用户仍可在主题选择器里切回任何上游预设。

### 2.3 前端渠道类型注册

| 文件 | 改动 |
|---|---|
| `web/src/features/channels/constants.ts` | 新增 `CHANNEL_TYPE_GROQ = 64`；`CHANNEL_TYPES`、`CHANNEL_PROVIDER_PRESENTATION` 增加条目；加入 `CHANNEL_TYPE_DISPLAY_ORDER`、`FIELD_PASSTHROUGH_TYPES`、`OPENAI_FIELD_PASSTHROUGH_TYPES` |
| `web/src/i18n/locales/{en,fr,ja,ru,vi,zh,zh-TW}.json` | 各新增 `"Groq"` 与 `"Connect to Groq LPU inference services"` 两条翻译 |
| `web/scripts/sync-i18n.mjs` | `BRAND_AND_LITERAL_KEYS` 增加 `'Groq'`，与 `'Replicate'`/`'Gemini'` 同级，避免被 i18n 检查器判为「未翻译」 |

> **i18n 已验证**：运行上游自带的 `bun run i18n:sync` 后，
> 7 个语言全部 `missingCount: 0`、`extrasCount: 0`，
> 且脚本规范化格式化后与手工插入结果**零差异**（说明插入位置与格式即项目规范格式）。

### 2.4 部署套件（`deploy/`，全新目录，不影响构建）

| 文件 | 用途 |
|---|---|
| `deploy/docker-compose.zumeng.yml` | 品牌化 Compose：强制密钥、日志独立库、健康检查依赖、Redis 持久化 |
| `deploy/.env.template` | 环境变量模板（必填项标注 + 生成命令） |
| `deploy/Caddyfile` | 自动 HTTPS 反代，含 SSE 关键配置 |
| `deploy/nginx.conf` | Nginx 备选反代，含 SSE 关键配置 |

## 3. 上游源码获取方式（AGPLv3 §13 合规）

任何人运行你部署的实例时，都必须能拿到**完整对应源码**。请确保：

1. 本仓库保持 **public**（最省事的合规方式）；或
2. 在系统设置的「自定义内容 / 页脚」中放上源码获取链接；或
3. 随产品分发一份源码归档或书面要约。

上游源码地址（必须一并提供）：
<https://github.com/QuantumNous/new-api>

## 4. 为什么没改 Go module 路径

`go.mod` 仍为 `github.com/QuantumNous/new-api`。这是**刻意决定**：

- 保留 module 路径可以让 `git merge upstream/main` / `git rebase` 几乎无冲突，
  持续获得上游的安全修复与新渠道；
- 改名需要动 ~1000 个文件的 import，且每次上游合并都要重做，成本极高、收益仅是「好看」。

品牌化完全由**部署层与主题层**承载（镜像名、域名、主题、站点名配置），
不需要侵入 module 路径。若未来确实需要改名，建议单独开一个分支做，并预期合并成本上升。

## 5. 如何持续跟进上游

```bash
git remote add upstream https://github.com/QuantumNous/new-api.git
git fetch upstream

# 常规跟进（推荐，保留自己的提交在顶部）
git checkout feat/brand-customization
git rebase upstream/main

# 或合并式跟进（冲突更少，历史更清晰）
git merge upstream/main --no-edit
```

由于本 Fork 只做了**加法**且集中在少数文件，冲突概率很低。
合并后必须重新执行 `VERIFY.md` 的检查清单（尤其 `ChannelTypeDummy` 计数与倍率表重复 key）。

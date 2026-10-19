# BRANDING.md — 品牌化实施指南与合规边界

本指南回答两个问题：**哪些品牌化可以做、怎么做**，以及**哪些绝对不能做**。

---

## 0. 先读这一条：合规边界

上游 LICENSE（AGPLv3 + 附加限制）明确规定：

> 在仅使用 AGPLv3 开源许可证的情况下，您**必须完整保留**项目代码中原有的品牌标识、LOGO
> 及版权声明信息。**禁止以任何形式修改、移除或遮盖**这些信息。如需移除，必须获取商业许可证。

| 操作 | 是否允许（纯 AGPLv3） | 说明 |
|---|---|---|
| 修改/删除/遮盖代码内的 New API LOGO | ❌ **禁止** | 需商业授权 |
| 删除版权声明、`LICENSE`、文件头 | ❌ **禁止** | 需商业授权 |
| 用自定义 CSS 遮罩上游 LOGO | ❌ **禁止** | 属「遮盖」 |
| **新增**自己的 LOGO / 品牌元素 | ✅ 允许 | 加法式共存 |
| 通过系统设置上传站点 LOGO | ✅ 允许 | 运行时配置，非改代码 |
| 改站点名称、公告、页脚自定义内容 | ✅ 允许 | 运行时配置 |
| 新增主题配色（如本 Fork 的 `zumeng-spotlight`） | ✅ 允许 | 加法，上游预设保留 |
| 换域名、镜像名、仓库名 | ✅ 允许 | 部署层 |
| 加自己的版权行 | ✅ 允许 | 建议加，表明你的修改 |

> **结论：品牌化的正确姿势是「叠加」而不是「替换」。**
> 你的品牌与上游品牌并存；若要做到「完全看不到上游」，必须购买商业授权。

---

## 1. 三层品牌化模型

```
┌─────────────────────────────────────────────────────────┐
│ L1 部署层（零代码，最安全，优先做）                          │
│   域名 / 镜像名 / 站点名称 / 公告 / 页脚 / 登录页文案         │
├─────────────────────────────────────────────────────────┤
│ L2 主题层（少量代码，加法式）                               │
│   新增主题预设（本 Fork 的 zumeng-spotlight）/ 默认字体       │
├─────────────────────────────────────────────────────────┤
│ L3 代码层（谨慎，仅新增，不修改上游品牌元素）                  │
│   自研渠道 / 自研计费规则 / 自研插件                          │
└─────────────────────────────────────────────────────────┘
```

---

## 2. L1 部署层：上线后第一件事（全部在管理后台完成）

路径：**系统设置 → 运营设置 → 通用设置**

| 配置项 | 建议值 | 备注 |
|---|---|---|
| 站点名称 `SystemName` | `Zumeng API` | 显示在标题栏、邮件、通知里 |
| 站点 LOGO | 上传自己的 PNG/SVG | 上传即替换展示，**代码内 LOGO 不动** |
| 首页公告 | 你的产品公告 | 支持 Markdown |
| 页脚自定义 HTML | `<a href="https://your.site">Powered by Zumeng</a>` | 追加在原有版权信息之后 |
| 用户协议 / 隐私政策 | 链接到你的域名 | 有独立页面 |

**主题**：系统设置 → 外观，选择 `Zumeng Spotlight`（本 Fork 新增的默认主题）。
想换配色：直接改 `web/src/styles/theme-presets.css` 里 `zumeng-spotlight` 的
`--primary` / `--secondary` 两个值即可，全站（含图表、侧边栏、暗色模式）自动联动。

---

## 3. L2 主题层：改品牌色

只需改两个变量（`web/src/styles/theme-presets.css`）：

```css
[data-theme-preset='zumeng-spotlight'] {
  --primary: oklch(0.6812 0.1738 55.01);   /* ← 主色：聚光灯琥珀 */
  --secondary: oklch(0.5392 0.1445 330.2); /* ← 辅色：梅子紫 */
  /* ...其余变量保持不动... */
}
.dark [data-theme-preset='zumeng-spotlight'] {
  --primary: oklch(0.7623 0.1526 62.44);   /* ← 暗色主色 */
  --secondary: oklch(0.6318 0.1356 328.71);
}
```

取色方法：用任意 OKLCH 取色器（如 <https://oklch.com/>）选色后填入。
不要直接删掉这个 preset——把它改成你的颜色即可，这样上游合并时不会冲突。

---

## 4. L3 代码层：加自己的东西

### 4.1 新增一个渠道（以上游未覆盖的服务为例）

完整范式见 `FORK.md` §2.1 的 Groq 实现，共 5 处改动：

```text
constant/channel.go                  ① ChannelTypeGroq = <下一个可用编号，插在 Dummy 之前>
constant/channel.go                  ② ChannelTypeNames / ChannelBaseURLs 加条目
common/api_type.go                   ③ ChannelType2APIType 加映射（OpenAI 兼容 → APITypeOpenAI）
relay/channel/<name>/constant.go     ④ 新增包：ChannelName + ModelList
relay/channel/openai/adaptor.go      ⑤ GetModelList()/GetChannelName() 加分支
setting/ratio_setting/model_ratio.go ⑥ 加默认倍率（先 grep 确认 key 不重复！）
web/src/features/channels/constants.ts ⑦ 前端注册 + i18n 翻译
```

**最容易踩的三个坑：**

1. **编号必须插在 `ChannelTypeDummy` 之前**，且不能与已用编号重复。
   合并上游后上游可能也加了新类型——先 `grep ChannelType constant/channel.go` 看最大编号。
2. **倍率表 key 不能重复**。Go 的 map 字面量出现重复 key 会**直接编译失败**。
   加之前务必 `grep -c '"模型名":' setting/ratio_setting/model_ratio.go`。
3. **i18n 是嵌套结构**：条目必须插在 `"translation": { ... }` **里面**，
   插到外层 i18next 找不到（本 Fork 开发时就踩过，见 git 历史）。

### 4.2 自定义计费规则

上游已内置**计费表达式 DSL**（`pkg/billingexpr`），支持阶梯价、固定价、参数联动：

```text
(len <= 32000 ? tier("short", fixed(0.01)) : tier("long", p * 2 + c * 8)) * (param("fast") == true ? 2 : 1)
```

在 **系统设置 → 倍率设置 → 模型倍率** 里直接填表达式即可，**不需要改代码**。
短剧场景建议预设：

| 场景 | 建议计费 |
|---|---|
| 字幕翻译（短输入长输出） | 按 token，`c` 倍率高于 `p` |
| 剧名/标签生成（极短） | 用 `fixed()` 按次计费，避免几分钱的订单 |
| 配音脚本改写（长上下文） | 开缓存倍率，命中缓存的 prompt 打折 |
| 批量翻译（夜间跑量） | 单独分组 + 分组倍率 0.8，与线上流量隔离 |

---

## 5. 绝对不要做的事

1. ❌ 删/改 `LICENSE`、任何文件头的版权声明、`web/public/logo.png`、
   `web/src/assets/logo.tsx`、`clerk-logo.tsx`、`waffo-logo-*.svg`
2. ❌ 用 CSS/JS 把上游 LOGO 遮掉
3. ❌ 把仓库改成 private 同时对外提供服务（违反 AGPLv3 §13 源码提供义务）
4. ❌ 宣称这是「你的原创项目」而不注明基于 New API

## 6. 想要「完全去品牌化」怎么办？

只有一条路：联系 **support@quantumnous.com** 购买商业授权。
在询价时可以把本 Fork 的改造量（约 +250 行、纯加法）作为谈判依据——
商业授权通常远比「自己重写一遍」便宜。

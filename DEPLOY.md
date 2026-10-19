# DEPLOY.md — Zumeng API 部署指南

## 0. 前置条件

| 项 | 要求 |
|---|---|
| Docker | 24+，含 Compose v2 |
| 内存 | ≥ 2 GB（应用 + PG + Redis） |
| 磁盘 | ≥ 10 GB（日志增长快，见 §5） |
| 域名 | 已解析到服务器（HTTPS 必须） |

## 1. 快速开始（5 分钟）

```bash
# 1. 准备配置
cd deploy
cp .env.template .env

# 2. 生成三个必填密钥（把输出填进 .env）
openssl rand -base64 48   # → SESSION_SECRET
openssl rand -base64 48   # → CRYPTO_SECRET
openssl rand -base64 24   # → POSTGRES_PASSWORD / REDIS_PASSWORD

# 3. 编辑 .env：填齐所有必填项，设置 SESSION_COOKIE_TRUSTED_URL 为你的域名

# 4. 启动
docker compose --env-file .env -f docker-compose.zumeng.yml up -d

# 5. 查看状态
docker compose --env-file .env -f docker-compose.zumeng.yml ps
docker logs -f zumeng-api
```

首次启动会自动建表并进入初始化向导。访问 <http://服务器IP:3000>，
按向导设置管理员账号与数据库即可。

## 2. 构建自己的镜像

默认镜像是占位符 `ghcr.io/your-org/zumeng-api:latest`，必须先自己构建。
上游 `Dockerfile` 已是三阶段构建（bun 构建前端 → golang:1.26-alpine 编译 → debian:bookworm-slim 运行时），
**直接用即可，无需改动**：

```bash
# 在仓库根目录（需要 bun + go 1.26）
docker build -t ghcr.io/your-org/zumeng-api:latest .

# 推送（如在 GHCR）
echo "$GITHUB_TOKEN" | docker login ghcr.io -u your-user --password-stdin
docker push ghcr.io/your-org/zumeng-api:latest
```

版本号由仓库根目录的 `VERSION` 文件提供（该文件默认为空，CI 会写入）。
本地构建想带版本号：

```bash
git describe --tags --always > VERSION   # 或手动写，如 v1.0.0-zumeng.1
docker build -t ghcr.io/your-org/zumeng-api:$(cat VERSION) .
```

> 上游 CI（`.github/workflows/release.yml`）支持打 tag 自动多架构构建，
> Fork 后只需把镜像仓库名改成自己的即可直接复用。

## 3. 接入 HTTPS 反代（必做）

### 3.1 Caddy（推荐，自动证书）

```bash
# 1. 编辑 deploy/Caddyfile，把 api.example.com 全部替换成你的域名
# 2. 让 Caddy 加入同一个 docker network
docker run -d --name zumeng-caddy --restart always \
  --network zumeng-api_zumeng-network \
  -p 80:80 -p 443:443 -p 443:443/udp \
  -v "$PWD/deploy/Caddyfile":/etc/caddy/Caddyfile:ro \
  -v zumeng_caddy_data:/data -v zumeng_caddy_config:/config \
  caddy:2-alpine
```

### 3.2 Nginx（备选）

```bash
docker run -d --name zumeng-nginx --restart always \
  --network zumeng-api_zumeng-network \
  -p 80:80 -p 443:443 \
  -v "$PWD/deploy/nginx.conf":/etc/nginx/conf.d/default.conf:ro \
  -v "$PWD/tls":/etc/nginx/tls:ro \
  nginx:1.27-alpine
```

证书用 certbot 申请后放到 `./tls/`。

### 3.3 为什么 `SESSION_COOKIE_SECURE=true` 是默认值

反代终止 TLS 后，应用层必须发 `Secure` Cookie，否则会话可被中间人窃取。
配套要求：

- `SESSION_COOKIE_TRUSTED_URL` 必须是**精确的 HTTPS Origin**（不含路径、不支持通配符），
  多个 Origin 用英文逗号分隔；
- `TRUSTED_PROXIES` 要覆盖反代所在网段（Compose 默认 `172.20.0.0/16`，
  自定义网络时用 `docker network inspect <网络名> -f '{{.IPAM.Config}}'` 确认）。

> 本地 HTTP 调试时才设 `SESSION_COOKIE_SECURE=false`。

## 4. 上线后必做清单

- [ ] 改 `.env` 全部默认密码（Compose 已用 `:?` 强制，缺失即启动失败）
- [ ] 系统设置里改站点名称、上传 LOGO、配置公告与页脚
- [ ] 外观选择 `Zumeng Spotlight` 主题
- [ ] 关掉注册或加验证码（系统设置 → 运营设置 → 登录注册）
- [ ] 配 SMTP，用于额度不足提醒与密码重置
- [ ] 加一个渠道并测试：渠道管理 → 测试
- [ ] 配 `LOG_SQL_DSN` 独立日志库（Compose 已默认开启）
- [ ] 设定日志清理策略（系统设置 → 运营设置 → 日志）
- [ ] 加监控：`GET /api/status` 健康探针 + 容器 healthcheck 已配

## 5. 运维要点

### 备份

```bash
# 数据库（最关键）
docker exec zumeng-postgres pg_dump -U zumeng zumeng | gzip > backup-$(date +%F).sql.gz

# 日志库
docker exec zumeng-postgres pg_dump -U zumeng zumeng_log | gzip > backup-log-$(date +%F).sql.gz

# Redis（AOF 已开启，容器内 /data）
docker exec zumeng-redis redis-cli -a "$REDIS_PASSWORD" BGSAVE
```

建议每天 `pg_dump` + 保留 14 天；`/data` 与 `./logs` 纳入主机级备份。

### 日志增长

流式网关的日志表增长极快（一次对话一行，含 prompt 摘要）。
务必：

1. `LOG_SQL_DSN` 指向独立库（Compose 已默认）；
2. 系统设置里配置日志保留天数，或定期调用日志清理接口；
3. `./logs` 目录配 logrotate 或 Docker logging driver 的 `max-size`。

### 升级

```bash
git fetch upstream && git merge upstream/main --no-edit   # 先合并上游
docker compose --env-file .env -f deploy/docker-compose.zumeng.yml pull
docker compose --env-file .env -f deploy/docker-compose.zumeng.yml up -d
```

升级后**必须**跑 `VERIFY.md` 的检查清单。数据库迁移是自动的（GORM AutoMigrate），
但升级前请先 `pg_dump`。

### 多节点

- 所有节点 `SESSION_SECRET` 与 `CRYPTO_SECRET` 必须一致；
- 共用同一个 PostgreSQL 与 Redis；
- 只让一个节点设为主节点（默认即主节点，`NODE_TYPE=slave` 的是从节点）；
- 多节点时给每个节点不同的 `NODE_NAME`，便于审计日志区分。

## 6. 故障速查

| 现象 | 原因 | 处理 |
|---|---|---|
| 启动即退出，提示 secret | `.env` 必填项未填 | 按提示补全 |
| 登录后立刻掉登录态 | `SESSION_COOKIE_SECURE=true` 但走了 HTTP | 本地调试设 `false`，或走 HTTPS |
| SSE 流式卡住不输出 | 反代开了缓冲 | Nginx 设 `proxy_buffering off`；Caddy 用 `flush_interval -1` |
| 渠道测试超时 | `STREAMING_TIMEOUT` 太小或上游慢 | 调到 600 |
| 额度扣了但请求失败 | 正常，`ReturnPreConsumedQuota` 会返还 | 看日志确认；若未返还是 bug，上报 |
| 429 限流 | 触发了全局限流 | 系统设置 → 限流设置调整 |
| 上游报模型无倍率 | 新模型未配倍率 | 系统设置 → 模型倍率添加，或开启「允许未配置倍率的模型」 |

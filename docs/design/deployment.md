# デプロイ設計

## 構成概要

```
┌─────────────────────────────────────────┐
│  docker-compose                         │
│                                         │
│  ┌──────────┐  ┌──────────┐  ┌───────┐  │
│  │ frontend │  │ backend  │  │  db   │  │
│  │  Nginx   │  │  Go+Gin  │─▶│  PG   │  │
│  │  :80     │  │  :8080   │  │ :5432 │  │
│  └────┬─────┘  └────┬─────┘  └───────┘  │
└───────┼─────────────┼───────────────────┘
        │             │
   静的配信        /api/* プロキシ
```

フロントエンドのNginxが `/api/*` リクエストをバックエンドにリバースプロキシする。

## docker-compose.yml

```yaml
services:
  frontend:
    build: ./frontend
    ports:
      - "80:80"
    depends_on:
      - backend

  backend:
    build: ./backend
    env_file: .env
    depends_on:
      db:
        condition: service_healthy

  db:
    image: mysql:8.4
    environment:
      MYSQL_DATABASE: moneyrabbit
      MYSQL_USER: moneyrabbit
      MYSQL_PASSWORD: ${MYSQL_PASSWORD}
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
    volumes:
      - mysql_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 5s
      retries: 5

volumes:
  mysql_data:
```

## Dockerfile

### フロントエンド（Vite + Nginx）

```dockerfile
FROM node:22-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
```

```nginx
# frontend/nginx.conf
server {
  listen 80;

  location / {
    root /usr/share/nginx/html;
    try_files $uri $uri/ /index.html;
  }

  location /api/ {
    proxy_pass http://backend:8080;
  }
}
```

### バックエンド（Go）

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

## 環境変数（.env）

| 変数名 | 内容 |
|---|---|
| `DATABASE_URL` | `moneyrabbit:password@tcp(db:3306)/moneyrabbit?parseTime=true` |
| `MYSQL_PASSWORD` | MySQLユーザーパスワード |
| `MYSQL_ROOT_PASSWORD` | MySQLルートパスワード |
| `ANTHROPIC_API_KEY` | Claude API キー |

`.env` はGit管理外（`.gitignore`）。`.env.example` をリポジトリに含める。

## DBマイグレーション

gooseでバージョン管理する。バックエンド起動時に自動実行する。

```go
// backend起動時に実行
goose.Up(db, "db/migrations")
```

スキーマ変更時のフローは [data-model.md](./data-model.md) を参照。

## リポジトリ構成

```
MoneyRabbit/
├── frontend/           # Vite + React
│   ├── src/
│   ├── nginx.conf
│   └── Dockerfile
├── backend/            # Go + Gin
│   ├── cmd/server/
│   ├── ent/
│   │   └── schema/     # entスキーマ定義
│   └── Dockerfile
├── docker-compose.yml
├── .env.example
└── docs/
```

## 本番環境（Kubernetes）

Docker Composeと同じイメージを使用する。  
Kubernetes固有のマニフェストは別途管理する。

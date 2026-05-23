# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Chart 概要

`charts/moneyrabbit/` は Kubernetes 向け Helm Chart。frontend・backend・MariaDB の3リソースを含む。

GitHub Pages でホスト: `https://itk13201.github.io/MoneyRabbit`

## リリースフロー

`charts/**` を `main` に push すると `helm-release.yml` が自動実行され、GitHub Release と `gh-pages` の `index.yaml` が更新される。

Chart のリリース時は **必ず** `Chart.yaml` の `version` を上げること（上げないと chart-releaser がスキップする）。`appVersion` と `values.yaml` のデフォルトイメージタグも合わせて更新する。

| 変数 | 場所 | 例 |
|---|---|---|
| chart version | `Chart.yaml` | `0.3.0` → `0.4.0` |
| appVersion | `Chart.yaml` | `"1.0.2"` → `"1.1.0"` |
| image tag | `values.yaml` | `tag: "1.0.2"` → `tag: "1.1.0"` |

## Kubernetes へのインストール

```bash
helm repo add moneyrabbit https://itk13201.github.io/MoneyRabbit
helm repo update

# 新規インストール
helm install mr moneyrabbit/moneyrabbit \
  --namespace moneyrabbit --create-namespace \
  --set mariadb.auth.password=YOUR_PASSWORD \
  --set mariadb.auth.rootPassword=YOUR_ROOT_PASSWORD \
  --set backend.anthropicApiKey=YOUR_API_KEY

# アップグレード
helm upgrade mr moneyrabbit/moneyrabbit \
  --namespace moneyrabbit \
  --reuse-values \
  --set frontend.image.tag=1.1.0 \
  --set backend.image.tag=1.1.0
```

## existingSecret（本番環境推奨）

パスワードを values に直書きせず、既存の Kubernetes Secret を参照できる。

**backend 用 Secret**（キー名は固定）:
```bash
kubectl create secret generic moneyrabbit-backend \
  --namespace moneyrabbit \
  --from-literal=DATABASE_URL='moneyrabbit:pass@tcp(mr-moneyrabbit-mariadb:3306)/moneyrabbit?parseTime=true' \
  --from-literal=ANTHROPIC_API_KEY='sk-ant-...'
```

**mariadb 用 Secret**（キー名は固定）:
```bash
kubectl create secret generic moneyrabbit-mariadb \
  --namespace moneyrabbit \
  --from-literal=MYSQL_ROOT_PASSWORD='rootpass' \
  --from-literal=MYSQL_USER='moneyrabbit' \
  --from-literal=MYSQL_PASSWORD='pass' \
  --from-literal=MYSQL_DATABASE='moneyrabbit'
```

```bash
helm install mr moneyrabbit/moneyrabbit \
  --namespace moneyrabbit --create-namespace \
  --set backend.existingSecret=moneyrabbit-backend \
  --set mariadb.auth.existingSecret=moneyrabbit-mariadb
```

## ローカルでのテンプレート確認

```bash
helm template test charts/moneyrabbit \
  --set mariadb.auth.password=p \
  --set mariadb.auth.rootPassword=r \
  --set backend.anthropicApiKey=k

helm lint charts/moneyrabbit \
  --set mariadb.auth.password=p \
  --set mariadb.auth.rootPassword=r \
  --set backend.anthropicApiKey=k
```

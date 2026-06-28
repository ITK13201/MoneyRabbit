# リリース手順

## 概要

git push だけで全て自動化されている。

| トリガー | 自動実行される処理 |
|---|---|
| `v*` タグを push | `build.yml` — frontend/backend Docker イメージをビルドして GHCR にプッシュ |
| `charts/**` を main に push | `helm-release.yml` — GitHub Pages の Helm Chart リポジトリを更新 |

## バージョン番号の管理

| ファイル | フィールド | 例 |
|---|---|---|
| `frontend/package.json` | `version` | `"1.0.4"` |
| `charts/moneyrabbit/Chart.yaml` | `appVersion` | `"1.0.4"` |
| `charts/moneyrabbit/Chart.yaml` | `version` | `0.6.0`（Chart 自体のバージョン） |
| `charts/moneyrabbit/values.yaml` | `frontend.image.tag` | `"1.0.4"` |
| `charts/moneyrabbit/values.yaml` | `backend.image.tag` | `"1.0.4"` |

`appVersion` / `package.json` はアプリのバージョン、`Chart.yaml` の `version` は Chart 構造の変更に対してインクリメントする（アプリのみ更新する場合は patch だけ上げれば良い）。

## 手順

`/release 1.1.0` を実行すると以下を自動で行う。

### Step 1 — バージョン番号を一括更新してコミット & push

```bash
# frontend/package.json, Chart.yaml, values.yaml を編集後:
git add frontend/package.json charts/moneyrabbit/Chart.yaml charts/moneyrabbit/values.yaml
git commit -m "chore: bump version to 1.1.0"
git push origin main  # → helm-release.yml が自動実行
```

### Step 2 — git タグを push

```bash
git tag v1.1.0
git push origin v1.1.0  # → build.yml が自動実行
```

ビルドの進捗は [GitHub Actions](https://github.com/itk13201/MoneyRabbit/actions) で確認する。

## チェックリスト

```
[ ] frontend/package.json version を更新
[ ] Chart.yaml version を更新（patch / minor / major）
[ ] Chart.yaml appVersion を更新
[ ] values.yaml frontend.image.tag / backend.image.tag を更新
[ ] コミット & push → helm-release.yml が自動実行
[ ] git tag v<version> を push → build.yml が自動実行
```

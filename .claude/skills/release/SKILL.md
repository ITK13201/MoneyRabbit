---
name: release
description: MoneyRabbit のリリース手順を実行する。バージョン番号の更新・コミット・git タグの push までを行い、以降は GitHub Actions が自動でリリースする。
disable-model-invocation: true
---

# release スキル

引数にリリースするバージョン番号を指定して呼び出す（例: `/release 1.1.0`）。

## 前提確認

開始前に以下を確認し、問題があればユーザーに報告して中断する。

- `$ARGUMENTS` が指定されているか確認する。未指定の場合は「バージョン番号を指定してください（例: /release 1.1.0）」と伝えて終了する。
- 作業ディレクトリが clean か確認する（`git status --porcelain`）。未コミットの変更がある場合はユーザーに確認を求める。
- 現在のブランチが `main` か確認する（`git branch --show-current`）。異なる場合はユーザーに確認を求める。
- 指定バージョンのタグが既に存在しないか確認する（`git tag -l "v$ARGUMENTS"`）。存在する場合は中断する。

## Step 1 — バージョン番号の更新

以下の3ファイルを編集する。

**`frontend/package.json`**
```json
"version": "<ARGUMENTS>"
```

**`charts/moneyrabbit/Chart.yaml`**
- `appVersion` を `"<ARGUMENTS>"` に更新
- `version` を patch インクリメントする（例: `0.6.0` → `0.7.0`）
  - Chart のテンプレートや values スキーマも変更した場合は minor または major を検討する
  - ユーザーに現在の Chart version を提示し、変更後の値を確認してから編集する

**`charts/moneyrabbit/values.yaml`**
- `frontend.image.tag` を `"<ARGUMENTS>"` に更新
- `backend.image.tag` を `"<ARGUMENTS>"` に更新

編集後、差分を表示してユーザーに確認を求める（`git diff`）。

## Step 2 — コミット & push

```bash
git add frontend/package.json charts/moneyrabbit/Chart.yaml charts/moneyrabbit/values.yaml
git commit -m "chore: bump version to <ARGUMENTS>"
git push origin main
```

`charts/**` の変更が main に push されることで `helm-release.yml` が自動実行される。

## Step 3 — git タグを push

```bash
git tag v<ARGUMENTS>
git push origin v<ARGUMENTS>
```

`v*` タグの push により `build.yml` が自動実行され、以下の Docker イメージが GHCR にプッシュされる：

- `ghcr.io/itk13201/moneyrabbit-frontend:<ARGUMENTS>`
- `ghcr.io/itk13201/moneyrabbit-backend:<ARGUMENTS>`

## 完了報告

以下をまとめて報告する：

- リリースされたバージョン
- 作成された git タグ
- 更新された Chart version
- GitHub Actions の確認 URL: https://github.com/itk13201/MoneyRabbit/actions

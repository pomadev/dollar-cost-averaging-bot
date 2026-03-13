# dollar-cost-averaging-bot

bitbank で BTC を定期購入し、結果を LINE に通知する Cloud Run Job 向けのボットです。

## 概要

- 1 回の実行で BTC を 1 回だけ成行購入します
- 購入額は `BTC_YEN` で指定します
- 実行結果を LINE Messaging API で通知します
- `DRY_RUN=true` の場合は注文せず、想定数量だけ計算して通知します

## 前提

- Go 1.26.1
- bitbank の API キー
- LINE Messaging API のチャネルシークレット
- LINE Messaging API のチャネルアクセストークン
- 通知先の `LINE_USER_ID`
- GCP で定期実行する場合:
  - GCP project 作成済み
  - 課金アカウント紐付け済み
  - `gcloud` インストール済み
  - Cloud Run
  - Cloud Scheduler
  - Secret Manager
  - Cloud Build
  - Artifact Registry

## 環境変数

| 名前 | 必須 | 説明 |
| --- | --- | --- |
| `BTC_YEN` | yes | 1 回の購入金額。例: `5000` |
| `BITBANK_ACCESS_KEY` | yes | bitbank API key |
| `BITBANK_API_SECRET` | yes | bitbank API secret |
| `LINE_CHANNEL_SECRET` | yes | LINE channel secret |
| `LINE_CHANNEL_ACCESS_TOKEN` | yes | LINE channel access token |
| `LINE_USER_ID` | yes | Push 通知先の LINE user ID |
| `DRY_RUN` | no | `true` の場合は実注文しない。デフォルトは `false` |

## ローカル実行

まずは `DRY_RUN=true` で確認してください。

```bash
export BTC_YEN=5000
export DRY_RUN=true
export BITBANK_ACCESS_KEY=xxxxx
export BITBANK_API_SECRET=xxxxx
export LINE_CHANNEL_SECRET=xxxxx
export LINE_CHANNEL_ACCESS_TOKEN=xxxxx
export LINE_USER_ID=Uxxxxxxxx

go run ./cmd
```

## Cloud Run Job での実行

### 初期セットアップ

```bash
gcloud auth login
gcloud auth application-default login
gcloud config set project YOUR_PROJECT_ID
```

プロジェクト番号の確認:

```bash
gcloud projects describe YOUR_PROJECT_ID --format='value(projectNumber)'
```

必要な API を有効化します。

```bash
gcloud services enable \
  run.googleapis.com \
  cloudbuild.googleapis.com \
  artifactregistry.googleapis.com \
  secretmanager.googleapis.com \
  cloudscheduler.googleapis.com
```

### サービスアカウントと IAM

Cloud Run Job の実行用サービスアカウントを作成します。Job が Secret Manager から値を読むために使います。

```bash
gcloud iam service-accounts create dca-bot-runner \
  --display-name "DCA Bot Runner"
```

そのサービスアカウントに Secret Manager の読み取り権限を付与します。

```bash
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member "serviceAccount:dca-bot-runner@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role "roles/secretmanager.secretAccessor"
```

次に、Cloud Scheduler から Cloud Run Job を起動するための専用サービスアカウントを作成します。

```bash
gcloud iam service-accounts create dca-bot-scheduler \
  --display-name "DCA Bot Scheduler"
```

そのサービスアカウントに Cloud Run Job 実行 API を呼ぶ権限を付与します。

```bash
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member "serviceAccount:dca-bot-scheduler@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role "roles/run.invoker"
```

Cloud Run Job の実行サービスアカウントを分けたい場合は、別途作成して `gcloud run jobs deploy` に `--service-account` で指定してください。

Secret Manager に機密値を登録します。

```bash
printf '%s' 'your-bitbank-access-key' | gcloud secrets create bitbank-access-key --data-file=-
printf '%s' 'your-bitbank-api-secret' | gcloud secrets create bitbank-api-secret --data-file=-
printf '%s' 'your-line-channel-secret' | gcloud secrets create line-channel-secret --data-file=-
printf '%s' 'your-line-channel-access-token' | gcloud secrets create line-channel-access-token --data-file=-
printf '%s' 'Uxxxxxxxx' | gcloud secrets create line-user-id --data-file=-
```

Cloud Run Job をデプロイします。

```bash
gcloud run jobs deploy dca-bot \
  --source . \
  --region asia-northeast1 \
  --service-account dca-bot-runner@YOUR_PROJECT_ID.iam.gserviceaccount.com \
  --tasks 1 \
  --max-retries 0 \
  --task-timeout 300s \
  --set-env-vars BTC_YEN=5000,DRY_RUN=false \
  --set-secrets BITBANK_ACCESS_KEY=bitbank-access-key:latest \
  --set-secrets BITBANK_API_SECRET=bitbank-api-secret:latest \
  --set-secrets LINE_CHANNEL_SECRET=line-channel-secret:latest \
  --set-secrets LINE_CHANNEL_ACCESS_TOKEN=line-channel-access-token:latest \
  --set-secrets LINE_USER_ID=line-user-id:latest
```

`--service-account` を省略すると default compute service account が使われ、Secret Manager の権限不足で deploy に失敗することがあります。

手動実行:

```bash
gcloud run jobs execute dca-bot --region asia-northeast1
```

## Cloud Scheduler での定期実行

毎月 1 日 09:00 JST に実行する例です。

```bash
gcloud scheduler jobs create http dca-bot-monthly \
  --location asia-northeast1 \
  --schedule '0 9 1 * *' \
  --time-zone 'Asia/Tokyo' \
  --uri "https://run.googleapis.com/v2/projects/YOUR_PROJECT_ID/locations/asia-northeast1/jobs/dca-bot:run" \
  --http-method POST \
  --oauth-service-account-email dca-bot-scheduler@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

既存 Job を更新したい場合は `gcloud scheduler jobs update http ...` を使ってください。

## 注意点

- `DRY_RUN=false` では実際に注文が発生します
- 少額すぎると bitbank 側の最小数量に届かず失敗します
- 購入後の約定価格は API の板状況で変動します

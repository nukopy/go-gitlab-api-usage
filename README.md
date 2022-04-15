# GitLab API Usage

Go で GitLab API を叩く。サンプルとして、特定のグループの全てのプロジェクトを取得して CSV ファイルに吐き出すやつを書いてみた。

## 環境

- OS: macOS 11.6
- CPU: Intel(R) Core(TM) i9-9880H CPU 2.30GHz
- Go 1.18

## 実行方法

1. [GitLab のプロフィールページにて Personal Access Token を発行](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token)
2. `.env` の作成
3. 実行

### 1. GitLab のプロフィールページにて Personal Access Token を発行

下記を参照。

- [GitLab Docs: Create a personal access tokens](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token)

### 2. `.env` の作成

- `.env`

```sh
GITLAB_TOKEN="*****"
GITLAB_GROUP_ID="*****"
```

### 2. 実行

```sh
cd go-gitlab-api-usage/
go run .
```

package main

import "log"

// GitLab の特定のグループの全てのプロジェクト（リポジトリ）の情報を CSV に吐き出す
func main() {
	// 環境変数の読み込み
	envs := loadEnvs()
	gitlabToken := envs.GITLAB_TOKEN
	gitlabGroupId := envs.GITLAB_GROUP_ID

	// GitLab API クライアントの作成
	cli, err := createGitlabClient(gitlabToken)
	if err != nil {
		log.Fatal(err)
	}

	// 特定のグループの全ての Gitlab プロジェクトを取得// 特定のグループの全ての Gitlab プロジェクトを取得
	gitlabProjects, err := fetchAllGitlabProjectsInGroup(cli, gitlabGroupId)
	if err != nil {
		log.Fatal(err)
	}

	// CSV ファイルへ出力
	err = outputGitlabProjectsToCsv(gitlabGroupId, gitlabProjects)
	if err != nil {
		log.Fatal(err)
	}
}

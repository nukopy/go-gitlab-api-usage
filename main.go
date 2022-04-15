package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/xanzy/go-gitlab"
)

// GitLab の特定のグループの全てのプロジェクト（リポジトリ）の情報を CSV に吐き出す
func main() {
	// 環境変数の読み込み
	envs := loadEnvs()
	gitlabToken := envs.GITLAB_TOKEN
	gitlabGroupId := envs.GITLAB_GROUP_ID

	outputAllGitLabProjectsInfoToCsv(gitlabToken, gitlabGroupId)
}

func outputAllGitLabProjectsInfoToCsv(gitlabToken string, gitlabGroupId string) {
	// GitLab API クライアントの作成
	git, err := gitlab.NewClient(gitlabToken)
	if err != nil {
		log.Fatal(err)
	}

	// 大元のグループの取得
	group, _, err := git.Groups.GetGroup(gitlabGroupId, &gitlab.GetGroupOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Group ID   : %d\n", group.ID)
	log.Printf("Group Name : %s\n", group.Name)

	// 全てのグループ、サブグループの取得
	log.Printf("Getting all groups or subgroups \"%s\"...\n", group.Name)

	var allGroups []*gitlab.Group
	for {
		groups, resp, err := git.Groups.ListDescendantGroups(gitlabGroupId, &gitlab.ListDescendantGroupsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    1,
				PerPage: 100,
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		allGroups = append(allGroups, groups...)

		// paging
		if resp.NextPage == 0 {
			break
		}
	}

	allGroups = append([]*gitlab.Group{group}, allGroups...) // 大元のグループを先頭に追加

	// 各グループ、各サブグループのプロジェクト情報の取得
	log.Printf("Getting all projects in every group \"%s\"...\n", group.Name)

	var allProjects []*gitlab.Project
	for _, group := range allGroups {
		for {
			projects, resp, err := git.Groups.ListGroupProjects(group.ID, &gitlab.ListGroupProjectsOptions{
				ListOptions: gitlab.ListOptions{
					Page:    1,
					PerPage: 100,
				},
			})
			if err != nil {
				log.Fatal(err)
			}

			allProjects = append(allProjects, projects...)

			// paging
			if resp.NextPage == 0 {
				break
			}
		}
	}

	sort.Slice(allProjects, func(i, j int) bool { // sort projects alphabetically
		return allProjects[i].Namespace.FullPath < allProjects[j].Namespace.FullPath
	})

	// GitLab のプロジェクト情報を CSV へ出力
	filename := fmt.Sprintf("output/gitlab_projects_%s.csv", gitlabGroupId)
	outputGitlabProjectsToCsv(filename, allProjects)
}

func outputGitlabProjectsToCsv(filename string, projects []*gitlab.Project) {
	// Create CSV file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cw := csv.NewWriter(file)
	defer cw.Flush()

	// GitLab プロジェクト情報の CSV ファイルへの書き込み
	// write header
	cw.Write([]string{
		"Project Name (Repository Name)",
		"Group Name",
		"Group Depth",
		"Last Activity At (JST)",
		"Created At (JST)",
		"Project URL",
	})
	// write data
	for _, project := range projects {
		cw.Write([]string{
			project.Namespace.FullPath, // グループ名、サブグループ名
			fmt.Sprintf("%d", len(strings.Split(project.Namespace.FullPath, "/"))), // グループ、サブグループの階層の深さ
			project.PathWithNamespace,             // プロジェクト名（リポジトリ名）
			TimeToJSTString(*project.LastActivityAt), // 最終更新日時
			TimeToJSTString(*project.CreatedAt),      // 作成日時
			project.WebURL,                        // プロジェクト URL
		})
	}
}

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"
)

// GitLab API クライアントの作成
func createGitlabClient(gitlabToken string) (*gitlab.Client, error) {
	cli, err := gitlab.NewClient(gitlabToken)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// 特定のグループの全ての Gitlab プロジェクトを取得// 特定のグループの全ての Gitlab プロジェクトを取得
func fetchAllGitlabProjectsInGroup(gitlabToken string, gitlabGroupId string) ([]*gitlab.Project, error) {
	// GitLab API クライアントの作成
	cli, err := createGitlabClient(gitlabToken)
	if err != nil {
		return nil, err
	}

	// 大元のグループの取得
	log.Printf("Getting the group info of \"%s\"...\n", gitlabGroupId)
	groupOrigin, _, err := cli.Groups.GetGroup(gitlabGroupId, &gitlab.GetGroupOptions{})
	if err != nil {
		return nil, err
	}
	log.Printf("Group ID   : %d\n", groupOrigin.ID)
	log.Printf("Group Name : %s\n", groupOrigin.Name)

	// 全てのグループ、サブグループの取得
	log.Printf("Getting all subgroups info in \"%s\"...\n", groupOrigin.Name)

	var allGroups []*gitlab.Group
	for page := 1; ; page++ {
		groups, resp, err := cli.Groups.ListDescendantGroups(gitlabGroupId, &gitlab.ListDescendantGroupsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		if err != nil {
			return nil, err
		}

		allGroups = append(allGroups, groups...)

		// paging
		if resp.NextPage == 0 {
			break
		}
	}

	allGroups = append([]*gitlab.Group{groupOrigin}, allGroups...) // 大元のグループを先頭に追加

	// 各グループ、各サブグループのプロジェクト情報の取得
	log.Printf("Getting all projects in every group of \"%s\"...\n", groupOrigin.Name)

	var allProjects []*gitlab.Project
	for _, group := range allGroups {
		log.Printf("Group / Subgroup Name: \"%s\"\n", group.FullPath)

		for page := 1; ; page++ {
			projects, resp, err := cli.Groups.ListGroupProjects(group.ID, &gitlab.ListGroupProjectsOptions{
				ListOptions: gitlab.ListOptions{
					Page:    page,
					PerPage: 100,
				},
			})
			if err != nil {
				return nil, err
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

	return allProjects, nil
}

// GitLab のプロジェクト情報を CSV へ出力
func outputGitlabProjectsToCsv(gitlabGroupId string, projects []*gitlab.Project) (error) {
	log.Printf("Writing GitLab projects to CSV file...\n")

	// Create filename
	now := time.Now()
	filename := fmt.Sprintf("output/gitlab_projects_%s_%s.csv", gitlabGroupId, timeToString(now, layoutForFilename))

	// Create CSV file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	cw := csv.NewWriter(file)
	defer cw.Flush()

	// GitLab プロジェクト情報の CSV ファイルへの書き込み
	// write header
	cw.Write([]string{
		"Project Name (Repository Name)",
		"Group / Subgroup Name",
		"Group Depth",
		"Last Activity At (JST)",
		"Created At (JST)",
		"Project URL",
	})
	// write data
	for _, project := range projects {
		cw.Write([]string{
			project.PathWithNamespace,             // プロジェクト名（リポジトリ名）
			fmt.Sprintf("%d", len(strings.Split(project.Namespace.FullPath, "/"))), // グループ、サブグループの階層の深さ
			project.Namespace.FullPath, // グループ名、サブグループ名
			timeToJSTString(*project.LastActivityAt, layoutDefault), // 最終更新日時
			timeToJSTString(*project.CreatedAt, layoutDefault),      // 作成日時
			project.WebURL,                        // プロジェクト URL
		})
	}

	log.Printf("Complete writing GitLab projects of the group \"%s\" to CSV file 🎉\n", gitlabGroupId)
	log.Printf("Filename: %s\n", filename)

	return nil
}

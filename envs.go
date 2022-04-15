package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Envs struct {
	GITLAB_TOKEN string
	GITLAB_GROUP_ID     string
}

func loadEnvs() *Envs {
	// load .env
	filename := ".env"
	err := godotenv.Load(filename)
	if err != nil {
		log.Fatalf("Cannot load %s: %v", filename, err)
	}

	return &Envs{
		GITLAB_TOKEN: getGitLabToken(),
		GITLAB_GROUP_ID:     getGitlabGroupId(),
	}
}

func getGitLabToken() string {
	GITLAB_TOKEN := os.Getenv("GITLAB_TOKEN")
	if GITLAB_TOKEN == "" {
		log.Fatal("GITLAB_TOKEN is not set")
	}

	return GITLAB_TOKEN
}

func getGitlabGroupId() string {
	GITLAB_GROUP_ID := os.Getenv("GITLAB_GROUP_ID")
	if GITLAB_GROUP_ID == "" {
		log.Fatal("GITLAB_GROUP_ID is not set")
	}

	return GITLAB_GROUP_ID
}

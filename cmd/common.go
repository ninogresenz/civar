package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func getToken() string {
	token = viper.GetString("token")
	if token == "" {
		token = viper.GetString("GITLAB_TOKEN")
	}
	if token == "" {
		fmt.Printf("No Gitlab token found. You can set it via 3 options:\n" +
			"* set flag '--token xxx'\n" +
			"* export GITLAB_TOKEN=xxx\n" +
			"* set token property in $HOME/.civar.yml\n")
		os.Exit(1)
	}
	return token
}

func getGitlabUrl() string {
	url := viper.GetString("url")
	if url == "" {
		url = viper.GetString("GITLAB_URL")
	}
	if url == "" {
		fmt.Printf("No Gitlab url found. You can set it via 3 options:" +
			"* set flag '--url xxx'" +
			"* export GITLAB_URL=xxx" +
			"* set url property in $HOME/.civar.yml")
		os.Exit(1)
	}
	if url[len(url)-1:] == "/" {
		return url[:len(url)-1]
	}
	return url
}

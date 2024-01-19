package cmd

import (
	"net/http"

	"github.com/spf13/cobra"

	"github.com/ninogresenz/civar/gitlab"
	"github.com/ninogresenz/civar/service"
)

// getCmd represents the get command
var searchCmd = &cobra.Command{
	Use:     "search [term]",
	Example: "civar search my-project ",
	Short:   "Searches Gitlab for project names",
	Long:    "Prints all Gitlab Projects for the given search term to stdout.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := gitlab.New(getGitlabUrl(), getToken(), http.DefaultClient)
		service := service.NewService(api, cmd, args)
		service.Search()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

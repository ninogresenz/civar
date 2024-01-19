package cmd

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ninogresenz/civar/gitlab"
	"github.com/ninogresenz/civar/service"
)

var createCmd = &cobra.Command{
	Use:     "create group/project",
	Example: "cat .env | civar create group/project -d",
	Short:   "Creates CI/CD variables",
	Long:    "Reads data from stdin or file and creates all variables in a Gitlab project. Already existent variables will be skipped.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := gitlab.New(getGitlabUrl(), getToken(), http.DefaultClient)
		service := service.NewService(api, cmd, args)
		service.Create(format, k8s, fileFlag)
	},
}

func init() {
	createCmd.Flags().StringVarP(&format, "format", "f", "dotenv", "format is one of [ json | dotenv ]")
	_ = viper.BindPFlag("format", createCmd.Flags().Lookup("format"))
	createCmd.Flags().BoolVarP(&k8s, "k8s", "k", false, "creates variables with K8S_SECRET_ prefix")
	createCmd.Flags().StringVarP(&fileFlag, "file", "F", "", "reads input from a file")
	rootCmd.AddCommand(createCmd)
}

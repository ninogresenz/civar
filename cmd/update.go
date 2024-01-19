package cmd

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ninogresenz/civar/gitlab"
	"github.com/ninogresenz/civar/service"
)

var updateCmd = &cobra.Command{
	Use:     "update group/project",
	Example: "cat .env | civar update group/project -d",
	Short:   "Updates CI/CD variables",
	Long:    "Reads data from stdin or file and updates already existing variables in a Gitlab project. Non existent variables will be skipped.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := gitlab.New(getGitlabUrl(), getToken(), http.DefaultClient)
		service := service.NewService(api, cmd, args)
		service.Update(format, k8s, fileFlag)
	},
}

func init() {
	updateCmd.Flags().StringVarP(&format, "format", "f", "dotenv", "format is one of [ json | dotenv | pretty ]")
	_ = viper.BindPFlag("format", createCmd.Flags().Lookup("format"))
	updateCmd.Flags().BoolVarP(&k8s, "k8s", "k", false, "update variables with K8S_SECRET_ prefix")
	updateCmd.Flags().StringVarP(&fileFlag, "file", "F", "", "reads input from a file")
	rootCmd.AddCommand(updateCmd)
}

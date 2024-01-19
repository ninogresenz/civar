package cmd

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ninogresenz/civar/gitlab"
	"github.com/ninogresenz/civar/service"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "get group/project",
	Example: "civar get group/project -d",
	Short:   "Shows CI/CD variables",
	Long:    "Prints CI/CD variables for the given Gitlab project to stdout.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := gitlab.New(getGitlabUrl(), getToken(), http.DefaultClient)
		service := service.NewService(api, cmd, args)
		if pretty {
			format = "pretty"
		}
		if dotenv {
			format = "dotenv"
		}
		service.Get(format, scopeFilter)
	},
}

func init() {
	getCmd.Flags().StringVarP(&scopeFilter, "scope", "s", "", "scope filter [ * | staging | production | prodtest ]")

	getCmd.Flags().StringVarP(&format, "format", "f", "dotenv", "format is one of [ json | dotenv | pretty ]")
	_ = viper.BindPFlag("format", getCmd.Flags().Lookup("format"))

	// TODO therse should become format
	getCmd.Flags().BoolVarP(&pretty, "pretty", "p", false, "alias for --format pretty")
	getCmd.Flags().BoolVarP(&dotenv, "dotenv", "d", false, "alias for --format dotenv")

	// make sure they're not used together
	getCmd.MarkFlagsMutuallyExclusive(
		"format",
		"pretty",
		"dotenv",
	)
	rootCmd.AddCommand(getCmd)
}

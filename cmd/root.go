package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var token string
var gitlabUrl string

var format string
var pretty bool
var dotenv bool
var scopeFilter string
var k8s bool
var fileFlag string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "civar",
	Short: "Show / Create CI/CD Variables in Gitlab projects",
	Long:  `CLI tool for fetching and creating CI/CD Variables in Gitlab projects`,
	Example: `- Print all CI/CD variables in a table
	civar get -p 1
	
- Copy all CI/CD variables from one project to another:
	civar get 1 | civar create 2
	
- Save output to a file:
	civar get 1 > vars.txt
	
- Create CI/CD variables from a file:
	cat vars.txt | civar create 2
`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.civar.yml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "sets your token")
	_ = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))

	rootCmd.PersistentFlags().StringVarP(&gitlabUrl, "url", "u", "", "sets your gitlab url")
	_ = viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".civar" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".civar")
		viper.SetConfigType("yml")
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error reading config file:", viper.ConfigFileUsed())
	}
}

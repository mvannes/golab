package cmd

import (
	"log"
	"os"

	mergerequest "github.com/mvannes/golab/cmd/merge_request"
	"github.com/mvannes/golab/cmd/project"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "golab",
	Short: "Tooling for gitlab requests",
	Long:  `Golab exposes functionality to monitor your gitlab review requests, among other gitlab related functionality.`,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	project.Add(rootCmd)
	mergerequest.Add(rootCmd)
	dir, err := os.UserHomeDir()
	if nil != err {
		log.Fatal(err)
	}
	viper.SetConfigName("golab-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(dir)

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

}

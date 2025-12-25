/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configFile = "./config.yaml" //"$HOME/.config/go-onedrive/config.yaml"
)

var graphClient *msgraphsdkgo.GraphServiceClient

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-onedrive",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := readConfig(cmd); err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}
		return initializeGraphClient()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return writeConfig(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.new.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

func readConfig(_ *cobra.Command) error {

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	return nil
}

func writeConfig(_ *cobra.Command) error {
	return viper.WriteConfig()
}

func initializeGraphClient() error {
	var authConfig config.AuthenticationConfigImpl
	var localClient *msgraphsdkgo.GraphServiceClient
	var err error

	viperConfig := viper.Sub("auth")

	if err = viperConfig.Unmarshal(&authConfig); err != nil {
		return fmt.Errorf("failed to unmarshal auth config: %w", err)
	}

	if localClient, err = ClientFactory(&authConfig); err != nil {
		return fmt.Errorf("failed to create graph client: %w", err)
	}

	graphClient = localClient

	return nil
}

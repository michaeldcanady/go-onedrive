/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		credential, err := getCredential(cmd)
		if err != nil {
			fmt.Println("Error getting credential:", err)
			return
		}

		if authenticator, ok := credential.(Authenticator); ok {
			record, err := authenticator.Authenticate(context.Background(), &policy.TokenRequestOptions{
				Scopes: []string{
					"Files.ReadWrite.All",
					//"Sites.ReadWrite.All",
					//"offline_access",
					"User.Read",
				},
			})
			if err != nil {
				fmt.Println("Error authenticating:", err)
				return
			}
			fmt.Println(record)
		}
	},
}

func init() {
	authCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getCredential(_ *cobra.Command) (azcore.TokenCredential, error) {
	var authConfig config.AuthenticationConfigImpl
	var err error

	viperConfig := viper.Sub("auth")

	if err = viperConfig.Unmarshal(&authConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal auth config: %w", err)
	}

	return CredentialFactory(&authConfig)
}

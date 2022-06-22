package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/CheckPointSW/infinity-next-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type discardResponseData struct {
	DiscardChanges bool `json:"discardChanges"`
}

// discardCmd represents the discard command
var discardCmd = &cobra.Command{
	Use:   "discard",
	Short: "Discard changes of a session",
	Long:  `Discard changes of a session`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().StringVarP(&clientID, "client-id", "c", "", "Client ID of the API key")
		cmd.Flags().StringVarP(&accessKey, "access-key", "k", "", "Access key of the API key")
		cmd.Flags().StringVarP(&region, "region", "r", "eu", "Region of Infinity Next API")
		cmd.Flags().StringVarP(&token, "token", "t", "", "Authorization token of the API key")

		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
				if err := cmd.Flags().Set(f.Name, viper.GetString(f.Name)); err != nil {
					fmt.Println(err)
				}
			}
		})

		if err := cmd.MarkFlagRequired("client-id"); err != nil {
			return err
		}

		if err := cmd.MarkFlagRequired("access-key"); err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var URL string
		var API string
		switch region {
		case "eu":
			URL = EUCIURL
			API = CIAPIV1
		case "us":
			URL = USCIURL
			API = CIAPIV1
		case "dev":
			URL = DevCIURL
			API = DevCIAPIV1
		case "preprod":
			URL = DevCIURL
			API = CIAPIV1
		default:
			fmt.Printf("Invalid region %s, expected eu or us\n", region)
			os.Exit(1)
		}

		authForm := url.Values{}
		authForm.Add("clientId", clientID)
		authForm.Add("accessKey", accessKey)
		authReq, err := http.NewRequest(http.MethodPost, URL+CIAuthPath, strings.NewReader(authForm.Encode()))
		if err != nil {
			return err
		}

		authReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		client := http.Client{
			Timeout: 1 * time.Minute,
		}
		authResp, err := client.Do(authReq)
		if err != nil {
			return err
		}

		if authResp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed authenticating to Infinity Next with status %s\n", authResp.Status)
		}

		bResp, err := io.ReadAll(authResp.Body)
		if err != nil {
			return err
		}

		var auth externalAuthResponse
		if err := json.Unmarshal(bResp, &auth); err != nil {
			return err
		}

		var graphqlReq graphqlRequest
		graphqlReq.Query = `
		mutation discardChanges{
			discardChanges
		}
 `
		graphqlReq.Variables = map[string]interface{}{}

		bReq, err := json.Marshal(graphqlReq)
		if err != nil {
			return err
		}

		enforceReq, err := http.NewRequest(http.MethodPost, URL+API, bytes.NewBuffer(bReq))
		if err != nil {
			return err
		}

		enforceReq.Header.Set("Authorization", "Bearer "+auth.Data.Token)
		enforceReq.Header.Set("Content-Type", "application/json")

		var discardChanges graphqlResponse[discardResponseData]
		discardResp, err := utils.HTTPRequestUnmarshal(&client, enforceReq, &discardChanges)
		if err != nil {
			return err
		}

		if discardResp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed discarding changes with status %s\n", discardResp.Status)
		}

		if !discardChanges.Data.DiscardChanges {
			return errors.New("failed discarding changes")
		}

		fmt.Println("Successfully discarded changes")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(discardCmd)
}

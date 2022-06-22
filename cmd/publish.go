package main

import (
	"bytes"
	"encoding/json"
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

type validationMsg struct {
	Message string `json:"message"`
}

type publishResponsePublishChanges struct {
	IsValid  bool            `json:"isValid"`
	Errors   []validationMsg `json:"errors"`
	Warnings []validationMsg `json:"warnings"`
}

type publishResponseData struct {
	PublishChanges publishResponsePublishChanges `json:"publishChanges"`
}

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish changes of a session",
	Long:  `Publish changes of a session`,
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
			API = ProdCIAPIV1
		case "us":
			URL = USCIURL
			API = ProdCIAPIV1
		case "dev":
			URL = DevCIURL
			API = DevCIAPIV1
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
		mutation publishChanges {
		 publishChanges {
		   isValid
		   errors {
			 message
		   }
		   warnings {
			 message
		   }
		 }
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

		var publishChanges graphqlResponse[publishResponseData]
		publishResp, err := utils.HTTPRequestUnmarshal(&client, enforceReq, &publishChanges)
		if err != nil {
			return err
		}

		if publishResp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed publishing changes with status %s\n", publishResp.Status)
		}

		if !publishChanges.Data.PublishChanges.IsValid {
			return fmt.Errorf("failed publishing changes with errors: %s", strings.Join(utils.Map(publishChanges.Data.PublishChanges.Errors, func(t validationMsg) string {
				return t.Message
			}), ", "))
		}

		if len(publishChanges.Data.PublishChanges.Warnings) > 0 {
			fmt.Printf("published changes with warnings: %s\n", strings.Join(utils.Map(publishChanges.Data.PublishChanges.Warnings, func(t validationMsg) string {
				return t.Message
			}), ", "))

			return nil
		}

		fmt.Println("Successfully published changes")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
}

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

type enforcePolicyResponseEnforcePolicy struct {
	ID string `json:"id"`
}

type enforcePolicyResponseData struct {
	EnforcePolicy enforcePolicyResponseEnforcePolicy `json:"enforcePolicy"`
}

// enforceCmd represents the enforce command
var enforceCmd = &cobra.Command{
	Use:   "enforce",
	Short: "Enforce a policy",
	Long:  `Enforce a policy`,

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
		switch region {
		case "eu":
			URL = EUCIURL
		case "us":
			URL = USCIURL
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
		client := http.Client{}
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
		graphqlReq.Query = `mutation {enforcePolicy {id}}`
		graphqlReq.Variables = map[string]interface{}{}

		bReq, err := json.Marshal(graphqlReq)
		if err != nil {
			return err
		}

		enforceReq, err := http.NewRequest(http.MethodPost, URL+CIAPIV1, bytes.NewBuffer(bReq))
		if err != nil {
			return err
		}

		enforceReq.Header.Set("Authorization", "Bearer "+auth.Data.Token)
		enforceReq.Header.Set("Content-Type", "application/json")

		var enforcePolicy graphqlResponse[enforcePolicyResponseData]
		enforceResp, err := utils.HTTPRequestUnmarshal(&client, enforceReq, &enforcePolicy)
		if err != nil {
			return err
		}

		if enforceResp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed enforcing policy with status %s\n", enforceResp.Status)
		}

		taskStatus := "InProgress"
		graphqlReq.Query = "query {getTask(id: \"" + enforcePolicy.Data.EnforcePolicy.ID + "\") {id\r\nstatus}}"
		bReq, err = json.Marshal(graphqlReq)
		if err != nil {
			return err
		}

		errch := make(chan error, 1)
		go func() {
			// Poll for the enforce policy task's status until it's done
			for taskStatus == "InProgress" {
				taskReq, err := http.NewRequest(http.MethodPost, URL+CIAPIV1, bytes.NewBuffer(bReq))
				if err != nil {
					errch <- err
				}

				taskReq.Header.Set("Authorization", "Bearer "+auth.Data.Token)
				taskReq.Header.Set("Content-Type", "application/json")

				var getTask getTaskResponse
				taskResp, err := utils.HTTPRequestUnmarshal(&client, taskReq, &getTask)
				if err != nil {
					errch <- err
				}

				if taskResp.StatusCode != http.StatusOK {
					errch <- fmt.Errorf("failed to get task %s with status %s\n", enforcePolicy.Data.EnforcePolicy.ID, enforceResp.Status)
				}

				taskStatus = getTask.Data.GetTask.Status
				time.Sleep(200 * time.Millisecond)
			}

			errch <- nil
		}()

		// Timeout of 10 seconds for the polling routine to finish
		select {
		case err := <-errch:
			if err != nil {
				return err
			}

			switch taskStatus {
			case "Succeeded":
				fmt.Printf("Enforce policy task %s succeeded\n", enforcePolicy.Data.EnforcePolicy.ID)
			case "Failed":
				return fmt.Errorf("enforce policy task %s failed", enforcePolicy.Data.EnforcePolicy.ID)
			default:
				return fmt.Errorf("enforce policy task %s done with unknown status %s", enforcePolicy.Data.EnforcePolicy.ID, taskStatus)
			}
		case <-time.After(10 * time.Second):
			return fmt.Errorf("enforce policy task did not finish after 10 seconds, quiting")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(enforceCmd)
}

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
	"github.com/golang-jwt/jwt/v5"
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
		API := policyPath
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

		discardReq, err := http.NewRequest(http.MethodPost, URL+API, bytes.NewBuffer(bReq))
		if err != nil {
			return err
		}

		token, _, err := jwt.NewParser().ParseUnverified(auth.Data.Token, jwt.MapClaims{})
		if err != nil {
			return fmt.Errorf("failed to parse token: %w", err)
		}

		tokenMapClaims := token.Claims.(jwt.MapClaims)
		if appID, ok := tokenMapClaims[appIDClaim]; ok {
			switch appID.(string) {
			case wafAppID:
				if API != wafPath {
					API = wafPath
					discardReq, err = http.NewRequest(http.MethodPost, URL+API, bytes.NewBuffer(bReq))
					if err != nil {
						return err
					}
				}
			case policyAppID:
				if API != policyPath {
					API = policyPath
					discardReq, err = http.NewRequest(http.MethodPost, URL+API, bytes.NewBuffer(bReq))
					if err != nil {
						return err
					}
				}
			}
		}

		discardReq.Header.Set("Authorization", "Bearer "+auth.Data.Token)
		discardReq.Header.Set("Content-Type", "application/json")

		var discardChanges graphqlResponse[discardResponseData]
		discardResp, err := utils.HTTPRequestUnmarshal(&client, discardReq, &discardChanges)
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

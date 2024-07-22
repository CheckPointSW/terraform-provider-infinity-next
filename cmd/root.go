package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	viperpkg "github.com/spf13/viper"
)

const (
	EUCIURL     = "https://cloudinfra-gw.portal.checkpoint.com"
	USCIURL     = "https://cloudinfra-gw-us.portal.checkpoint.com"
	DevCIURL    = "https://dev-cloudinfra-gw.kube1.iaas.checkpoint.com"
	CIAuthPath  = "/auth/external"
	DevCIAPIV1  = "/app/infinity2gem/graphql/V1"
	appIDClaim  = "appId"
	wafAppID    = "64488de9-f813-42a7-93e7-f3fe25dd9011"
	policyAppID = "f47b536c-a990-42fb-9ab2-ec38f8c2dcff"
	wafPath     = "/app/waf/graphql/V1"
	policyPath  = "/app/i2/graphql/V1"
)

var (
	clientID  string
	accessKey string
	region    string
	token     string
	cfgFile   string
)

type graphqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type graphqlResponse[T any] struct {
	Data T `json:"data"`
}

type externalAuthResponseData struct {
	Token     string `json:"token"`
	CSRF      string `json:"csrf"`
	Expires   string `json:"expires"`
	ExpiresIn int    `json:"expiresIn"`
}

type externalAuthResponse struct {
	Success bool                     `json:"success"`
	Data    externalAuthResponseData `json:"data"`
}

type lowerCaseStringEnvKeyReplacer struct{}

func (r *lowerCaseStringEnvKeyReplacer) Replace(key string) string {
	return strings.ReplaceAll(key, "-", "_")
}

type getTaskResponseGetTask struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type getTaskResponseData struct {
	GetTask getTaskResponseGetTask `json:"getTask"`
}

type getTaskResponse struct {
	Data getTaskResponseData `json:"data"`
}

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "inext",
		Short: "Infinity Next API Command Line Interface",
		Long: `Infinity Next API Command Line Interface
For example:
inext publish && inext enforce
`,
	}

	viper = viperpkg.NewWithOptions(viperpkg.EnvKeyReplacer(&lowerCaseStringEnvKeyReplacer{}))
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.inext.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".inext" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".inext")
	}

	viper.SetEnvPrefix("inext")
	viper.AutomaticEnv() // read in environment variables that match INEXT_*

	// If a config file is found, read it in.
	viper.ReadInConfig()
}

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/CheckPointSW/terraform-provider-infinity-next/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

// Same us plugin.Debug with simpler output to parse
func debug(ctx context.Context, providerAddr string, opts *plugin.ServeOpts) error {
	ctx, cancel := context.WithCancel(ctx)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer func() {
		signal.Stop(sigCh)
		cancel()
	}()
	config, closeCh, err := plugin.DebugServe(ctx, opts)
	if err != nil {
		return fmt.Errorf("Error launching debug server: %w", err)
	}
	go func() {
		select {
		case <-sigCh:
			cancel()
		case <-ctx.Done():
		}
	}()
	reattachBytes, err := json.Marshal(map[string]plugin.ReattachConfig{
		providerAddr: config,
	})
	if err != nil {
		return fmt.Errorf("Error building reattach string: %w", err)
	}

	reattachStr := string(reattachBytes)

	switch runtime.GOOS {
	case "windows":
		fmt.Printf("Command Prompt:\tset \"TF_REATTACH_PROVIDERS=%s\"\n", reattachStr)
		fmt.Printf("PowerShell:\t$env:TF_REATTACH_PROVIDERS='%s'\n", strings.ReplaceAll(reattachStr, `'`, `''`))
	case "linux", "darwin":
		fmt.Printf("TF_REATTACH_PROVIDERS='%s'\n", strings.ReplaceAll(reattachStr, `'`, `'"'"'`))
	default:
		fmt.Println(reattachStr)
	}
	fmt.Println("")

	// wait for the server to be done
	<-closeCh
	return nil
}

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: provider.Provider}

	if debugMode {
		err := debug(context.Background(), "checkpointsw/infinity-next", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}

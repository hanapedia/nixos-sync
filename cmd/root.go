package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

type Config struct {
	Hosts map[string]string `toml:"hosts"`
}

var (
	tomlPath  string
	flakePath string
	hostNames []string
)

var rootCmd = &cobra.Command{
	Use:   "nixos-sync",
	Short: "Sync NixOS configurations to remote hosts via SSH",
	Run: func(cmd *cobra.Command, args []string) {
		// Read and parse TOML
		var cfg Config
		data, err := os.ReadFile(tomlPath)
		if err != nil {
			log.Fatalf("Failed to read %s: %v", tomlPath, err)
		}
		if err := toml.Unmarshal(data, &cfg); err != nil {
			log.Fatalf("Failed to parse TOML: %v", err)
		}

		if len(hostNames) > 0 {
			for _, host := range hostNames {
				target, ok := cfg.Hosts[host]
				if !ok {
					log.Fatalf("Host %s not found in %s", host, tomlPath)
				}
				runRebuild(host, target)
			}
		} else {
			for host, target := range cfg.Hosts {
				runRebuild(host, target)
			}
		}
	},
}

func runRebuild(host, target string) {
	fmt.Printf("Deploying to %s (%s)...\n", host, target)
	cmd := exec.Command("nixos-rebuild", "switch",
		"--flake", fmt.Sprintf("%s#%s", flakePath, host),
		"--target-host", target,
		"--build-host", target,
		"--use-remote-sudo",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		log.Printf("Failed to deploy to %s: %v\n", host, err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&tomlPath, "toml", "t", "", "Path to hosts.toml file (required)")
	rootCmd.Flags().StringVarP(&flakePath, "flake", "f", "", "Path to flake directory (required)")
	rootCmd.Flags().StringArrayVarP(&hostNames, "host", "h", nil, "Hostname(s) to deploy (repeatable)")

	// Mark required flags
	rootCmd.MarkFlagRequired("toml")
	rootCmd.MarkFlagRequired("flake")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package commands

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "cli is a CLI that manages your dream cli.",
	Long:  `A CLI that keeps track of the valuation of your cli, and other smart home automation.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal(err)
		}
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute(version string) {
	rootCmd.AddCommand(versionCmd(version))
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func versionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "View the installed version of crib",
		Long:  "View the installed version of crib",
		Run: func(cmd *cobra.Command, args []string) {
			pterm.DefaultBasicText.WithWriter(os.Stderr).Printfln("version: %v", version)
		},
	}
}

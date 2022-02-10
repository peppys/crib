package commands

import (
	"fmt"
	"github.com/leekchan/accounting"
	"github.com/peppys/crib/internal/redfin"
	"github.com/peppys/crib/internal/zillow"
	"github.com/peppys/crib/pkg/crib"
	"github.com/peppys/crib/pkg/crib/estimators"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
)

var valueCmd = &cobra.Command{
	Use:   "value",
	Short: "Checks the estimated valuation of your property",
	Long:  "Checks the estimated valuation of your property",
	PreRun: func(cmd *cobra.Command, args []string) {
		if strings.TrimSpace(address) == "" {
			pterm.Error.Println("address must be provided")
			cmd.Help()
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = pterm.DefaultBigText.WithLetters(pterm.NewLettersFromStringWithStyle("crib", pterm.NewStyle(pterm.FgLightMagenta))).Render()
		introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(true).WithRemoveWhenDone(true).Start(fmt.Sprintf("estimating valuation for crib '%s'...", address))

		estimates, err := theCrib.Estimate(address)
		introSpinner.Stop()
		if err != nil {
			pterm.Error.Println(err)
			os.Exit(1)
		}
		data := pterm.TableData{
			{"Vendor", "Estimate"},
		}
		ac := accounting.Accounting{Symbol: "$", Precision: 2}
		for _, estimate := range estimates {
			data = append(data, []string{strings.ToLower(string(estimate.Vendor)), ac.FormatMoney(estimate.Value)})
		}
		pterm.DefaultTable.WithHasHeader().WithData(data).Render()
	},
}

var theCrib *crib.Crib
var address string

func init() {
	valueCmd.Flags().StringVarP(&address, "address", "a", "", "Address of your crib")
	valueCmd.MarkFlagRequired("address")
	rootCmd.AddCommand(valueCmd)
	theCrib = crib.NewCrib(
		crib.WithEstimators(
			estimators.NewZillowEstimator(zillow.NewClient(http.DefaultClient)),
			estimators.NewRedfinEstimator(redfin.NewClient(http.DefaultClient)),
		),
	)
}

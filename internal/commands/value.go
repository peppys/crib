package commands

import (
	"fmt"
	"github.com/leekchan/accounting"
	"github.com/peppys/crib/internal/services/property"
	"github.com/peppys/crib/internal/services/property/estimators"
	"github.com/peppys/crib/pkg/redfin"
	"github.com/peppys/crib/pkg/zillow"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strings"
)

var valueCmd = &cobra.Command{
	Use:   "value",
	Short: "Checks the estimated valuation of your property",
	Long:  "Checks the estimated valuation of your property",
	Run: func(cmd *cobra.Command, args []string) {
		_ = pterm.DefaultBigText.WithLetters(pterm.NewLettersFromStringWithStyle("crib", pterm.NewStyle(pterm.FgLightMagenta))).Render()
		introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(true).WithRemoveWhenDone(true).Start(fmt.Sprintf("estimating valuation for crib '%s'...", address))
		estimates, err := manager.Valuation(address)
		introSpinner.Stop()
		if err != nil {
			pterm.Error.Println(err)
			log.Fatal(err)
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

var manager *property.Manager
var address string

func init() {
	valueCmd.Flags().StringVarP(&address, "address", "a", "", "Address of your crib")
	valueCmd.MarkFlagRequired("address")
	rootCmd.AddCommand(valueCmd)
	manager = property.NewManager(
		property.WithEstimator(
			estimators.NewBulkEstimator(
				estimators.NewZillowEstimator(zillow.NewClient(http.DefaultClient)),
				estimators.NewRedfinEstimator(redfin.NewClient(http.DefaultClient)),
			),
		),
	)
}

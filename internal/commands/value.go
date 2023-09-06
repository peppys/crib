package commands

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/leekchan/accounting"
	"github.com/peppys/crib/internal/redfin"
	"github.com/peppys/crib/internal/zillow"
	"github.com/peppys/crib/pkg/crib"
	"github.com/peppys/crib/pkg/crib/estimators"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var valueCmd = &cobra.Command{
	Use:   "value",
	Short: "Checks the estimated valuation of your property",
	Long:  "Checks the estimated valuation of your property",
	PreRun: func(cmd *cobra.Command, args []string) {
		if strings.TrimSpace(address) == "" {
			pterm.Error.WithWriter(os.Stderr).Println("address must be provided")
			cmd.Help()
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !slices.Contains(validFormats, format) {
			pterm.Error.WithWriter(os.Stderr).Println(fmt.Sprintf("format must be one of: [%v]", strings.Join(validFormats, ", ")))
			os.Exit(1)
		}
		_ = pterm.DefaultBigText.WithWriter(os.Stderr).WithLetters(putils.LettersFromStringWithStyle("crib", pterm.NewStyle(pterm.FgLightMagenta))).Render()
		introSpinner, _ := pterm.DefaultSpinner.WithWriter(os.Stderr).WithShowTimer(true).WithRemoveWhenDone(true).Start(fmt.Sprintf("estimating valuation for crib '%s'...", address))

		estimates, err := theCrib.Estimate(address)
		introSpinner.Stop()
		if err != nil {
			pterm.Error.WithWriter(os.Stderr).Println(err)
			os.Exit(1)
		}

		switch format {
		case jsonFormat:
			encoded, err := json.MarshalIndent(estimates, "", "  ")
			if err != nil {
				pterm.Error.WithWriter(os.Stderr).Println(err)
				os.Exit(1)
			}

			fmt.Println(string(encoded))
		case csvFormat:
			writer := csv.NewWriter(os.Stdout)
			records := [][]string{
				{"Provider", "Estimate"},
			}
			for _, estimate := range estimates {
				records = append(records, []string{
					string(estimate.Provider),
					strconv.FormatInt(estimate.Value, 10),
				})
			}
			writer.WriteAll(records)
			writer.Flush()
		case tableFormat:
			fallthrough
		default:
			data := pterm.TableData{
				{"Provider", "Estimate"},
			}
			ac := accounting.Accounting{Symbol: "$", Precision: 2}
			for _, estimate := range estimates {
				data = append(data, []string{strings.ToLower(string(estimate.Provider)), ac.FormatMoneyBigRat(big.NewRat(estimate.Value, 100))})
			}
			pterm.DefaultTable.WithWriter(os.Stderr).WithHasHeader().WithData(data).Render()
		}
	},
}

var theCrib *crib.Crib
var address string
var format string

const (
	tableFormat string = "table"
	jsonFormat         = "json"
	csvFormat          = "csv"
)

var validFormats = []string{tableFormat, jsonFormat, csvFormat}

func init() {
	valueCmd.Flags().StringVarP(&address, "address", "a", "", "Address of your crib")
	valueCmd.Flags().StringVarP(&format, "format", "f", tableFormat, "Format of output (table/json)")
	valueCmd.MarkFlagRequired("address")
	rootCmd.AddCommand(valueCmd)
	theCrib = crib.New(
		crib.WithEstimators(
			estimators.NewZillowEstimator(zillow.NewClient(http.DefaultClient)),
			estimators.NewRedfinEstimator(redfin.NewClient(http.DefaultClient)),
		),
	)
}

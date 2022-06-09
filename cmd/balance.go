/*
Copyright © 2022 kubetrail.io

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/kubetrail/ethio/pkg/flags"
	"github.com/kubetrail/ethio/pkg/run"
	"github.com/spf13/cobra"
)

// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Check address balance",
	RunE:  run.Balance,
}

func init() {
	rootCmd.AddCommand(balanceCmd)
	f := balanceCmd.Flags()

	f.String(flags.Addr, "", "Address")
	f.Int64(flags.BlockNumber, -1, "Block number (-1 = current)")
	f.String(flags.Unit, flags.UnitEth, "Amount unit (eth, wei or gwei)")

	_ = balanceCmd.RegisterFlagCompletionFunc(
		flags.Unit,
		func(
			cmd *cobra.Command,
			args []string,
			toComplete string,
		) (
			[]string,
			cobra.ShellCompDirective,
		) {
			return []string{
					flags.UnitEth,
					flags.UnitWei,
					flags.UnitGwei,
				},
				cobra.ShellCompDirectiveDefault
		},
	)
}

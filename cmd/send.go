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

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send ether",
	Long:  ``,
	RunE:  run.Send,
}

func init() {
	rootCmd.AddCommand(sendCmd)
	f := sendCmd.Flags()

	f.String(flags.Addr, "", "Address of the receiver")
	f.String(flags.Key, "", "Private key of sender")
	f.String(flags.Unit, flags.UnitEth, "Amount unit (eth, wei or gwei)")
	f.Float64(flags.Amount, 0, "Amount to send")
	f.Float64(flags.Gas, 30, "Gas price in gwei (-1 for auto set)")

	_ = sendCmd.RegisterFlagCompletionFunc(
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

	_ = sendCmd.RegisterFlagCompletionFunc(
		flags.Gas,
		func(
			cmd *cobra.Command,
			args []string,
			toComplete string,
		) (
			[]string,
			cobra.ShellCompDirective,
		) {
			return []string{
					"30",
					"40",
					"50",
				},
				cobra.ShellCompDirectiveDefault
		},
	)
}

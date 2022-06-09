package run

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kubetrail/bip32/pkg/keys"
	"github.com/kubetrail/bip39/pkg/prompts"
	"github.com/kubetrail/ethio/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Balance(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.Addr, cmd.Flag(flags.Addr))
	_ = viper.BindPFlag(flags.BlockNumber, cmd.Flag(flags.BlockNumber))
	_ = viper.BindPFlag(flags.Unit, cmd.Flag(flags.Unit))
	addr := viper.GetString(flags.Addr)
	blockNumber := viper.GetInt64(flags.BlockNumber)
	unit := strings.ToLower(viper.GetString(flags.Unit))

	prompt, err := prompts.Status()
	if err != nil {
		return fmt.Errorf("failed to get prompt status: %w", err)
	}

	switch unit {
	case flags.UnitEth, flags.UnitWei, flags.UnitGwei:
	default:
		return fmt.Errorf("invalid unit, it can be either eth, wei or gwei")
	}

	switch strings.ToLower(persistentFlags.OutputFormat) {
	case flags.OutputFormatNative, flags.OutputFormatYaml, flags.OutputFormatJson:
	default:
		return fmt.Errorf("invalid output format, it can be native, yaml or json")
	}

	if len(addr) == 0 {
		if len(args) == 0 {
			if prompt {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter address: "); err != nil {
					return fmt.Errorf("failed to write to output: %w", err)
				}
			}
			addr, err = keys.Read(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read pub addr from input: %w", err)
			}
		} else {
			addr = args[0]
		}
	}

	var blockNumberBigInt *big.Int
	if blockNumber > -1 {
		blockNumberBigInt = big.NewInt(blockNumber)
	}

	client, err := ethclient.Dial(persistentFlags.RPCEndpoint)
	if err != nil {
		return fmt.Errorf("failed to dial eth client: %w", err)
	}
	defer client.Close()

	account := common.HexToAddress(addr)
	balance, err := client.BalanceAt(ctx, account, blockNumberBigInt)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	type output[T float64 | int64] struct {
		Amount T      `json:"amount" yaml:"amount"`
		Unit   string `json:"unit" yaml:"unit"`
	}

	switch unit {
	case flags.UnitEth:
		out := &output[float64]{
			Amount: float64(balance.Int64()) / 1e18,
			Unit:   unit,
		}

		switch persistentFlags.OutputFormat {
		case flags.OutputFormatNative:
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), out.Amount); err != nil {
				return fmt.Errorf("failed to print to output: %w", err)
			}
		case flags.OutputFormatYaml:
			b, err := yaml.Marshal(out)
			if err != nil {
				return fmt.Errorf("failed to yaml marshal output: %w", err)
			}
			if _, err := fmt.Fprint(cmd.OutOrStdout(), string(b)); err != nil {
				return fmt.Errorf("failed to print yaml to output: %w", err)
			}
		case flags.OutputFormatJson:
			b, err := json.Marshal(out)
			if err != nil {
				return fmt.Errorf("failed to json marshal output: %w", err)
			}
			if _, err := fmt.Fprint(cmd.OutOrStdout(), string(b)); err != nil {
				return fmt.Errorf("failed to print json to output: %w", err)
			}
		}
	case flags.UnitWei:
		out := &output[int64]{
			Amount: balance.Int64(),
			Unit:   unit,
		}

		switch persistentFlags.OutputFormat {
		case flags.OutputFormatNative:
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), out.Amount); err != nil {
				return fmt.Errorf("failed to print to output: %w", err)
			}
		case flags.OutputFormatYaml:
			b, err := yaml.Marshal(out)
			if err != nil {
				return fmt.Errorf("failed to yaml marshal output: %w", err)
			}
			if _, err := fmt.Fprint(cmd.OutOrStdout(), string(b)); err != nil {
				return fmt.Errorf("failed to print yaml to output: %w", err)
			}
		case flags.OutputFormatJson:
			b, err := json.Marshal(out)
			if err != nil {
				return fmt.Errorf("failed to json marshal output: %w", err)
			}
			if _, err := fmt.Fprint(cmd.OutOrStdout(), string(b)); err != nil {
				return fmt.Errorf("failed to print json to output: %w", err)
			}
		}
	case flags.UnitGwei:
		out := &output[int64]{
			Amount: balance.Int64() / 1000000000,
			Unit:   unit,
		}

		switch persistentFlags.OutputFormat {
		case flags.OutputFormatNative:
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), out.Amount); err != nil {
				return fmt.Errorf("failed to print to output: %w", err)
			}
		case flags.OutputFormatYaml:
			b, err := yaml.Marshal(out)
			if err != nil {
				return fmt.Errorf("failed to yaml marshal output: %w", err)
			}
			if _, err := fmt.Fprint(cmd.OutOrStdout(), string(b)); err != nil {
				return fmt.Errorf("failed to print yaml to output: %w", err)
			}
		case flags.OutputFormatJson:
			b, err := json.Marshal(out)
			if err != nil {
				return fmt.Errorf("failed to json marshal output: %w", err)
			}
			if _, err := fmt.Fprint(cmd.OutOrStdout(), string(b)); err != nil {
				return fmt.Errorf("failed to print json to output: %w", err)
			}
		}
	}

	return nil
}

package run

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kubetrail/bip32/pkg/keys"
	"github.com/kubetrail/bip39/pkg/prompts"
	"github.com/kubetrail/ethio/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	weiPerEth  = 1e18
	weiPerGwei = 1e9
)

func Send(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.Addr, cmd.Flag(flags.Addr))
	_ = viper.BindPFlag(flags.Unit, cmd.Flag(flags.Unit))
	_ = viper.BindPFlag(flags.Key, cmd.Flag(flags.Key))
	_ = viper.BindPFlag(flags.Amount, cmd.Flag(flags.Amount))
	_ = viper.BindPFlag(flags.Gas, cmd.Flag(flags.Gas))
	addr := viper.GetString(flags.Addr)
	key := viper.GetString(flags.Key)
	unit := strings.ToLower(viper.GetString(flags.Unit))
	amount := viper.GetFloat64(flags.Amount)
	gas := viper.GetFloat64(flags.Gas)

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
		if prompt {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter address: "); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}
		addr, err = keys.Read(cmd.InOrStdin())
		if err != nil {
			return fmt.Errorf("failed to read pub addr from input: %w", err)
		}
	}

	if len(key) == 0 {
		if prompt {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter sender private key: "); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}
		key, err = keys.Read(cmd.InOrStdin())
		if err != nil {
			return fmt.Errorf("failed to read key from input: %w", err)
		}
	}

	client, err := ethclient.Dial(persistentFlags.RPCEndpoint)
	if err != nil {
		return fmt.Errorf("failed to dial eth client: %w", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	toAddress := common.HexToAddress(addr)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return fmt.Errorf("failed to gen pending nonce: %w", err)
	}

	var value *big.Int
	var gasPrice *big.Int

	switch unit {
	case flags.UnitEth:
		value = big.NewInt(int64(amount * weiPerEth))
	case flags.UnitWei:
		value = big.NewInt(int64(amount))
	case flags.UnitGwei:
		value = big.NewInt(int64(amount * weiPerGwei))
	}

	gasLimit := uint64(21000) // in units

	if gas >= 0 {
		gasPrice = big.NewInt(int64(gas * weiPerGwei))
	} else {
		gasPrice, err = client.SuggestGasPrice(ctx)
		if err != nil {
			return fmt.Errorf("failed to get suggested gas price: %w", err)
		}
	}

	//tx2 := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasLimit,
			To:       &toAddress,
			Value:    value,
			Data:     nil,
			V:        nil,
			R:        nil,
			S:        nil,
		},
	)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get network ID: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	type output struct {
		FromAddr string `json:"fromAddr,omitempty" yaml:"fromAddr,omitempty"`
		ToAddr   string `json:"toAddr,omitempty" yaml:"toAddr,omitempty"`
		Amount   int64  `json:"amount" yaml:"amount"`
		TxHash   string `json:"txHash,omitempty" yaml:"txHash,omitempty"`
	}

	out := &output{
		FromAddr: fromAddress.String(),
		ToAddr:   toAddress.String(),
		Amount:   value.Int64(),
		TxHash:   signedTx.Hash().Hex(),
	}

	switch persistentFlags.OutputFormat {
	case flags.OutputFormatNative:
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), out.TxHash); err != nil {
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

	return nil
}

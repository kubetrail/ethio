package flags

const (
	RPCEndpoint  = "rpc-endpoint"
	OutputFormat = "output-format"
	Amount       = "amount"
	Key          = "key"
	Addr         = "addr"
	BlockNumber  = "block-number"
	Unit         = "unit"
	Gas          = "gas"
)

const (
	DefaultRPCEndpoint          = "https://goerli.infura.io"
	DefaultRPCEndpointEnvVarKey = "ETHIO_RPC_ENDPOINT"

	OutputFormatNative = "native"
	OutputFormatJson   = "json"
	OutputFormatYaml   = "yaml"

	UnitEth  = "eth"
	UnitWei  = "wei"
	UnitGwei = "gwei"
)

# ethio
Tool to perform Ethereum cryptocurrency transactions, check balances etc.

## disclaimer
> The use of this tool does not guarantee security or usability for any
> particular purpose. Please review the code and use at your own risk.

## installation
This step assumes you have [Go compiler toolchain](https://go.dev/dl/)
installed on your system.

```bash
go install github.com/kubetrail/ethio@latest
```
Add autocompletion for `bash` to your `.bashrc`
```bash
source <(ethio completion bash)
```

## rpc endpoint
In order to send transactions or query balance you will need an API endpoint.
An option is to run a node yourself, other is to leverage endpoint from
[infura.io](https://infura.io/). Below are general guidelines for creating an
endpoint:
* Signup on infura.io
* Create a project
* Select ETH network
* Copy endpoint for it and set an environment variable

```bash
export ETHIO_RPC_ENDPOINT=https://goerli.infura.io/v3/your-account-number
```

## ethereum keys
For the sake of this readme, we will create two keys, one as sender and the
other as receiver. In order to create new keys, you will need to download
[bip39](https://github.com/kubetrail/bip39) and 
[ethkey](https://github.com/kubetrail/ethkey)

### for sender
Generate a new mnemonic sentence:
```bash
bip39 gen --length=12
```
```text
predict cancel fun split recycle expand wise mixed unfold bulb festival fox
```

Generate ethereum keys
```bash
ethkey gen predict cancel fun split recycle expand wise mixed unfold bulb festival fox
```
```yaml
seed: 063015b36ea348db3ecad0469bd125791338bdff18a31a5d3ba97c06c4602cf4126315d9f354fad7f99c45a98eece00e877d3c54825868887646145911433f06
prvHex: 5040f993d59dd4dbba2bbcb0286371bdb02dc0bacc73ce56b989d07e9ef8bb24
pubHex: 8324573c4d4c423c75ab0488abb74790f163d46936ba78dd32a91b4ef7cfc9a5722bc251ddac4c001c33815b01aa94afc87cd0de53cae09aeca601ebaf1cdae6
addr: 0x07B5e1b5fB3746b117241e493b6fD42b4FC74a76
keyType: ecdsa
derivationPath: m/44'/60'/0'/0/0
```

### for receiver
Generate a new mnemonic sentence:
```bash
bip39 gen --length=12
```
```text
leader hen juice nut story shed gentle later vault agree snake that
```

Generate ethereum keys
```bash
ethkey gen leader hen juice nut story shed gentle later vault agree snake that
```
```yaml
seed: bcf18b367fd4204eff5bec890d7ed850588217a9a572a29819ddf5a567db9e4f04c4901f24ab9faac2bf34aa67d0bc5373e8c719b1dab36f01626cbe88ef5c8e
prvHex: 901eb4d266daa3b45a88d41d37e0d8d12897913b4175e4b2bde1acb0c651dced
pubHex: 83ad922e870595644e35f7c5046d24d47b07158835783eee4f616251be202720e5c5ccd71ed424596ccfd31112e0fa9e9684d449a720da0c626e0916797d4ca7
addr: 0xb313e2A167a44BF1b110fD2f882d5Bd86A9aa022
keyType: ecdsa
derivationPath: m/44'/60'/0'/0/0
```

## fund sender address
Depending on the ethereum network you may have different options to fund the sender's address
with some ether coins.

For instance, [goerli faucet](https://goerlifaucet.com/) can be used to fund the sender's
account assuming transactions will also be performed on this network.

## check balance
Assuming you have the env. var `ETHIO_RPC_ENDPOINT`, a call can be made to
the network to get current balance for an account:
```bash
ethio balance 0xb313e2A167a44BF1b110fD2f882d5Bd86A9aa022
```
```text
0
```

## send ether
Send ether using sender's private key
```bash
ethio send \
  --key=5040f993d59dd4dbba2bbcb0286371bdb02dc0bacc73ce56b989d07e9ef8bb24 \
  --addr=0xb313e2A167a44BF1b110fD2f882d5Bd86A9aa022 \
  --amount=0.1 \
  --unit=eth \
  --gas=30
```

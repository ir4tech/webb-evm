package precompile

import (
	"github.com/ethereum/go-ethereum/common"
)

var (
	_                 StatefulPrecompileConfig = &TxAssetConfig{}
	TxAssetPrecompile                          = createAllowListPrecompile(TxAllowListAddress)
)

type TxAssetConfig struct {
	AssetConfig
}

func (c *TxAssetConfig) Address() common.Address {
	return TxAssetAddress
}

func (c *TxAssetConfig) Configure(_ ChainConfig, state StateDB, _ BlockContext) {
	c.AssetConfig.Configure(state, TxAssetAddress)
}

func (c *TxAssetConfig) Contract() StatefulPrecompiledContract {
	return TxAssetPrecompile
}

func GetTxLocation(stateDB StateDB) AssetLocation {
	return getLocation(stateDB, TxAssetAddress)
}

func SetTxLocation(stateDB StateDB, location AssetLocation) {
	setLocation(stateDB, TxAssetAddress, location)
}

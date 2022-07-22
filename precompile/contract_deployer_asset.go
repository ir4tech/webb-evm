package precompile

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ir4tech/webb-evm/core/types"
)

var (
	_                               StatefulPrecompileConfig = &ContractDeployerAssetConfig{}
	ContractDeployerAssetPrecompile                          = createAssetPrecompile(ContractDeployerAssetAddress)
)

type ContractDeployerAssetConfig struct {
	AssetConfig
}

func (c *ContractDeployerAssetConfig) Address() common.Address {
	return ContractDeployerAssetAddress
}

func (c *ContractDeployerAssetConfig) Configure(_ ChainConfig, state StateDB, _ BlockContext) {
	c.AssetConfig.Configure(state, ContractDeployerAssetAddress)
}

func (c *ContractDeployerAssetConfig) Contract() StatefulPrecompiledContract {
	return ContractDeployerAssetPrecompile
}

func GetContractDeployerAssetStatus(assetDB AssetDB, assetId common.Hash) (types.Asset, error) {
	return getAsset(assetDB, assetId)
}

func RegisterContractDeployerAssetStatus(assetDB AssetDB, owner common.Address, location string) {
	registerAsset(assetDB, owner, location)
}

package precompile

import "github.com/ethereum/go-ethereum/common"

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

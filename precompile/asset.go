package precompile

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ir4tech/webb-evm/core/types"
	"github.com/ir4tech/webb-evm/vmerrs"
	"math/big"
)

type AssetLocation string

var (
	DefaultAsset = types.Asset{}

	registerAssetSignature     = CalculateFunctionSelector("registerAsset(address, string)")
	getAssetSignature          = CalculateFunctionSelector("getAsset(string)")
	getAssetByAddressSignature = CalculateFunctionSelector("getAssetByAddress(address)")
	updateLocationSignature    = CalculateFunctionSelector("updateLocation(string, string)")
	updateNameSignature        = CalculateFunctionSelector("updateName(string, string)")
)

type AssetConfig struct {
	BlockTimestamp *big.Int `json:"blockTimestamp"`
	Location       string   `json:"location"`
}

func (asset *AssetConfig) Timestamp() *big.Int { return asset.BlockTimestamp }

func (asset *AssetConfig) Configure(_ StateDB, _ common.Address) {

}

func registerAsset(assetDB AssetDB, owner common.Address, name string) {
	assetDB.RegisterAsset(owner, name)
}

func getAsset(assetDB AssetDB, assetId common.Hash) (types.Asset, error) {
	asset, err := assetDB.GetAsset(assetId)
	if err != nil {
		return DefaultAsset, err
	}
	return asset, nil
}

func getAssetByAddress(assetDB AssetDB, address common.Address) []types.Asset {
	return assetDB.GetAssetByOwner(address)
}

func updateLocation(assetDB AssetDB, assetId common.Hash, location string) {
	assetDB.UpdateLocation(assetId, location)
}

func updateName(assetDB AssetDB, assetId common.Hash, name string) {
	assetDB.UpdateName(assetId, name)
}

func createAssetPrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	// Construct the contract with no fallback function.
	allowListFuncs := createAssetFunctions(precompileAddr)
	contract := newStatefulPrecompileWithFunctionSelectors(nil, allowListFuncs)
	return contract
}

func createAssetFunctions(precompileAddr common.Address) []*statefulPrecompileFunction {
	registerAsset := newStatefulPrecompileFunction(registerAssetSignature, createAssetSetter(precompileAddr, DefaultAsset))
	getAsset := newStatefulPrecompileFunction(getAssetSignature, createAssetSetter(precompileAddr, DefaultAsset))
	getAssetByAddress := newStatefulPrecompileFunction(getAssetByAddressSignature, createAssetSetter(precompileAddr, DefaultAsset))
	updateLocation := newStatefulPrecompileFunction(updateLocationSignature, createAssetSetter(precompileAddr, DefaultAsset))
	updateName := newStatefulPrecompileFunction(updateNameSignature, createAssetSetter(precompileAddr, DefaultAsset))
	return []*statefulPrecompileFunction{registerAsset, getAsset, getAssetByAddress, updateLocation, updateName}
}

func createAssetSetter(precompileAddr common.Address, asset types.Asset) RunStatefulPrecompileFunc {
	return func(evm PrecompileAccessibleState, callerAddr, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
		if remainingGas, err = deductGas(suppliedGas, ModifyAllowListGasCost); err != nil {
			return nil, 0, err
		}

		if len(input) != allowListInputLen {
			return nil, remainingGas, fmt.Errorf("invalid input length for modifying allow list: %d", len(input))
		}

		if readOnly {
			return nil, remainingGas, vmerrs.ErrWriteProtection
		}

		assetDB := evm.GetAssetDB()
		stateDB := evm.GetStateDB()

		// Verify that the caller is in the allow list and therefore has the right to modify it
		callerStatus := getAllowListStatus(stateDB, precompileAddr, callerAddr)
		if !callerStatus.IsAdmin() {
			return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotModifyAllowList, callerAddr)
		}

		registerAsset(assetDB(), precompileAddr, "")
		// Return an empty output and the remaining gas
		return []byte{}, remainingGas, nil
	}
}

package precompile

import (
	"fmt"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"
	"math/big"
)

type AssetLocation string

var (
	DefaultAsset = commontype.Asset{}

	registerAssetSignature     = CalculateFunctionSelector("registerAsset()")
	getAssetSignature          = CalculateFunctionSelector("getAsset(string)")
	getAllAssetsSignature      = CalculateFunctionSelector("getAll()")
	getAssetByAddressSignature = CalculateFunctionSelector("getAssetByAddress(address)")
	updateLocationSignature    = CalculateFunctionSelector("updateLocation(string,string)")
	updateNameSignature        = CalculateFunctionSelector("updateName(string,string)")

	getAssetByIdInputLen      = common.HashLength
	getAssetByAddressInputLen = common.AddressLength
)

type AssetConfig struct {
	BlockTimestamp *big.Int `json:"blockTimestamp"`
}

func (asset *AssetConfig) Timestamp() *big.Int { return asset.BlockTimestamp }

func (asset *AssetConfig) Configure(_ StateDB, _ common.Address) {

}

func getAll(assetDB AssetDB) []commontype.Asset {
	return assetDB.GetAll()
}

func registerAsset(assetDB AssetDB, owner common.Address, name string) {
	assetId := common.BytesToHash([]byte(uuid.New().String()))
	assetDB.RegisterAsset(assetId, owner, name)
}

func getAsset(assetDB AssetDB, assetId common.Hash) (commontype.Asset, error) {
	asset, err := assetDB.GetAsset(assetId)
	if err != nil {
		return DefaultAsset, err
	}
	return asset, nil
}

func getAssetByAddress(assetDB AssetDB, address common.Address) []commontype.Asset {
	return assetDB.GetAssetByOwner(address)
}

func updateLocation(assetDB AssetDB, assetId common.Hash, location string) {
	assetDB.UpdateLocation(assetId, location)
}

func updateName(assetDB AssetDB, stateDB StateDB, precompileAddr common.Address, assetId common.Hash, name string) {
	assetDB.UpdateName(assetId, name)
}

func createAssetPrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	// Construct the contract with no fallback function.
	allowListFuncs := createAssetFunctions(precompileAddr)
	contract := newStatefulPrecompileWithFunctionSelectors(nil, allowListFuncs)
	return contract
}

func createAssetFunctions(precompileAddr common.Address) []*statefulPrecompileFunction {
	registerAsset := newStatefulPrecompileFunction(registerAssetSignature, createRegisterAsset(precompileAddr))
	getAsset := newStatefulPrecompileFunction(getAssetSignature, createGetAsset(precompileAddr))
	getAllAssets := newStatefulPrecompileFunction(getAllAssetsSignature, createGetAllAssets(precompileAddr))
	getAssetByAddress := newStatefulPrecompileFunction(getAssetByAddressSignature, createGetAssetByAddress(precompileAddr))
	updateLocation := newStatefulPrecompileFunction(updateLocationSignature, createGetAsset(precompileAddr)) // TODO: needs the correct create function
	updateName := newStatefulPrecompileFunction(updateNameSignature, createGetAsset(precompileAddr))         // TODO: needs the correct create function
	return []*statefulPrecompileFunction{registerAsset, getAsset, getAllAssets, getAssetByAddress, updateLocation, updateName}
}

func createRegisterAsset(precompileAddr common.Address) RunStatefulPrecompileFunc {
	return func(evm PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
		if remainingGas, err = deductGas(suppliedGas, RegisterAssetGasCost); err != nil {
			return nil, 0, err
		}

		assetDB := evm.GetAssetDB()

		registerAsset(assetDB, caller, "-")
		// Return an empty output and the remaining gas
		return []byte{}, remainingGas, nil
	}
}

func createGetAsset(_ common.Address) RunStatefulPrecompileFunc {
	return func(evm PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
		if remainingGas, err = deductGas(suppliedGas, GetAssetGasCost); err != nil {
			return nil, 0, err
		}

		if len(input) != getAssetByIdInputLen {
			return nil, remainingGas, fmt.Errorf("invalid input length for getting an asset: %d", len(input))
		}

		assetDB := evm.GetAssetDB()
		assetId := common.BytesToHash(input)

		asset, err := getAsset(assetDB, assetId)
		if err != nil {
			log.Info("precompile: not found asset: ", assetId)
			return nil, remainingGas, fmt.Errorf("error: %d", err)
		}

		assetBytes := []byte(fmt.Sprintf("%v", asset))

		log.Info("precompile: found asset = ", asset)
		log.Info("precompile: asset to bytes = ", assetBytes)

		return assetBytes, remainingGas, nil
	}
}

func createGetAllAssets(_ common.Address) RunStatefulPrecompileFunc {
	return func(evm PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
		if remainingGas, err = deductGas(suppliedGas, GetAssetGasCost); err != nil {
			return nil, 0, err
		}

		assetDB := evm.GetAssetDB()
		asset := getAll(assetDB)

		return []byte(fmt.Sprintf("%v", asset)), remainingGas, nil
	}
}

func createGetAssetByAddress(_ common.Address) RunStatefulPrecompileFunc {
	return func(evm PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
		if remainingGas, err = deductGas(suppliedGas, GetAssetGasCost); err != nil {
			return nil, 0, err
		}

		if len(input) != getAssetByAddressInputLen {
			return nil, remainingGas, fmt.Errorf("invalid input length for getting an asset by address: %d", len(input))
		}

		assetDB := evm.GetAssetDB()
		owner := common.BytesToAddress(input)
		assets := getAssetByAddress(assetDB, owner)
		assetBytes := []byte(fmt.Sprintf("%v", assets))

		log.Info("precompile: found assets = ", assets)
		log.Info("precompile: asset to bytes = ", assetBytes)

		return assetBytes, remainingGas, nil
	}
}

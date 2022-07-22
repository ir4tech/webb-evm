package asset

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/ir4tech/webb-evm/core/types"
)

var (
	defaultLocation  = "-"
	emptyAsset       = types.Asset{}
	errAssetNotFound = errors.New("asset not found")
)

type AssetDB struct {
	assets map[common.Hash]types.Asset
	owners map[common.Address][]common.Hash
}

func New() *AssetDB {
	return &AssetDB{
		assets: make(map[common.Hash]types.Asset),
		owners: make(map[common.Address][]common.Hash),
	}
}

func (store AssetDB) RegisterAsset(owner common.Address, name string) {
	assetId := common.BytesToHash([]byte(uuid.New().String()))
	store.assets[assetId] = types.Asset{
		Id:       assetId,
		Name:     name,
		Owner:    owner,
		Location: defaultLocation,
	}
	store.owners[owner] = append(store.owners[owner], assetId)
}

func (store AssetDB) GetAsset(assetId common.Hash) (types.Asset, error) {
	asset, ok := store.assets[assetId]
	if ok {
		return asset, nil
	}
	return emptyAsset, errAssetNotFound
}

func (store AssetDB) GetAssetByOwner(owner common.Address) []types.Asset {
	assetIds := store.owners[owner]
	var assets []types.Asset
	for _, assetId := range assetIds {
		assets = append(assets, store.assets[assetId])
	}
	return assets
}

func (store AssetDB) UpdateLocation(assetId common.Hash, location string) error {
	asset, ok := store.assets[assetId]
	if ok {
		asset.Location = location
		store.assets[assetId] = asset
		return nil
	}
	return errAssetNotFound
}

func (store AssetDB) UpdateName(assetId common.Hash, name string) error {
	asset, ok := store.assets[assetId]
	if ok {
		asset.Name = name
		store.assets[assetId] = asset
		return nil
	}
	return errAssetNotFound
}

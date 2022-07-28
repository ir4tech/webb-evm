package assetdb

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ethereum/go-ethereum/log"
)

var (
	defaultLocation  = "-"
	emptyAsset       = commontype.Asset{}
	errAssetNotFound = errors.New("asset not found")
)

type AssetDB struct {
	assets map[common.Hash]commontype.Asset
	owners map[common.Address][]common.Hash
}

func New() *AssetDB {
	return &AssetDB{
		assets: make(map[common.Hash]commontype.Asset),
		owners: make(map[common.Address][]common.Hash),
	}
}

func (store AssetDB) GetAll() []commontype.Asset {
	log.Info("assetdb: fetching all assets in the db")
	assets := []commontype.Asset{}
	for _, v := range store.assets {
		assets = append(assets, v)
	}
	log.Info("assetdb: assets=", assets)
	return assets
}

func (store AssetDB) RegisterAsset(assetId common.Hash, owner common.Address, name string) {
	log.Info("assetdb: registering new asset with assetId=", assetId, "owner=", owner, "name=", name)
	store.assets[assetId] = commontype.Asset{
		Id:       assetId,
		Name:     name,
		Owner:    owner,
		Location: defaultLocation,
	}
	store.owners[owner] = append(store.owners[owner], assetId)
}

func (store AssetDB) GetAsset(assetId common.Hash) (commontype.Asset, error) {
	log.Info("assetdb: looking for asset with id=", assetId)
	asset, ok := store.assets[assetId]
	if ok {
		log.Info("assetdb: found asset=", asset)
		return asset, nil
	}
	log.Warn("assetdb: did not find asset with id=", assetId)
	return emptyAsset, errAssetNotFound
}

func (store AssetDB) GetAssetByOwner(owner common.Address) []commontype.Asset {
	log.Info("assetdb: looking for assets owned by=", owner)
	assetIds := store.owners[owner]
	var assets []commontype.Asset
	for _, assetId := range assetIds {
		assets = append(assets, store.assets[assetId])
	}
	log.Info("assetdb: found=", assets)
	return assets
}

func (store AssetDB) UpdateLocation(assetId common.Hash, location string) error {
	log.Info("assetdb: will update location of asset=", assetId, "with location=", location)
	asset, ok := store.assets[assetId]
	if ok {
		asset.Location = location
		store.assets[assetId] = asset
		log.Info("assetdb: location was updated")
		return nil
	}
	log.Warn("assetdb: did not find asset with id=", assetId)
	return errAssetNotFound
}

func (store AssetDB) UpdateName(assetId common.Hash, name string) error {
	log.Info("assetdb: will update name of asset=", assetId, "with name=", name)
	asset, ok := store.assets[assetId]
	if ok {
		asset.Name = name
		store.assets[assetId] = asset
		return nil
	}
	log.Warn("assetdb: did not find asset with id=", assetId)
	return errAssetNotFound
}

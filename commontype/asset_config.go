package commontype

import "github.com/ethereum/go-ethereum/common"

type Asset struct {
	Id       common.Hash    `json:"id"`
	Name     string         `json:"name"`
	Owner    common.Address `json:"owner"`
	Location string         `json:"location"`
}

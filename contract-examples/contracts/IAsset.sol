//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

struct Asset {
    string id;
    string name;
    address owner;
    string location;
}

interface IAsset {
    function getAll() external view returns (Asset[] memory);
    function registerAsset() external;
    function getAsset(string memory assetId) external view returns (Asset memory);
    function getAssetByAddress(address owner) external view returns (Asset memory);
    function updateLocation(string memory assetId, string memory location) external;
    function updateName(string memory assetId, string memory name) external;
}
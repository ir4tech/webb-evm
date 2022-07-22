//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IAsset {
    function getLocation() external view returns (string memory);
    function setLocation(string memory location) external;
}
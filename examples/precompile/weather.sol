pragma solidity ^0.8.19;

contract weather {
    function getCurrentGameWeather() external view returns(uint8) {
        return gameWeather();
    }
}
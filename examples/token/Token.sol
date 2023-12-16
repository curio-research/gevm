pragma solidity 0.8.19;

import "./ERC20Flattened.sol";

contract Token is ERC20 {
    constructor() ERC20("Test Coin", "TEST") {}

    function mint(uint256 amount) public {
        _mint(msg.sender, amount);
    }
}

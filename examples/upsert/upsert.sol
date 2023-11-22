pragma solidity ^0.6.0;

contract upsert {
      // Function to send Ether. The `payable` keyword is used to allow the function to receive Ether.
    function sendETH(address payable _to) external payable {
        // Check that the function received some Ether
        require(msg.value > 0, "No Ether sent");

        // Transfer the Ether to the specified address
        _to.transfer(msg.value);
    }
}

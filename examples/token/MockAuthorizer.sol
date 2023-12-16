pragma solidity ^0.7.0;

interface IAuthorizer {
    /**
     * @dev Returns true if `account` can perform the action described by `actionId` in the contract `where`.
     */
    function canPerform(
        bytes32 actionId,
        address account,
        address where
    ) external view returns (bool);

}


contract MockAuthorizer is IAuthorizer {
    function canPerform(bytes32, address, address) external view returns (bool) {
        return true;
    }
}

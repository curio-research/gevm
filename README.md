# gevm
extremely minimal POC of an in-memory blockchain, (headers)
with persistent state via PebbleDB.



### multiplication example

create new EVMContext.
create new PebbleDB.
create new StateDB. (statedb)
set balances in stateDB.

create new EVM.
create new Contract within EVM.
call contract using evm.Call.
check state after Call. 
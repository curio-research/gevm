# gevm
extremely minimal POC of an in-memory blockchain, (headers)
with persistent state via PebbleDB.




### multiplication example

`./run-sum`
or
`cd examples/sum && go run main.go`

create new EVMContext.
create new PebbleDB.
create new StateDB. (statedb)
set balances in stateDB.

create new EVM.
create new Contract within EVM.
call contract using evm.Call.
check state after Call. 

### token example
 create token contract
 mint a token
 transfer the toke

`cd examples/token && go run main.go`

### precompile example
 create weather contract
 call precompile function

`cd examples/precompile && go run main.go`

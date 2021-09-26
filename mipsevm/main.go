package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

type StateDB struct {
	Bytecode []byte
}

func (s *StateDB) AddAddressToAccessList(addr common.Address)                {}
func (s *StateDB) AddBalance(addr common.Address, amount *big.Int)           {}
func (s *StateDB) AddLog(log *types.Log)                                     {}
func (s *StateDB) AddPreimage(hash common.Hash, preimage []byte)             {}
func (s *StateDB) AddRefund(gas uint64)                                      {}
func (s *StateDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {}
func (s *StateDB) AddressInAccessList(addr common.Address) bool              { return true }
func (s *StateDB) CreateAccount(addr common.Address)                         {}
func (s *StateDB) Empty(addr common.Address) bool                            { return false }
func (s *StateDB) Exist(addr common.Address) bool                            { return true }
func (b *StateDB) ForEachStorage(addr common.Address, cb func(key, value common.Hash) bool) error {
	return nil
}
func (s *StateDB) GetBalance(addr common.Address) *big.Int { return common.Big0 }
func (s *StateDB) GetCode(addr common.Address) []byte {
	fmt.Println("GetCode", addr)
	return s.Bytecode
}
func (s *StateDB) GetCodeHash(addr common.Address) common.Hash { return common.Hash{} }
func (s *StateDB) GetCodeSize(addr common.Address) int         { return 0 }
func (s *StateDB) GetCommittedState(addr common.Address, hash common.Hash) common.Hash {
	return common.Hash{}
}
func (s *StateDB) GetNonce(addr common.Address) uint64 { return 0 }
func (s *StateDB) GetRefund() uint64                   { return 0 }
func (s *StateDB) GetState(addr common.Address, hash common.Hash) common.Hash {
	fmt.Println("GetState", addr, hash)
	return common.Hash{}
}
func (s *StateDB) HasSuicided(addr common.Address) bool { return false }
func (s *StateDB) PrepareAccessList(sender common.Address, dst *common.Address, precompiles []common.Address, list types.AccessList) {
}
func (s *StateDB) RevertToSnapshot(revid int)                           {}
func (s *StateDB) SetCode(addr common.Address, code []byte)             {}
func (s *StateDB) SetNonce(addr common.Address, nonce uint64)           {}
func (s *StateDB) SetState(addr common.Address, key, value common.Hash) {}
func (s *StateDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressPresent bool, slotPresent bool) {
	return true, true
}
func (s *StateDB) Snapshot() int                                   { return 0 }
func (s *StateDB) SubBalance(addr common.Address, amount *big.Int) {}
func (s *StateDB) SubRefund(gas uint64)                            {}
func (s *StateDB) Suicide(addr common.Address) bool                { return true }

// **** stub Tracer ****
type Tracer struct{}

// CaptureStart implements the Tracer interface to initialize the tracing operation.
func (jst *Tracer) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
}

// CaptureState implements the Tracer interface to trace a single step of VM execution.
func (jst *Tracer) CaptureState(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	fmt.Println(pc, op, gas)
}

// CaptureFault implements the Tracer interface to trace an execution fault
func (jst *Tracer) CaptureFault(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}

// CaptureEnd is called after the call finishes to finalize the tracing.
func (jst *Tracer) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
}

type jsoncontract struct {
	Bytecode         string `json:"bytecode"`
	DeployedBytecode string `json:"deployedBytecode"`
}

func main() {
	fmt.Println("hello")

	/*var parent types.Header
	database := state.NewDatabase(parent)
	statedb, _ := state.New(parent.Root, database, nil)*/

	var jj jsoncontract
	mipsjson, _ := ioutil.ReadFile("../artifacts/contracts/MIPS.sol/MIPS.json")
	json.NewDecoder(bytes.NewReader(mipsjson)).Decode(&jj)
	bytecode := common.Hex2Bytes(jj.DeployedBytecode[2:])
	//fmt.Println(bytecode, jj.Bytecode)
	bytecodehash := crypto.Keccak256Hash(bytecode)
	statedb := &StateDB{Bytecode: bytecode}

	bc := core.NewBlockChain()
	var header types.Header
	header.Number = big.NewInt(13284469)
	header.Difficulty = common.Big0
	author := common.Address{}
	blockContext := core.NewEVMBlockContext(&header, bc, &author)
	txContext := vm.TxContext{}
	config := vm.Config{}
	config.Debug = true
	tracer := Tracer{}
	config.Tracer = &tracer
	evm := vm.NewEVM(blockContext, txContext, statedb, params.MainnetChainConfig, config)
	fmt.Println(evm)

	from := common.Address{}
	to := common.HexToAddress("0x1337")
	/*ret, gas, err := evm.Call(vm.AccountRef(from), to, []byte{}, 20000000, common.Big0)
	fmt.Println(ret, gas, err)*/

	interpreter := vm.NewEVMInterpreter(evm, config)
	input := []byte{0x69, 0x37, 0x33, 0x72} // Step(bytes32)
	input = append(input, common.Hash{}.Bytes()...)
	contract := vm.NewContract(vm.AccountRef(from), vm.AccountRef(to), common.Big0, 20000000)
	//fmt.Println(bytecodehash, bytecode)
	contract.SetCallCode(&to, bytecodehash, bytecode)
	ret, err := interpreter.Run(contract, input, false)
	fmt.Println(ret, err, contract.Gas)
	if err != nil {
		log.Fatal(err)
	}
}
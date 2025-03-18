package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/evm/x/vm/core/vm"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// EmitEventArgs are the arguments required to emit an authorization event.
//
// The event type can be:
//   - ApprovalEvent
//   - GenericApprovalEvent
//   - AllowanceChangeEvent
//   - ...
type EmitEventArgs struct {
	Ctx            sdk.Context
	StateDB        vm.StateDB
	ContractAddr   common.Address
	ContractEvents map[string]abi.Event
	EventData      interface{}
}

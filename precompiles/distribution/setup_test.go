package distribution_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/evm/precompiles/distribution"
	testconstants "github.com/cosmos/evm/testutil/constants"
	"github.com/cosmos/evm/testutil/integration/os/factory"
	"github.com/cosmos/evm/testutil/integration/os/grpc"
	testkeyring "github.com/cosmos/evm/testutil/integration/os/keyring"
	"github.com/cosmos/evm/testutil/integration/os/network"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/stretchr/testify/suite"
)

type PrecompileTestSuite struct {
	suite.Suite

	network     *network.UnitTestNetwork
	factory     factory.TxFactory
	grpcHandler grpc.Handler
	keyring     testkeyring.Keyring

	precompile           *distribution.Precompile
	bondDenom            string
	baseDenom            string
	validatorsKeys       []testkeyring.Key
	withValidatorSlashes bool
}

func TestPrecompileUnitTestSuite(t *testing.T) {
	suite.Run(t, new(PrecompileTestSuite))
}

func (s *PrecompileTestSuite) SetupTest() {
	keyring := testkeyring.New(2)
	s.validatorsKeys = generateKeys(3)
	customGen := network.CustomGenesisState{}

	// set some slashing events for integration test
	distrGen := distrtypes.DefaultGenesisState()
	if s.withValidatorSlashes {
		distrGen.ValidatorSlashEvents = []distrtypes.ValidatorSlashEventRecord{
			{
				ValidatorAddress:    sdk.ValAddress(s.validatorsKeys[0].Addr.Bytes()).String(),
				Height:              0,
				Period:              1,
				ValidatorSlashEvent: distrtypes.NewValidatorSlashEvent(1, math.LegacyNewDecWithPrec(5, 2)),
			},
			{
				ValidatorAddress:    sdk.ValAddress(s.validatorsKeys[0].Addr.Bytes()).String(),
				Height:              1,
				Period:              1,
				ValidatorSlashEvent: distrtypes.NewValidatorSlashEvent(1, math.LegacyNewDecWithPrec(5, 2)),
			},
		}
	}
	customGen[distrtypes.ModuleName] = distrGen

	// set non-zero inflation for rewards to accrue (use defaults from SDK for values)
	mintGen := minttypes.DefaultGenesisState()
	mintGen.Params.MintDenom = testconstants.ExampleAttoDenom
	customGen[minttypes.ModuleName] = mintGen

	operatorsAddr := make([]sdk.AccAddress, 3)
	for i, k := range s.validatorsKeys {
		operatorsAddr[i] = k.AccAddr
	}

	nw := network.NewUnitTestNetwork(
		network.WithPreFundedAccounts(keyring.GetAllAccAddrs()...),
		network.WithCustomGenesis(customGen),
		network.WithValidatorOperators(operatorsAddr),
	)
	grpcHandler := grpc.NewIntegrationHandler(nw)
	txFactory := factory.New(nw, grpcHandler)

	ctx := nw.GetContext()
	sk := nw.App.StakingKeeper
	bondDenom, err := sk.BondDenom(ctx)
	if err != nil {
		panic(err)
	}

	s.bondDenom = bondDenom
	// TODO: check if this is correct?
	s.baseDenom = evmtypes.GetEVMCoinDenom()

	s.factory = txFactory
	s.grpcHandler = grpcHandler
	s.keyring = keyring
	s.network = nw
	s.precompile, err = distribution.NewPrecompile(
		s.network.App.DistrKeeper,
		*s.network.App.StakingKeeper,
		s.network.App.AuthzKeeper,
		s.network.App.EVMKeeper,
	)
	if err != nil {
		panic(err)
	}
}

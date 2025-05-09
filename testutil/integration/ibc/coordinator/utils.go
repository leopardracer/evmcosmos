package coordinator

import (
	"strconv"
	"testing"

	"github.com/cosmos/evm/testutil/integration/common/network"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// getIBCChains returns a map of TestChain's for the given network interface.
func getIBCChains(t *testing.T, coord *ibctesting.Coordinator, chains []network.Network) map[string]*ibctesting.TestChain {
	t.Helper()
	ibcChains := make(map[string]*ibctesting.TestChain)
	for _, chain := range chains {
		ibcChains[chain.GetChainID()] = chain.GetIBCChain(t, coord)
	}
	return ibcChains
}

// generateDummyChains returns a map of dummy chains to complement IBC connections for integration tests.
func generateDummyChains(t *testing.T, coord *ibctesting.Coordinator, numberOfChains int) (map[string]*ibctesting.TestChain, []string) {
	t.Helper()
	ibcChains := make(map[string]*ibctesting.TestChain)
	ids := make([]string, numberOfChains)
	// dummy chains use the ibc testing chain setup
	// that uses the default sdk address prefix ('cosmos')
	// Update the prefix configs to use that prefix
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	// Also need to disable address cache to avoid using modules
	// accounts with 'evmos' addresses (because Cosmos EVM chain setup is first)
	sdk.SetAddrCacheEnabled(false)
	for i := 1; i <= numberOfChains; i++ {
		chainID := "dummychain-" + strconv.Itoa(i)
		ids[i-1] = chainID
		ibcChains[chainID] = ibctesting.NewTestChain(t, coord, chainID)
	}

	return ibcChains, ids
}

// mergeMaps merges two maps of TestChain's.
func mergeMaps(m1, m2 map[string]*ibctesting.TestChain) map[string]*ibctesting.TestChain {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

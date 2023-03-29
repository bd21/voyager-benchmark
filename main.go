package main

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

var ENDPOINTS = []string{
	//"{LCD_HOST}/cosmos/auth/v1beta1/account_info/{address}",
	//"{LCD_HOST}/cosmos/auth/v1beta1/accounts",
	"{LCD_HOST}/cosmos/auth/v1beta1/accounts/{address}",
	//"{LCD_HOST}/cosmos/auth/v1beta1/address_by_id/{id}",
	//"{LCD_HOST}/cosmos/auth/v1beta1/bech32",
	//"{LCD_HOST}/cosmos/auth/v1beta1/bech32/{address_bytes}",
	//"{LCD_HOST}/cosmos/auth/v1beta1/bech32/{address_string}",
	//"{LCD_HOST}/cosmos/auth/v1beta1/module_accounts",
	"{LCD_HOST}/cosmos/auth/v1beta1/module_accounts/{name}", // TODO 404
	"{LCD_HOST}/cosmos/auth/v1beta1/params",
	"{LCD_HOST}/cosmos/bank/v1beta1/balances/{address}",
	"{LCD_HOST}/cosmos/bank/v1beta1/balances/{address}/by_denom", // TODO bad request
	//"{LCD_HOST}/cosmos/bank/v1beta1/denom_owners/{denom}",
	"{LCD_HOST}/cosmos/bank/v1beta1/denoms_metadata",
	"{LCD_HOST}/cosmos/bank/v1beta1/denoms_metadata/{denom}",
	"{LCD_HOST}/cosmos/bank/v1beta1/params",
	//"{LCD_HOST}/cosmos/bank/v1beta1/send_enabled",
	"{LCD_HOST}/cosmos/bank/v1beta1/spendable_balances/{address}",
	//"{LCD_HOST}/cosmos/bank/v1beta1/spendable_balances/{address}/by_denom",
	"{LCD_HOST}/cosmos/bank/v1beta1/supply",
	//"{LCD_HOST}/cosmos/bank/v1beta1/supply/by_denom", // TODO 500
	//"{LCD_HOST}/cosmos/base/tendermint/v1beta1/abci_query",
	"{LCD_HOST}/cosmos/base/tendermint/v1beta1/blocks/latest",
	"{LCD_HOST}/cosmos/base/tendermint/v1beta1/blocks/{height}",
	"{LCD_HOST}/cosmos/base/tendermint/v1beta1/node_info",
	"{LCD_HOST}/cosmos/base/tendermint/v1beta1/syncing",
	"{LCD_HOST}/cosmos/base/tendermint/v1beta1/validatorsets/latest",
	"{LCD_HOST}/cosmos/base/tendermint/v1beta1/validatorsets/{height}",
	"{LCD_HOST}/cosmos/distribution/v1beta1/community_pool",
	//"{LCD_HOST}/cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards", // TODO 400
	"{LCD_HOST}/cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards/{validator_address}",
	"{LCD_HOST}/cosmos/distribution/v1beta1/delegators/{delegator_address}/validators",
	"{LCD_HOST}/cosmos/distribution/v1beta1/delegators/{delegator_address}/withdraw_address",
	"{LCD_HOST}/cosmos/distribution/v1beta1/params",
	//"{LCD_HOST}/cosmos/distribution/v1beta1/validators/{validator_address}",
	//"{LCD_HOST}/cosmos/distribution/v1beta1/validators/{validator_address}/commission", // TODO 400
	//"{LCD_HOST}/cosmos/distribution/v1beta1/validators/{validator_address}/outstanding_rewards", // TODO 400
	// "{LCD_HOST}/cosmos/distribution/v1beta1/validators/{validator_address}/slashes", // TODO 400
	"{LCD_HOST}/cosmos/evidence/v1beta1/evidence",
	// "{LCD_HOST}/cosmos/evidence/v1beta1/evidence/{hash}", // TODO 404
	// "{LCD_HOST}/cosmos/gov/v1beta1/params/{params_type}", // TODO 400
	"{LCD_HOST}/cosmos/gov/v1beta1/proposals",
	"{LCD_HOST}/cosmos/gov/v1beta1/proposals/{proposal_id}",
	"{LCD_HOST}/cosmos/gov/v1beta1/proposals/{proposal_id}/deposits",
	"{LCD_HOST}/cosmos/gov/v1beta1/proposals/{proposal_id}/deposits/{depositor}",
	"{LCD_HOST}/cosmos/gov/v1beta1/proposals/{proposal_id}/tally",
	"{LCD_HOST}/cosmos/gov/v1beta1/proposals/{proposal_id}/votes",
	"{LCD_HOST}/cosmos/gov/v1beta1/proposals/{proposal_id}/votes/{voter}",
	"{LCD_HOST}/cosmos/gov/v1/params/{params_type}",
	"{LCD_HOST}/cosmos/gov/v1/proposals",
	"{LCD_HOST}/cosmos/gov/v1/proposals/{proposal_id}",
	"{LCD_HOST}/cosmos/gov/v1/proposals/{proposal_id}/deposits",
	"{LCD_HOST}/cosmos/gov/v1/proposals/{proposal_id}/deposits/{depositor}",
	"{LCD_HOST}/cosmos/gov/v1/proposals/{proposal_id}/tally",
	"{LCD_HOST}/cosmos/gov/v1/proposals/{proposal_id}/votes",
	"{LCD_HOST}/cosmos/gov/v1/proposals/{proposal_id}/votes/{voter}",
	"{LCD_HOST}/cosmos/mint/v1beta1/annual_provisions",
	"{LCD_HOST}/cosmos/mint/v1beta1/inflation",
	"{LCD_HOST}/cosmos/mint/v1beta1/params",
	"{LCD_HOST}/cosmos/params/v1beta1/params",
	"{LCD_HOST}/cosmos/params/v1beta1/subspaces",
	"{LCD_HOST}/cosmos/slashing/v1beta1/params",
	"{LCD_HOST}/cosmos/slashing/v1beta1/signing_infos",
	"{LCD_HOST}/cosmos/slashing/v1beta1/signing_infos/{cons_address}",
	"{LCD_HOST}/cosmos/staking/v1beta1/delegations/{delegator_addr}",
	"{LCD_HOST}/cosmos/staking/v1beta1/delegators/{delegator_addr}/redelegations",
	"{LCD_HOST}/cosmos/staking/v1beta1/delegators/{delegator_addr}/unbonding_delegations",
	"{LCD_HOST}/cosmos/staking/v1beta1/delegators/{delegator_addr}/validators",
	"{LCD_HOST}/cosmos/staking/v1beta1/delegators/{delegator_addr}/validators/{validator_addr}",
	"{LCD_HOST}/cosmos/staking/v1beta1/historical_info/{height}",
	"{LCD_HOST}/cosmos/staking/v1beta1/params",
	"{LCD_HOST}/cosmos/staking/v1beta1/pool",
	"{LCD_HOST}/cosmos/staking/v1beta1/validators",
	"{LCD_HOST}/cosmos/staking/v1beta1/validators/{validator_addr}",
	"{LCD_HOST}/cosmos/staking/v1beta1/validators/{validator_addr}/delegations",
	"{LCD_HOST}/cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}",
	"{LCD_HOST}/cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}/unbonding_delegation",
	"{LCD_HOST}/cosmos/staking/v1beta1/validators/{validator_addr}/unbonding_delegations",
	"{LCD_HOST}/cosmos/tx/v1beta1/decode",
	"{LCD_HOST}/cosmos/tx/v1beta1/decode/amino",
	"{LCD_HOST}/cosmos/tx/v1beta1/encode",
	"{LCD_HOST}/cosmos/tx/v1beta1/encode/amino",
	"{LCD_HOST}/cosmos/tx/v1beta1/simulate",
	"{LCD_HOST}/cosmos/tx/v1beta1/txs",
	"{LCD_HOST}/cosmos/tx/v1beta1/txs/block/{height}",
	"{LCD_HOST}/cosmos/tx/v1beta1/txs/{hash}",
	"{LCD_HOST}/cosmos/upgrade/v1beta1/applied_plan/{name}",
	"{LCD_HOST}/cosmos/upgrade/v1beta1/authority",
	"{LCD_HOST}/cosmos/upgrade/v1beta1/current_plan",
	"{LCD_HOST}/cosmos/upgrade/v1beta1/module_versions",
	"{LCD_HOST}/cosmos/upgrade/v1beta1/upgraded_consensus_state/{last_height}",
	"{LCD_HOST}/cosmos/authz/v1beta1/grants",
	"{LCD_HOST}/cosmos/authz/v1beta1/grants/grantee/{grantee}",
	"{LCD_HOST}/cosmos/authz/v1beta1/grants/granter/{granter}",
	"{LCD_HOST}/cosmos/feegrant/v1beta1/allowance/{granter}/{grantee}",
	"{LCD_HOST}/cosmos/feegrant/v1beta1/allowances/{grantee}",
	"{LCD_HOST}/cosmos/feegrant/v1beta1/issued/{granter}",
	"{LCD_HOST}/cosmos/nft/v1beta1/balance/{owner}/{class_id}",
	"{LCD_HOST}/cosmos/nft/v1beta1/classes",
	"{LCD_HOST}/cosmos/nft/v1beta1/classes/{class_id}",
	"{LCD_HOST}/cosmos/nft/v1beta1/nfts",
	"{LCD_HOST}/cosmos/nft/v1beta1/nfts/{class_id}/{id}",
	"{LCD_HOST}/cosmos/nft/v1beta1/owner/{class_id}/{id}",
	"{LCD_HOST}/cosmos/nft/v1beta1/supply/{class_id}",
	"{LCD_HOST}/cosmos/group/v1/group_info/{group_id}",
	"{LCD_HOST}/cosmos/group/v1/group_members/{group_id}",
	"{LCD_HOST}/cosmos/group/v1/group_policies_by_admin/{admin}",
	"{LCD_HOST}/cosmos/group/v1/group_policies_by_group/{group_id}",
	"{LCD_HOST}/cosmos/group/v1/group_policy_info/{address}",
	"{LCD_HOST}/cosmos/group/v1/groups",
	"{LCD_HOST}/cosmos/group/v1/groups_by_admin/{admin}",
	"{LCD_HOST}/cosmos/group/v1/groups_by_member/{address}",
	"{LCD_HOST}/cosmos/group/v1/proposal/{proposal_id}",
	"{LCD_HOST}/cosmos/group/v1/proposals/{proposal_id}/tally",
	"{LCD_HOST}/cosmos/group/v1/proposals_by_group_policy/{address}",
	"{LCD_HOST}/cosmos/group/v1/vote_by_proposal_voter/{proposal_id}/{voter}",
	"{LCD_HOST}/cosmos/group/v1/votes_by_proposal/{proposal_id}",
	"{LCD_HOST}/cosmos/group/v1/votes_by_voter/{voter}",
	"{RPC_HOST}/broadcast_tx_sync",
	"{RPC_HOST}/broadcast_tx_async",
	"{RPC_HOST}/broadcast_tx_commit",
	"{RPC_HOST}/check_tx",
	"{RPC_HOST}/subscribe",
	"{RPC_HOST}/unsubscribe",
	"{RPC_HOST}/unsubscribe_all",
	"{RPC_HOST}/health",
	"{RPC_HOST}/status",
	"{RPC_HOST}/net_info",
	"{RPC_HOST}/dial_seeds",
	"{RPC_HOST}/dial_peers",
	"{RPC_HOST}/blockchain",
	"{RPC_HOST}/block",
	"{RPC_HOST}/block_by_hash",
	"{RPC_HOST}/block_results",
	"{RPC_HOST}/commit",
	"{RPC_HOST}/validators",
	"{RPC_HOST}/genesis",
	"{RPC_HOST}/genesis_chunked",
	"{RPC_HOST}/dump_consensus_state",
	"{RPC_HOST}/consensus_state",
	"{RPC_HOST}/consensus_params",
	"{RPC_HOST}/unconfirmed_txs",
	"{RPC_HOST}/num_unconfirmed_txs",
	"{RPC_HOST}/tx_search",
	"{RPC_HOST}/block_search",
	"{RPC_HOST}/tx",
	"{RPC_HOST}/abci_info",
	"{RPC_HOST}/abci_query",
	"{RPC_HOST}/broadcast_evidence",
}

func benchmarkHTTP(url string, numRequests int) (float64, float64, float64, float64, error) {
	var totalDuration time.Duration
	var responseTimes []float64

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	for i := 0; i < numRequests; i++ {
		startTime := time.Now()

		resp, err := client.Get(url)
		if err != nil {
			return 0, 0, 0, 0, err
		}

		if !strings.HasPrefix(resp.Status, "2") {
			fmt.Println("HTTP Error: " + resp.Status)
			return 0, 0, 0, 0, errors.New(resp.Status)
		}

		defer resp.Body.Close()

		duration := time.Since(startTime)
		totalDuration += duration

		responseTimes = append(responseTimes, duration.Seconds())
	}

	// Calculate time per request
	timePerRequest := totalDuration.Seconds() / float64(numRequests)

	// Calculate tp50
	sort.Float64s(responseTimes)
	tp50Index := int(float64(len(responseTimes)) * 0.5)
	tp90Index := int(float64(len(responseTimes)) * 0.9)
	tp99Index := int(float64(len(responseTimes)) * 0.99)
	tp50 := responseTimes[tp50Index]
	tp90 := responseTimes[tp90Index]
	tp99 := responseTimes[tp99Index]

	return timePerRequest, tp50, tp90, tp99, nil
}

func main() {
	numRequests := 20

	for _, url := range ENDPOINTS {

		url = strings.Replace(url, "{LCD_HOST}", "https://api.cosmoshub.strange.love", -1)
		url = strings.Replace(url, "{RPC_HOST}", "https://rpc.cosmoshub.strange.love", -1)

		url = strings.Replace(url, "{address_bytes}", "cosmos1uyhsywympl67amu02y44xakknjcfhmjl9kx2fr", -1)
		url = strings.Replace(url, "{address_string}", "cosmos1uyhsywympl67amu02y44xakknjcfhmjl9kx2fr", -1)
		url = strings.Replace(url, "{address}", "cosmos1uyhsywympl67amu02y44xakknjcfhmjl9kx2fr", -1)
		url = strings.Replace(url, "{admin}", "cosmos1uyhsywympl67amu02y44xakknjcfhmjl9kx2fr", -1)        // TODO
		url = strings.Replace(url, "{class_id}", "cosmos1uyhsywympl67amu02y44xakknjcfhmjl9kx2fr", -1)     // TODO
		url = strings.Replace(url, "{cons_address}", "cosmos1uyhsywympl67amu02y44xakknjcfhmjl9kx2fr", -1) // TODO
		url = strings.Replace(url, "{delegator_addr}", "cosmos1u753tf5t6hq9l8fq7wtsfsxc89955mvgsj5krk", -1)
		url = strings.Replace(url, "{delegator_address}", "cosmos1u753tf5t6hq9l8fq7wtsfsxc89955mvgsj5krk", -1)
		url = strings.Replace(url, "{denom}", "uatom", -1)
		url = strings.Replace(url, "{depositor}", "cosmos1hxv7mpztvln45eghez6evw2ypcw4vjmsmr8cdx", -1)
		url = strings.Replace(url, "{grantee}", "cosmos1uyhsywympl67amu02y44xakknjcfhmjl9kx2fr", -1)  // TODO
		url = strings.Replace(url, "{granter}", "cosmos1uyhsywympl67amu02y44xakknjcfhmjl9kx2fr", -1)  // TODO
		url = strings.Replace(url, "{group_id}", "cosmos1uyhsywympl67amu02y44xakknjcfhmjl9kx2fr", -1) // TODO
		url = strings.Replace(url, "{hash}", "C436256EA929158B6926EAB9987964895433F119CA0DC3497DC9345D280FC8FE", -1)
		url = strings.Replace(url, "{height}", "14659652", -1)
		url = strings.Replace(url, "{last_height}", "14659652", -1)
		url = strings.Replace(url, "{name}", "cosmos1uyhsywympl67amu02y44xakknjcfhmjl9kx2fr", -1)        // TODO module account
		url = strings.Replace(url, "{owner}", "cosmos1uyhsywympl67amu02y44xakknjcfhmjl9kx2fr", -1)       // TODO
		url = strings.Replace(url, "{params_type}", "cosmos1uyhsywympl67amu02y44xakknjcfhmjl9kx2fr", -1) // TODO
		url = strings.Replace(url, "{proposal_id}", "69", -1)
		url = strings.Replace(url, "{validator_addr}", "cosmos130mdu9a0etmeuw52qfxk73pn0ga6gawkryh2z6", -1)
		url = strings.Replace(url, "{validator_address}", "cosmos130mdu9a0etmeuw52qfxk73pn0ga6gawkryh2z6", -1)
		url = strings.Replace(url, "{voter}", "cosmos15sfx9qfccwrnfp342dh2w8wx2jvfrln4pwwsgp", -1)

		fmt.Printf("Endpoint: %s\n", url)

		timePerRequest, tp50, tp90, tp99, err := benchmarkHTTP(url, numRequests)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Time per request: %.0f ms ", timePerRequest*1000)
		fmt.Printf("tp50: %.0f ms ", tp50*1000)
		fmt.Printf("tp90: %.0f ms ", tp90*1000)
		fmt.Printf("tp99: %.0f ms\n\n", tp99*1000)
	}

}

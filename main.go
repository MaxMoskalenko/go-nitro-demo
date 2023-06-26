package main

import (
	"fmt"
	"math/big"
	"time"

	"go-nitro-demo/peer"
	"go-nitro-demo/stage"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

var (
	contractAddress = common.HexToAddress("0xD7b7829D9b9a2c362AF2500B5Fe66014e6a91D8c")
	assetAddress    = common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
)

func main() {
	alice := peer.New("0x518A41356aB936b2e8cCFE4C179dF460C83000AD", "0a142fd3bb44bbf5091ccf39e853e9f7cedfbb937b7dc22cfea1bfae011d1937")
	bob := peer.New("0xd92B92cc0b45CC1191071c705aEde07babd69281", "e78a8e3cd440035f9e8bfda19050907afd611908b457e5c2855e86e05f567258")

	fixedPart := state.FixedPart{
		Participants:      []common.Address{alice.PublicKey, bob.PublicKey},
		ChannelNonce:      123,
		AppDefinition:     contractAddress,
		ChallengeDuration: uint32(time.Hour.Milliseconds()),
	}

	initialOutcome := outcome.Exit{
		outcome.SingleAssetExit{
			Asset: assetAddress,
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: types.AddressToDestination(alice.PublicKey),
					Amount:      big.NewInt(5),
				},
				outcome.Allocation{
					Destination: types.AddressToDestination(bob.PublicKey),
					Amount:      big.NewInt(5),
				},
			},
		},
	}

	variablePart := state.VariablePart{
		Outcome: initialOutcome,
		TurnNum: 0,
		IsFinal: false,
	}

	var err error

	// Both peers create their versions of channel
	err = alice.CreateChannel(fixedPart, variablePart)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(alice.Channel.Id)

	err = bob.CreateChannel(fixedPart, variablePart)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Both peers exchange their agreements to open the channel
	err = stage.PreFundStage(&alice, &bob)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Alice and Bob deposit to adjudicator here

	// Both peers verify that funding by another party was correct
	err = stage.PostFundStage(&alice, &bob)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Alice decides to change outcome of the deal and propose Bob to sign a new state
	newState, err := alice.Channel.LatestSupportedState()
	if err != nil {
		fmt.Println(err)
		return
	}

	newState.Outcome = outcome.Exit{
		outcome.SingleAssetExit{
			Asset: assetAddress,
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: types.AddressToDestination(alice.PublicKey),
					Amount:      big.NewInt(6),
				},
				outcome.Allocation{
					Destination: types.AddressToDestination(bob.PublicKey),
					Amount:      big.NewInt(4),
				},
			},
		},
	}
	newState.TurnNum += 1

	err = stage.AddNewStateAndAgree(&alice, &bob, newState)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Alice proposes to finalize channel
	newState, err = alice.Channel.LatestSupportedState()
	if err != nil {
		fmt.Println(err)
		return
	}

	newState.IsFinal = true
	newState.TurnNum += 1

	err = stage.AddNewStateAndAgree(&alice, &bob, newState)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("alice and bob prefunds %+v %+v\n", alice.Channel.PreFundComplete(), bob.Channel.PreFundComplete())
	fmt.Printf("alice and bob postfund %+v %+v\n", alice.Channel.PostFundComplete(), bob.Channel.PostFundComplete())
	fmt.Printf("alice and bob finalized %+v %+v\n", alice.Channel.FinalCompleted(), bob.Channel.FinalCompleted())
}

package stage

import (
	"go-nitro-demo/peer"

	"github.com/statechannels/go-nitro/channel/state"
)

func PreFundStage(alice *peer.Peer, bob *peer.Peer) error {
	// Alice signs prefund state and "sends" signed state to Bob
	ssa, err := alice.Channel.SignAndAddPrefund(&alice.PrivateKey)
	if err != nil {
		return err
	}

	// Bob saves signed state by Alice, countersigns it and sends back his signature to Alice
	bob.Channel.AddSignedState(ssa)
	ssb, err := bob.Channel.SignAndAddPrefund(&bob.PrivateKey)
	if err != nil {
		return err
	}

	// Alice saves Bob's signature
	alice.Channel.AddSignedState(ssb)
	return nil
}

func PostFundStage(alice *peer.Peer, bob *peer.Peer) error {
	// Alice signs postfund state and "sends" signed state to Bob
	ssa, err := alice.Channel.SignAndAddPostfund(&alice.PrivateKey)
	if err != nil {
		return err
	}

	// Bob saves signed state by Alice, countersigns it and sends back his signature to Alice
	bob.Channel.AddSignedState(ssa)
	ssb, err := bob.Channel.SignAndAddPostfund(&bob.PrivateKey)
	if err != nil {
		return err
	}

	// Alice saves Bob's signature
	alice.Channel.AddSignedState(ssb)
	return nil
}

func AddNewStateAndAgree(alice *peer.Peer, bob *peer.Peer, s state.State) error {
	// Alice signs new state state and "sends" signed state to Bob
	ssa, err := alice.Channel.SignAndAddState(s, &alice.PrivateKey)
	if err != nil {
		return err
	}

	// Bob saves signed state by Alice, countersigns it and sends back his signature to Alice
	bob.Channel.AddSignedState(ssa)
	ssb, err := bob.Channel.SignAndAddState(ssa.State(), &bob.PrivateKey)
	if err != nil {
		return err
	}

	// Alice saves Bob's signature
	alice.Channel.AddSignedState(ssb)
	return nil
}

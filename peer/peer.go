package peer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
)

type Peer struct {
	PublicKey  common.Address
	PrivateKey []byte
	Channel    *channel.Channel
}

func New(pub string, priv string) Peer {
	return Peer{
		PublicKey:  common.HexToAddress(pub),
		PrivateKey: common.Hex2Bytes(priv),
	}
}

func (p *Peer) CreateChannel(fp state.FixedPart, vp state.VariablePart) error {
	state := state.StateFromFixedAndVariablePart(fp, vp)
	ch, err := channel.New(state, 0)
	if err != nil {
		return err
	}

	p.Channel = ch
	return nil
}

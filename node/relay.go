package node

import (
	"fmt"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/nknorg/nkn/block"
	"github.com/nknorg/nkn/chain"
	"github.com/nknorg/nkn/event"
	"github.com/nknorg/nkn/pb"
	"github.com/nknorg/nkn/por"
	"github.com/nknorg/nkn/transaction"
	"github.com/nknorg/nkn/util/address"
	"github.com/nknorg/nkn/util/config"
	"github.com/nknorg/nkn/util/log"
	"github.com/nknorg/nkn/vault"
)

type RelayService struct {
	sync.Mutex
	wallet    *vault.Wallet
	localNode *LocalNode
	porServer *por.PorServer
}

func NewRelayService(wallet *vault.Wallet, localNode *LocalNode) *RelayService {
	service := &RelayService{
		wallet:    wallet,
		localNode: localNode,
		porServer: por.GetPorServer(),
	}
	return service
}

func (rs *RelayService) Start() error {
	event.Queue.Subscribe(event.NewBlockProduced, rs.populateVRFCache)
	event.Queue.Subscribe(event.NewBlockProduced, rs.flushSigChain)
	event.Queue.Subscribe(event.PinSigChain, rs.startPinSigChain)
	event.Queue.Subscribe(event.BacktrackSigChain, rs.backtrackDestSigChain)
	rs.localNode.AddMessageHandler(pb.RELAY, rs.relayMessageHandler)
	rs.localNode.AddMessageHandler(pb.PIN_SIGNATURE_CHAIN, rs.pinSigChainMessageHandler)
	rs.localNode.AddMessageHandler(pb.BACKTRACK_SIGNATURE_CHAIN, rs.backtrackSigChainMessageHandler)
	return nil
}

// NewRelayMessage creates a RELAY message
func NewRelayMessage(srcIdentifier string, srcPubkey, destID, payload, blockHash, lastHash []byte, maxHoldingSeconds uint32) (*pb.UnsignedMessage, error) {
	msgBody := &pb.Relay{
		SrcIdentifier:     srcIdentifier,
		SrcPubkey:         srcPubkey,
		DestId:            destID,
		Payload:           payload,
		MaxHoldingSeconds: maxHoldingSeconds,
		BlockHash:         blockHash,
		LastHash:          lastHash,
		SigChainLen:       1,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pb.UnsignedMessage{
		MessageType: pb.RELAY,
		Message:     buf,
	}

	return msg, nil
}

// NewPinSigChainMessage creates a PIN_SIGNATURE_CHAIN message
func NewPinSigChainMessage(hash []byte) (*pb.UnsignedMessage, error) {
	msgBody := &pb.PinSignatureChain{
		Hash: hash,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pb.UnsignedMessage{
		MessageType: pb.PIN_SIGNATURE_CHAIN,
		Message:     buf,
	}

	return msg, nil
}

// NewBacktrackSigChainMessage creates a BACKTRACK_SIGNATURE_CHAIN message
func NewBacktrackSigChainMessage(sigChainElems []*pb.SigChainElem, hash []byte) (*pb.UnsignedMessage, error) {
	msgBody := &pb.BacktrackSignatureChain{
		SigChainElems: sigChainElems,
		Hash:          hash,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pb.UnsignedMessage{
		MessageType: pb.BACKTRACK_SIGNATURE_CHAIN,
		Message:     buf,
	}

	return msg, nil
}

// relayMessageHandler handles a RELAY message
func (rs *RelayService) relayMessageHandler(remoteMessage *RemoteMessage) ([]byte, bool, error) {
	msgBody := &pb.Relay{}
	err := proto.Unmarshal(remoteMessage.Message, msgBody)
	if err != nil {
		return nil, false, err
	}

	event.Queue.Notify(event.SendInboundMessageToClient, msgBody)

	return nil, false, nil
}

// pinSigChainMessageHandler handles a PIN_SIGNATURE_CHAIN message
func (rs *RelayService) pinSigChainMessageHandler(remoteMessage *RemoteMessage) ([]byte, bool, error) {
	msgBody := &pb.PinSignatureChain{}
	err := proto.Unmarshal(remoteMessage.Message, msgBody)
	if err != nil {
		return nil, false, err
	}

	err = rs.pinSigChain(msgBody.Hash, remoteMessage.Sender.PublicKey)
	if err != nil {
		return nil, false, err
	}

	return nil, false, nil
}

// backtrackSigChainMessageHandler handles a BACKTRACK_SIGNATURE_CHAIN message
func (rs *RelayService) backtrackSigChainMessageHandler(remoteMessage *RemoteMessage) ([]byte, bool, error) {
	msgBody := &pb.BacktrackSignatureChain{}
	err := proto.Unmarshal(remoteMessage.Message, msgBody)
	if err != nil {
		return nil, false, err
	}

	err = rs.backtrackSigChain(msgBody.SigChainElems, msgBody.Hash, remoteMessage.Sender.PublicKey)
	if err != nil {
		return nil, false, err
	}

	return nil, false, nil
}

func (rs *RelayService) pinSigChain(hash, senderPubkey []byte) error {
	prevHash, prevNodeID, err := rs.porServer.PinSigChain(hash, senderPubkey)
	if err != nil {
		return err
	}

	if prevNodeID == nil {
		err = rs.porServer.PinSrcSigChain(prevHash)
		if err != nil {
			return err
		}
	} else {
		nextHop := rs.localNode.GetNbrNode(chordIDToNodeID(prevNodeID))
		if nextHop == nil {
			return fmt.Errorf("cannot find next hop with id %x", prevNodeID)
		}

		nextMsg, err := NewPinSigChainMessage(prevHash)
		if err != nil {
			return err
		}

		buf, err := rs.localNode.SerializeMessage(nextMsg, false)
		if err != nil {
			return err
		}

		err = nextHop.SendBytesAsync(buf)
		if err != nil {
			return err
		}
	}

	rs.porServer.PinSigChainSuccess(hash)

	return nil
}

func (rs *RelayService) backtrackSigChain(sigChainElems []*pb.SigChainElem, hash, senderPubkey []byte) error {
	sigChainElems, prevHash, prevNodeID, err := rs.porServer.BacktrackSigChain(sigChainElems, hash, senderPubkey)
	if err != nil {
		return err
	}

	if prevNodeID == nil {
		sigChain, err := rs.porServer.PopSrcSigChainFromCache(prevHash)
		if err != nil {
			return err
		}

		sigChain.Elems = append(sigChain.Elems, sigChainElems...)

		err = rs.broadcastSigChain(sigChain)
		if err != nil {
			return err
		}
	} else {
		nextHop := rs.localNode.GetNbrNode(chordIDToNodeID(prevNodeID))
		if nextHop == nil {
			return fmt.Errorf("cannot find next hop with id %x", prevNodeID)
		}

		nextMsg, err := NewBacktrackSigChainMessage(sigChainElems, prevHash)
		if err != nil {
			return err
		}

		buf, err := rs.localNode.SerializeMessage(nextMsg, false)
		if err != nil {
			return err
		}

		err = nextHop.SendBytesAsync(buf)
		if err != nil {
			return err
		}
	}

	rs.porServer.BacktrackSigChainSuccess(hash)

	return nil
}

func (rs *RelayService) broadcastSigChain(sigChain *pb.SigChain) error {
	buf, err := proto.Marshal(sigChain)
	if err != nil {
		return err
	}

	txn, err := MakeSigChainTransaction(rs.wallet, buf)
	if err != nil {
		return err
	}

	currentHeight := chain.DefaultLedger.Store.GetHeight()

	err = chain.VerifyTransaction(txn, currentHeight+1)
	if err != nil {
		return err
	}

	porPkg, err := por.GetPorServer().AddSigChainFromTx(txn, currentHeight)
	if err != nil {
		return err
	}
	if porPkg == nil {
		return nil
	}

	err = rs.localNode.iHaveSignatureChainTransaction(porPkg.VoteForHeight, porPkg.SigHash, nil)
	if err != nil {
		return err
	}

	return nil
}

func (rs *RelayService) startPinSigChain(v interface{}) {
	sigChainInfo, ok := v.(*por.PinSigChainInfo)
	if !ok {
		log.Error("Decode pin sigchain info failed")
		return
	}

	err := rs.pinSigChain(sigChainInfo.PrevHash, nil)
	if err != nil {
		log.Errorf("Pin sigchain error: %v", err)
	}
}

func (rs *RelayService) backtrackDestSigChain(v interface{}) {
	sigChainInfo, ok := v.(*por.BacktrackSigChainInfo)
	if !ok {
		log.Error("Decode backtrack sigchain info failed")
		return
	}

	sigChainElems := []*pb.SigChainElem{sigChainInfo.DestSigChainElem}
	err := rs.backtrackSigChain(sigChainElems, sigChainInfo.PrevHash, nil)
	if err != nil {
		log.Errorf("Backtrack sigchain error: %v", err)
	}
}

func (rs *RelayService) updateRelayMessage(relayMessage *pb.Relay, nextHop, prevHop *RemoteNode) error {
	var nextPubkey []byte
	if nextHop != nil {
		nextPubkey = nextHop.GetPubKey()
	}

	mining := config.Parameters.Mining && rs.localNode.GetSyncState() == pb.PERSIST_FINISHED

	var prevNodeID []byte
	if prevHop != nil {
		prevNodeID = prevHop.Id
	}

	return rs.porServer.UpdateRelayMessage(relayMessage, nextPubkey, prevNodeID, mining)
}

func (localNode *LocalNode) startRelayer() {
	localNode.relayer.Start()
}

func (localNode *LocalNode) SendRelayMessage(srcAddr, destAddr string, payload, signature, blockHash []byte, nonce, maxHoldingSeconds uint32) error {
	srcID, srcPubkey, srcIdentifier, err := address.ParseClientAddress(srcAddr)
	if err != nil {
		return err
	}

	destID, destPubkey, _, err := address.ParseClientAddress(destAddr)
	if err != nil {
		return err
	}

	_, err = por.GetPorServer().CreateSigChainForClient(
		nonce,
		uint32(len(payload)),
		blockHash,
		srcID,
		srcPubkey,
		destID,
		destPubkey,
		signature,
		pb.SIGNATURE,
	)
	if err != nil {
		return err
	}

	msg, err := NewRelayMessage(srcIdentifier, srcPubkey, destID, payload, blockHash, signature, maxHoldingSeconds)
	if err != nil {
		return err
	}

	buf, err := localNode.SerializeMessage(msg, false)
	if err != nil {
		return err
	}

	_, err = localNode.nnet.SendBytesRelayAsync(buf, destID)
	if err != nil {
		return err
	}

	return nil
}

func MakeSigChainTransaction(wallet *vault.Wallet, sigChain []byte) (*transaction.Transaction, error) {
	account, err := wallet.GetDefaultAccount()
	if err != nil {
		return nil, err
	}
	txn, err := transaction.NewSigChainTransaction(sigChain, account.ProgramHash, 0)
	if err != nil {
		return nil, err
	}

	// sign transaction contract
	err = wallet.Sign(txn)
	if err != nil {
		return nil, err
	}

	return txn, nil
}

func (rs *RelayService) populateVRFCache(v interface{}) {
	block, ok := v.(*block.Block)
	if !ok {
		return
	}

	blockHash := block.Hash()
	rs.porServer.GetOrComputeVrf(blockHash.ToArray())
}

func (rs *RelayService) flushSigChain(v interface{}) {
	block, ok := v.(*block.Block)
	if !ok {
		return
	}

	height := block.Header.UnsignedHeader.Height - config.SigChainBlockDelay - 1
	if height < 0 {
		height = 0
	}
	blockHash := chain.DefaultLedger.Store.GetHeaderHashByHeight(height)

	rs.porServer.FlushSigChain(blockHash.ToArray())
}

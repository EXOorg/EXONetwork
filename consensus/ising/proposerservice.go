package ising

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	. "github.com/nknorg/nkn/common"
	"github.com/nknorg/nkn/consensus/ising/voting"
	"github.com/nknorg/nkn/core/ledger"
	"github.com/nknorg/nkn/core/transaction"
	"github.com/nknorg/nkn/crypto"
	"github.com/nknorg/nkn/events"
	"github.com/nknorg/nkn/net/message"
	"github.com/nknorg/nkn/net/protocol"
	"github.com/nknorg/nkn/util/config"
	"github.com/nknorg/nkn/util/log"
	"github.com/nknorg/nkn/vault"
)

const (
	TxnAmountToBePackaged      = 1024
	WaitingForFloodingFinished = time.Second * 1
	WaitingForVotingFinished   = time.Second * 8
	TimeoutTolerance           = time.Second * 2
)

type ProposerService struct {
	sync.RWMutex
	account              *vault.Account            // local account
	timer                *time.Timer               // timer for proposer node
	timeout              *time.Timer               // timeout for next round consensus
	proposerChangeTimer  *time.Timer               // timer for proposer change
	localNode            protocol.Noder            // local node
	txnCollector         *transaction.TxnCollector // collect transaction from where
	mining               Mining                    // built-in mining
	msgChan              chan interface{}          // get notice from probe thread
	consensusMsgReceived events.Subscriber         // consensus events listening
	blockPersisted       events.Subscriber         // block saved events
	syncFinished         events.Subscriber         // block syncing finished event
	proposerCache        *ProposerCache            // cache for block proposer
	syncCache            *SyncCache                // cache for block syncing
	voting               []voting.Voting           // array for sigchain and block voting
}

func NewProposerService(account *vault.Account, node protocol.Noder) *ProposerService {
	txnCollector := transaction.NewTxnCollector(node.GetTxnPool(), TxnAmountToBePackaged)

	service := &ProposerService{
		timer:               time.NewTimer(config.ConsensusTime),
		timeout:             time.NewTimer(config.ConsensusTime + TimeoutTolerance),
		proposerChangeTimer: time.NewTimer(config.ProposerChangeTime),
		account:             account,
		localNode:           node,
		txnCollector:        txnCollector,
		mining:              NewBuiltinMining(account, txnCollector),
		msgChan:             make(chan interface{}, MsgChanCap),
		syncCache:           NewSyncBlockCache(),
		proposerCache:       NewProposerCache(),
		voting: []voting.Voting{
			voting.NewBlockVoting(),
			voting.NewSigChainVoting(txnCollector),
		},
	}
	if !service.timer.Stop() {
		<-service.timer.C
	}
	if !service.timeout.Stop() {
		<-service.timeout.C
	}

	return service
}

func (ps *ProposerService) CurrentVoting(vType voting.VotingContentType) voting.Voting {
	for _, v := range ps.voting {
		if v.VotingType() == vType {
			return v
		}
	}

	return nil
}

func (ps *ProposerService) ConsensusRoutine(vType voting.VotingContentType, isProposer bool) {
	// initialization for height and voting weight
	ps.Initialize(vType)
	current := ps.CurrentVoting(vType)
	votingHeight := current.GetVotingHeight()
	votingPool := current.GetVotingPool()

	go func() {
		// waiting for flooding finished
		time.Sleep(WaitingForFloodingFinished)
		// send new proposal
		err := ps.SendNewProposal(votingHeight, vType, isProposer)
		if err != nil {
			log.Info("waiting for receiving proposed entity...")
			return
		}
	}()

	go func() {
		// waiting for voting finished
		time.Sleep(WaitingForVotingFinished)
		finalHash, ok := votingPool.GetMind(votingHeight)
		if !ok {
			return
		}
		// get the final entity from local cache or database
		content, err := current.GetVotingContent(finalHash, votingHeight)
		if err != nil {
			log.Errorf("get final entity error, hash: %s, type: %d, votingHeight: %d",
				BytesToHexString(finalHash.ToArrayReverse()), vType, votingHeight)
			log.Warn(err)
			return
		}
		// process final block and signature chain
		switch vType {
		case voting.SigChainTxnVote:
			if txn, ok := content.(*transaction.Transaction); ok {
				ps.proposerCache.Add(votingHeight, txn)
			}
		case voting.BlockVote:
			if block, ok := content.(*ledger.Block); ok {
				err = ledger.DefaultLedger.Blockchain.AddBlock(block)
				if err != nil {
					log.Error("saving block error: ", err)
					return
				}
			}
		}
	}()
}

// GetReceiverNode returns neighbors nodes according to neighbor node ID passed in.
// If 'nids' passed in is nil then returns all neighbor nodes.
func (ps *ProposerService) GetReceiverNode(nids []uint64) []protocol.Noder {
	if nids == nil {
		return ps.localNode.GetNeighborNoder()
	}
	var nodes []protocol.Noder
	neighbors := ps.localNode.GetNeighborNoder()
	for _, id := range nids {
		for _, node := range neighbors {
			if id == node.GetID() {
				nodes = append(nodes, node)
			}
		}
	}

	return nodes
}

func (ps *ProposerService) SendNewProposal(votingHeight uint32, vType voting.VotingContentType, isProposer bool) error {
	current := ps.CurrentVoting(vType)
	content, err := current.GetBestVotingContent(votingHeight)
	if err != nil {
		return err
	}
	if !isProposer && !current.VerifyVotingContent(content) {
		return errors.New("verify voting content error when send new proposal")
	}
	votingPool := current.GetVotingPool()
	ownMindHash := content.Hash()
	ownWeight, _ := ledger.DefaultLedger.Store.GetVotingWeight(Uint160{})
	ownNodeID := ps.localNode.GetID()
	log.Infof("own mind: %s, type: %d", BytesToHexString(ownMindHash.ToArrayReverse()), vType)
	if mind, ok := votingPool.GetMind(votingHeight); ok {
		log.Infof("neighbor mind: %s, type: %d ", BytesToHexString(mind.ToArrayReverse()), vType)
		// if own mind different with neighbors then change mind
		if ownMindHash.CompareTo(mind) != 0 {
			ownMindHash = mind
			log.Infof("own mind is affected by neighbor mind: %s, type: %d",
				BytesToHexString(ownMindHash.ToArrayReverse()), vType)
		}
		// add own vote to voting pool
		votingPool.AddToReceivePool(votingHeight, ownNodeID, ownWeight, ownMindHash)
	} else {
		// set initial mind if it has not been set
		votingPool.SetMind(votingHeight, ownMindHash)

		// if voting result is different with initial mind then change mind
		maybeFinalHash, _ := votingPool.AddVoteThenCounting(votingHeight, ownNodeID, ownWeight, ownMindHash)
		if maybeFinalHash != nil && maybeFinalHash.CompareTo(ownMindHash) != 0 {
			log.Infof("mind set when send proposal: %s, type: %d",
				BytesToHexString(ownMindHash.ToArrayReverse()), vType)
			votingPool.SetMind(votingHeight, ownMindHash)
		}
	}
	if !current.HasSelfState(ownMindHash, voting.ProposalSent) {
		log.Infof("proposing hash: %s, type: %d", BytesToHexString(ownMindHash.ToArrayReverse()), vType)
		// create new proposal
		proposalMsg := NewProposal(&ownMindHash, votingHeight, vType)
		// get nodes which should receive proposal message
		nodes := ps.GetReceiverNode(nil)
		// send proposal to neighbors
		ps.SendConsensusMsg(proposalMsg, nodes)
		// state changed for current hash
		current.SetSelfState(ownMindHash, voting.ProposalSent)
		// set confirming hash
		current.SetConfirmingHash(ownMindHash)
	}

	return nil
}

func (ps *ProposerService) ProduceNewBlock() {
	current := ps.CurrentVoting(voting.BlockVote)
	votingPool := current.GetVotingPool()
	votingHeight := current.GetVotingHeight()
	proposerInfo := ps.proposerCache.Get(votingHeight + 1)
	if proposerInfo == nil {
		proposerInfo = &ProposerInfo{
			winningHash:     EmptyUint256,
			winningHashType: ledger.WinningBlockHash,
		}
	}
	// build new block to be proposed
	block, err := ps.mining.BuildBlock(votingHeight, proposerInfo.winningHash, proposerInfo.winningHashType)
	if err != nil {
		log.Error("building block error: ", err)
	}
	err = current.AddToCache(block)
	if err != nil {
		log.Error("adding block to cache error: ", err)
		return
	}
	err = ledger.TransactionCheck(block)
	if err != nil {
		log.Error("found invalide transaction when produce new block")
		return
	}
	// generate BlockFlooding message
	blockFlooding := NewBlockFlooding(block)
	// get nodes which should receive this message
	nodes := ps.GetReceiverNode(nil)
	// send BlockFlooding message
	err = ps.SendConsensusMsg(blockFlooding, nodes)
	if err != nil {
		log.Error("sending consensus message error: ", err)
	}
	blockHash := block.Hash()
	// set initial mind for block proposer
	votingPool.SetMind(votingHeight, blockHash)
	log.Info("when produce new block mind set to: ", BytesToHexString(blockHash.ToArrayReverse()))
}

func (ps *ProposerService) IsBlockProposer() bool {
	localPublicKey, err := ps.account.PublicKey.EncodePoint(true)
	if err != nil {
		return false
	}
	current := ps.CurrentVoting(voting.BlockVote)
	votingHeight := current.GetVotingHeight()
	proposerInfo := ps.proposerCache.Get(votingHeight)
	if proposerInfo == nil {
		return false
	}
	if !IsEqualBytes(localPublicKey, proposerInfo.publicKey) {
		return false
	}

	return true
}

func (ps *ProposerService) ProposerRoutine() {
	for {
		select {
		case <-ps.timer.C:
			if ps.IsBlockProposer() {
				log.Info("-> I am Block Proposer")
				if ps.localNode.GetSyncState() != protocol.PersistFinished {
					ps.localNode.StopSyncBlock()
				}
				ps.ProduceNewBlock()
				for _, v := range ps.voting {
					go ps.ConsensusRoutine(v.VotingType(), true)
				}
				ps.timer.Reset(config.ConsensusTime)
			}
		}
	}
}

func (ps *ProposerService) TimeoutRoutine() {
	for {
		select {
		case <-ps.timeout.C:
			ps.timer.Stop()
			ps.timer.Reset(0)
		}
	}
}

func (ps *ProposerService) ProbeRoutine() {
	for {
		select {
		case msg := <-ps.msgChan:
			if notice, ok := msg.(*Notice); ok {
				heightHashMap := make(map[uint32]Uint256)
				heightNeighborMap := make(map[uint32]uint64)
				for k, v := range notice.BlockHistory {
					height, hash := StringToHeightHash(k)
					heightHashMap[height] = hash
					heightNeighborMap[height] = v
				}
				if height, ok := ledger.DefaultLedger.Store.CheckBlockHistory(heightHashMap); !ok {
					//TODO DB reverts to 'height' - 1 and request blocks from neighbor n[height]
					_ = height
				}
			}
		}
	}
}

func (ps *ProposerService) BlockPersistCompleted(v interface{}) {
	// reset timer when block persisted
	ps.proposerChangeTimer.Stop()
	ps.proposerChangeTimer.Reset(config.ProposerChangeTime)

	if ps.localNode.GetSyncState() == protocol.PersistFinished {
		if block, ok := v.(*ledger.Block); ok {
			// record time when persist block
			ledger.DefaultLedger.Blockchain.AddBlockTime(block.Hash(), time.Now().Unix())
			ps.txnCollector.Cleanup(block.Transactions)
		}
	}
}

func (ps *ProposerService) ChangeProposerRoutine() {
	for {
		select {
		case <-ps.proposerChangeTimer.C:
			now := time.Now().Unix()
			currentHeight := ledger.DefaultLedger.Store.GetHeight()
			var block *ledger.Block
			var err error
			if currentHeight < InitialBlockHeight {
				block, err = ledger.DefaultLedger.Store.GetBlockByHeight(0)
				if err != nil {
					log.Error("get genesis block error when change proposer")
				}
			} else {
				currentBlock, err := ledger.DefaultLedger.Store.GetBlockByHeight(currentHeight)
				if err != nil {
					log.Errorf("get latest block %d error when change proposer", currentHeight)
				}
				timestamp := currentBlock.Header.Timestamp
				index := (now - timestamp) / 60
				var height uint32 = 0
				if int64(currentHeight) > index {
					height = uint32(int64(currentHeight) - index)
				}
				block, err = ledger.DefaultLedger.Store.GetBlockByHeight(height)
				if err != nil {
					log.Errorf("get block %d error when change proposer", currentHeight)
				}
			}
			ps.proposerCache.Add(currentHeight+1, block)
			ps.timer.Stop()
			ps.timer.Reset(0)
			ps.proposerChangeTimer.Reset(config.ProposerChangeTime)
		}
	}
}

func (ps *ProposerService) BlockSyncingFinished(v interface{}) {
	log.Infof("process blocks, from height: %d, to height: %d, consensus height: %d, start height: %d",
		ps.syncCache.startHeight, ps.syncCache.nextHeight-1,
		ps.syncCache.consensusHeight, ps.syncCache.consensusHeight+1)
	for i := ps.syncCache.startHeight; i < ps.syncCache.nextHeight; i++ {
		if i > ps.syncCache.consensusHeight {
			vBlock, err := ps.syncCache.GetBlockFromSyncCache(i)
			if err != nil {
				//TODO: if found ambiguous block then re-sync block
				log.Error("persist cached block error: ", err)
				return
			}
			err = ledger.HeaderCheck(vBlock.Block.Header, vBlock.ReceiveTime)
			if err != nil {
				log.Error("verifying header of cached block error: ", err)
				return
			}
			err = ledger.TransactionCheck(vBlock.Block)
			if err != nil {
				log.Error("verifying transaction in cached block error: ", err)
				return
			}
			err = ledger.DefaultLedger.Blockchain.AddBlock(vBlock.Block)
			if err != nil {
				log.Error("saving cached block error: ", err)
				return
			}
			// Fixme: wait for block persisted
			time.Sleep(time.Millisecond * 300)
		}
		err := ps.syncCache.RemoveBlockFromCache(i)
		if err != nil {
			log.Warnf("sync cache cleanup failed for height %d, error: %v", i, err)
		}
		ps.syncCache.timeLock.RemoveForHeight(i)
	}
	log.Info("cached block saving finished")
	// switch syncing state
	ps.localNode.SetSyncState(protocol.PersistFinished)
}
func (ps *ProposerService) SyncBlock(isProposer bool) {
	if ps.localNode.GetSyncState() == protocol.PersistFinished ||
		ps.localNode.GetSyncState() == protocol.SyncFinished {
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	// start block syncing
	go func() {
		defer wg.Done()
		ps.localNode.SyncBlock(isProposer)
	}()
	wg.Add(1)
	// start monitor routine for block syncing
	go func() {
		defer wg.Done()
		ps.localNode.SyncBlockMonitor(isProposer)
	}()
	wg.Wait()
}

func (ps *ProposerService) Start() error {
	// register consensus message
	ps.consensusMsgReceived = ps.localNode.GetEvent("consensus").Subscribe(events.EventConsensusMsgReceived,
		ps.ReceiveConsensusMsg)
	// register block saving event
	ps.blockPersisted = ledger.DefaultLedger.Blockchain.BCEvents.Subscribe(events.EventBlockPersistCompleted,
		ps.BlockPersistCompleted)
	// register block syncing event
	ps.syncFinished = ps.localNode.GetEvent("sync").Subscribe(events.EventBlockSyncingFinished,
		ps.BlockSyncingFinished)

	// start block proposer routine
	go ps.ProposerRoutine()
	// start change proposer routine
	go ps.ChangeProposerRoutine()
	// start timeout routine
	go ps.TimeoutRoutine()

	ps.SyncBlock(false)
	// start probe routine
	go ps.ProbeRoutine()

	return nil
}

func (ps *ProposerService) SendConsensusMsg(msg IsingMessage, to []protocol.Noder) error {
	isingPld, err := BuildIsingPayload(msg, ps.account.PublicKey)
	if err != nil {
		return err
	}
	hash, err := isingPld.DataHash()
	if err != nil {
		return err
	}
	signature, err := crypto.Sign(ps.account.PrivateKey, hash)
	if err != nil {
		return err
	}
	isingPld.Signature = signature

	for _, node := range to {
		b, err := message.NewIsingConsensus(isingPld)
		if err != nil {
			return err
		}
		node.Tx(b)
	}

	return nil
}

func (ps *ProposerService) ReceiveConsensusMsg(v interface{}) {
	if payload, ok := v.(*message.IsingPayload); ok {
		sender := payload.Sender
		signature := payload.Signature
		hash, err := payload.DataHash()
		if err != nil {
			fmt.Println("get consensus payload hash error")
			return
		}
		err = crypto.Verify(*sender, hash, signature)
		if err != nil {
			fmt.Println("consensus message verification error")
			return
		}
		isingMsg, err := RecoverFromIsingPayload(payload)
		if err != nil {
			fmt.Println("Deserialization of ising message error")
			return
		}
		switch t := isingMsg.(type) {
		case *BlockFlooding:
			ps.HandleBlockFloodingMsg(t, sender)
		case *Request:
			ps.HandleRequestMsg(t, sender)
		case *Response:
			ps.HandleResponseMsg(t, sender)
		case *StateProbe:
			ps.HandleStateProbeMsg(t, sender)
		case *Proposal:
			ps.HandleProposalMsg(t, sender)
		case *MindChanging:
			ps.HandleMindChangingMsg(t, sender)
		}
	}
}

func (ps *ProposerService) HandleBlockFloodingMsg(bfMsg *BlockFlooding, sender *crypto.PubKey) {
	current := ps.CurrentVoting(voting.BlockVote)
	votingHeight := current.GetVotingHeight()
	blockHash := bfMsg.block.Hash()
	height := bfMsg.block.Header.Height

	// returns if receive duplicate block
	if current.HasSelfState(blockHash, voting.FloodingFinished) {
		return
	}
	// set state for flooding block
	current.SetSelfState(blockHash, voting.FloodingFinished)

	// relay block to neighbors
	var nodes []protocol.Noder
	senderID := publickKeyToNodeID(sender)
	for _, node := range ps.localNode.GetNeighborNoder() {
		if node.GetID() != senderID {
			nodes = append(nodes, node)
		}
	}
	err := ps.SendConsensusMsg(bfMsg, nodes)
	if err != nil {
		log.Error("broadcast block message error: ", err)
	}

	// if block syncing is not finished, cache received blocks in order
	if ps.localNode.GetSyncState() != protocol.PersistFinished {
		err = ps.syncCache.AddBlockToSyncCache(bfMsg.block)
		if err != nil {
			log.Error("add received block to sync cache error: ", err)
		}
		log.Infof("cached block: %s, block height: %d,  totally cached: %d",
			BytesToHexString(blockHash.ToArrayReverse()), height, ps.syncCache.CachedBlockHeight())
		return
	}

	// expect the height of received block is equal to voting height when block syncing finished
	if height != votingHeight {
		log.Warnf("receive block which height is invalid, consensus height: %d, received block height: %d,"+
			" hash: %s\n", votingHeight, height, BytesToHexString(blockHash.ToArrayReverse()))
		return
	}
	err = current.AddToCache(bfMsg.block)
	if err != nil {
		log.Error("add received block to local cache error")
		return
	}

	// trigger consensus when receive appropriate block
	if !ps.IsBlockProposer() {
		for _, v := range ps.voting {
			go ps.ConsensusRoutine(v.VotingType(), false)
		}
		// trigger block proposer changed when tolerance time not receive block
		ps.timeout.Reset(config.ConsensusTime + TimeoutTolerance)
	}
}

func (ps *ProposerService) HandleRequestMsg(req *Request, sender *crypto.PubKey) {
	if ps.localNode.GetSyncState() != protocol.PersistFinished {
		return
	}
	current := ps.CurrentVoting(req.contentType)
	votingHeight := current.GetVotingHeight()
	hash := *req.hash
	height := req.height
	if height < votingHeight {
		log.Warnf("receive invalid request, consensus height: %d, request height: %d,"+
			" hash: %s\n", votingHeight, height, BytesToHexString(hash.ToArrayReverse()))
		return
	}
	// TODO check if already sent Proposal to sender
	if hash.CompareTo(current.GetConfirmingHash()) != 0 {
		log.Warn("requested block doesn't match with local block in process")
		return
	}
	if !current.HasSelfState(hash, voting.ProposalSent) {
		log.Warn("receive invalid request for hash: ", BytesToHexString(hash.ToArrayReverse()))
		return
	}
	nodeID := publickKeyToNodeID(sender)
	// returns if receive duplicate request
	if current.HasNeighborState(nodeID, hash, voting.RequestReceived) {
		log.Warn("duplicate request received for hash: ", BytesToHexString(hash.ToArrayReverse()))
		return
	}
	// set state for request
	current.SetNeighborState(nodeID, hash, voting.RequestReceived)
	content, err := current.GetVotingContent(hash, height)
	if err != nil {
		return
	}
	// generate response message
	responseMsg := NewResponse(&hash, height, current.VotingType(), content)
	// get node which should receive response message
	nodes := ps.GetReceiverNode([]uint64{publickKeyToNodeID(sender)})
	// send response message
	ps.SendConsensusMsg(responseMsg, nodes)
}

func (ps *ProposerService) Initialize(vType voting.VotingContentType) {
	// initial total voting weight
	for _, v := range ps.voting {
		v.GetVotingPool().Reset()
		v.Reset()
	}
}

func (ps *ProposerService) HandleStateProbeMsg(msg *StateProbe, sender *crypto.PubKey) {
	switch msg.ProbeType {
	case BlockHistory:
		switch t := msg.ProbePayload.(type) {
		case *BlockHistoryPayload:
			history := ledger.DefaultLedger.Store.GetBlockHistory(t.StartHeight, t.StartHeight+t.BlockNum)
			s := &StateResponse{history}
			nodes := ps.GetReceiverNode([]uint64{publickKeyToNodeID(sender)})
			ps.SendConsensusMsg(s, nodes)
		}
	}
	return
}

func (ps *ProposerService) HandleResponseMsg(resp *Response, sender *crypto.PubKey) {
	if ps.localNode.GetSyncState() != protocol.PersistFinished {
		return
	}
	votingType := resp.contentType
	current := ps.CurrentVoting(votingType)
	votingHeight := current.GetVotingHeight()
	hash := resp.hash
	height := resp.height
	if height != votingHeight {
		log.Warnf("receive invalid response, consensus height: %d, response height: %d,"+
			" hash: %s\n", votingHeight, height, BytesToHexString(hash.ToArrayReverse()))
		return
	}
	nodeID := publickKeyToNodeID(sender)
	// returns if no request sent before
	if !current.HasNeighborState(nodeID, *hash, voting.RequestSent) {
		log.Warn("consensus state error in Response message handler")
		return
	}
	// returns if receive duplicate response
	if current.HasNeighborState(nodeID, *hash, voting.ProposalReceived) {
		log.Warn("duplicate response received for hash: ", BytesToHexString(hash.ToArrayReverse()))
		return
	}
	// TODO check if the sender is requested neighbor node
	err := current.AddToCache(resp.content)
	if err != nil {
		return
	}
	current.SetNeighborState(nodeID, *hash, voting.ProposalReceived)

	currentVotingPool := current.GetVotingPool()
	neighborWeight, _ := ledger.DefaultLedger.Store.GetVotingWeight(Uint160{})
	// Get voting result from voting pool. If votes is not enough then return.
	maybeFinalHash, err := currentVotingPool.AddVoteThenCounting(votingHeight, nodeID, neighborWeight, *hash)
	if err != nil {
		return
	}
	ps.SetOrChangeMind(votingType, votingHeight, maybeFinalHash)
}

func (ps *ProposerService) SetOrChangeMind(votingType voting.VotingContentType,
	votingHeight uint32, maybeFinalHash *Uint256) {
	current := ps.CurrentVoting(votingType)
	currentVotingPool := current.GetVotingPool()
	if mind, ok := currentVotingPool.GetMind(votingHeight); ok {
		// When current mind has been set, if voting result is different with
		// current mind then do mind changing.
		if mind.CompareTo(*maybeFinalHash) != 0 {
			log.Info("when receive proposal mind changed to neighbor mind: ",
				BytesToHexString(maybeFinalHash.ToArrayReverse()))
			history := currentVotingPool.SetMind(votingHeight, *maybeFinalHash)
			// generate mind changing message
			mindChangingMsg := NewMindChanging(maybeFinalHash, votingHeight, votingType)
			// get node which should receive request message
			var nids []uint64
			for n := range history {
				nids = append(nids, n)
			}
			nodes := ps.GetReceiverNode(nids)
			// send mind changing message
			ps.SendConsensusMsg(mindChangingMsg, nodes)
		}
	} else {
		// Set mind if current mind has not been set.
		currentVotingPool.SetMind(votingHeight, *maybeFinalHash)
		log.Info("when receive proposal mind set to neighbor mind: ",
			BytesToHexString(maybeFinalHash.ToArrayReverse()))
	}
}

func (ps *ProposerService) HandleProposalMsg(proposal *Proposal, sender *crypto.PubKey) {
	hash := *proposal.hash
	height := proposal.height
	nodeID := publickKeyToNodeID(sender)

	// handle proposal when block syncing
	if ps.localNode.GetSyncState() != protocol.PersistFinished {
		err := ps.syncCache.AddVoteForBlock(hash, height, nodeID)
		if err != nil {
			return
		}
		vBlock, err := ps.syncCache.GetBlockFromSyncCache(height)
		if err != nil {
			return
		}
		ps.localNode.SetSyncStopHash(vBlock.Block.Header.Hash(), vBlock.Block.Header.Height)
		ps.syncCache.SetConsensusHeight(height)
		return
	}

	current := ps.CurrentVoting(proposal.contentType)
	votingType := current.VotingType()
	votingHeight := current.GetVotingHeight()
	neighbors := ps.localNode.GetNeighborNoder()
	if height < votingHeight {
		log.Warnf("receive invalid proposal, consensus height: %d, proposal height: %d,"+
			" hash: %s\n", votingHeight, height, BytesToHexString(hash.ToArrayReverse()))
		return
	}
	if height > votingHeight {
		neighborHeight, count := current.CacheProposal(height)
		if 2*count > len(neighbors) {
			log.Errorf("state is different with neighbors, "+
				"current voting height: %d, neighbor height: %d (%d/%d), exits.",
				votingHeight, neighborHeight, count, len(neighbors))
			os.Exit(1)
		}
		return
	}
	// TODO check if the sender is neighbor
	if !current.Exist(hash, height) {
		// generate request message
		requestMsg := NewRequest(&hash, height, votingType)
		// get node which should receive request message
		nodes := ps.GetReceiverNode([]uint64{nodeID})
		// send request message
		ps.SendConsensusMsg(requestMsg, nodes)
		current.SetNeighborState(nodeID, hash, voting.RequestSent)
		log.Warnf("doesn't contain hash in local cache, requesting it from neighbor %s\n",
			BytesToHexString(hash.ToArrayReverse()))
		return
	}
	// returns if receive duplicated proposal
	if current.HasNeighborState(nodeID, hash, voting.ProposalReceived) {
		log.Warn("duplicate proposal received for hash: ", BytesToHexString(hash.ToArrayReverse()))
		return
	}
	// set state when receive a proposal from a neighbor
	current.SetNeighborState(nodeID, hash, voting.ProposalReceived)
	currentVotingPool := current.GetVotingPool()
	neighborWeight, _ := ledger.DefaultLedger.Store.GetVotingWeight(Uint160{})
	// Get voting result from voting pool. If votes is not enough then return.
	maybeFinalHash, err := currentVotingPool.AddVoteThenCounting(votingHeight, nodeID, neighborWeight, hash)
	if err != nil {
		return
	}
	ps.SetOrChangeMind(votingType, votingHeight, maybeFinalHash)
}

func (ps *ProposerService) HandleMindChangingMsg(mindChanging *MindChanging, sender *crypto.PubKey) {
	hash := *mindChanging.hash
	height := mindChanging.height
	nodeID := publickKeyToNodeID(sender)

	// handle mind changing when block syncing
	if ps.localNode.GetSyncState() != protocol.PersistFinished {
		ps.syncCache.AddVoteForBlock(hash, height, nodeID)
		return
	}

	current := ps.CurrentVoting(mindChanging.contentType)
	votingType := current.VotingType()
	votingHeight := current.GetVotingHeight()
	if height < votingHeight {
		log.Warnf("receive invalid mind changing, consensus height: %d, mind changing height: %d,"+
			" hash: %s\n", votingHeight, height, BytesToHexString(hash.ToArrayReverse()))
		return
	}
	currentVotingPool := current.GetVotingPool()
	if currentVotingPool.HasReceivedVoteFrom(votingHeight, nodeID) {
		log.Warn("no proposal received before, so mind changing is invalid")
		return
	}
	neighborWeight, _ := ledger.DefaultLedger.Store.GetVotingWeight(Uint160{})
	// recalculate votes
	maybeFinalHash, err := currentVotingPool.AddVoteThenCounting(votingHeight, nodeID, neighborWeight, hash)
	if err != nil {
		return
	}
	if mind, ok := currentVotingPool.GetMind(votingHeight); ok {
		// When current mind has been set, if voting result is different with
		// current mind then do mind changing.
		if mind.CompareTo(*maybeFinalHash) != 0 {
			log.Info("when receive mindchanging mind change to neighbor mind: ",
				BytesToHexString(maybeFinalHash.ToArrayReverse()))
			history := currentVotingPool.SetMind(votingHeight, *maybeFinalHash)
			// generate mind changing message
			mindChangingMsg := NewMindChanging(maybeFinalHash, votingHeight, votingType)
			// get node which should receive request message
			var nids []uint64
			for n := range history {
				nids = append(nids, n)
			}
			nodes := ps.GetReceiverNode(nids)
			// send mind changing message
			ps.SendConsensusMsg(mindChangingMsg, nodes)
		}
	}
}

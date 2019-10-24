package node

import (
	"fmt"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/nknorg/nkn/block"
	"github.com/nknorg/nkn/chain"
	"github.com/nknorg/nkn/common"
	"github.com/nknorg/nkn/node/consequential"
	"github.com/nknorg/nkn/pb"
	"github.com/nknorg/nkn/util/config"
	"github.com/nknorg/nkn/util/log"
)

const (
	concurrentSyncRequestPerNeighbor = 1
	maxSyncBlockHeadersBatchSize     = 1024
	maxSyncBlocksBatchSize           = 32
	syncReplyTimeout                 = 20 * time.Second
	maxSyncWorkerFails               = 3
	syncWorkerStartInterval          = 20 * time.Millisecond
)

// NewGetBlockHeadersMessage creates a GET_BLOCK_HEADERS message
func NewGetBlockHeadersMessage(startHeight, endHeight uint32) (*pb.UnsignedMessage, error) {
	msgBody := &pb.GetBlockHeaders{
		StartHeight: startHeight,
		EndHeight:   endHeight,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pb.UnsignedMessage{
		MessageType: pb.GET_BLOCK_HEADERS,
		Message:     buf,
	}

	return msg, nil
}

// NewGetBlockHeadersReply creates a GET_BLOCK_HEADERS_REPLY message in respond
// to GET_BLOCK_HEADERS message
func NewGetBlockHeadersReply(headers []*block.Header) (*pb.UnsignedMessage, error) {
	pbHeaders := make([]*pb.Header, len(headers))
	for i, header := range headers {
		pbHeaders[i] = header.Header
	}

	msgBody := &pb.GetBlockHeadersReply{
		BlockHeaders: pbHeaders,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pb.UnsignedMessage{
		MessageType: pb.GET_BLOCK_HEADERS_REPLY,
		Message:     buf,
	}

	return msg, nil
}

// NewGetBlocksMessage creates a GET_BLOCKS message
func NewGetBlocksMessage(startHeight, endHeight uint32) (*pb.UnsignedMessage, error) {
	msgBody := &pb.GetBlocks{
		StartHeight: startHeight,
		EndHeight:   endHeight,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pb.UnsignedMessage{
		MessageType: pb.GET_BLOCKS,
		Message:     buf,
	}

	return msg, nil
}

// NewGetBlocksReply creates a GET_BLOCKS_REPLY message in respond to GET_BLOCKS
// message
func NewGetBlocksReply(blocks []*block.Block) (*pb.UnsignedMessage, error) {
	msgBlocks := make([]*pb.Block, len(blocks))
	for i, block := range blocks {
		msgBlocks[i] = block.ToMsgBlock()
	}

	msgBody := &pb.GetBlocksReply{
		Blocks: msgBlocks,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &pb.UnsignedMessage{
		MessageType: pb.GET_BLOCKS_REPLY,
		Message:     buf,
	}

	return msg, nil
}

// getBlockHeadersMessageHandler handles a GET_BLOCK_HEADERS message
func (localNode *LocalNode) getBlockHeadersMessageHandler(remoteMessage *RemoteMessage) ([]byte, bool, error) {
	replyMsg, err := NewGetBlockHeadersReply(nil)
	if err != nil {
		return nil, false, err
	}

	replyBuf, err := localNode.SerializeMessage(replyMsg, false)
	if err != nil {
		return nil, false, err
	}

	msgBody := &pb.GetBlockHeaders{}
	err = proto.Unmarshal(remoteMessage.Message, msgBody)
	if err != nil {
		return replyBuf, false, err
	}

	startHeight, endHeight := msgBody.StartHeight, msgBody.EndHeight

	if endHeight < startHeight {
		return replyBuf, false, nil
	}

	if endHeight-startHeight > maxSyncBlockHeadersBatchSize {
		return replyBuf, false, nil
	}

	if endHeight > chain.DefaultLedger.Store.GetHeaderHeight() {
		return replyBuf, false, nil
	}

	headers := make([]*block.Header, endHeight-startHeight+1, endHeight-startHeight+1)
	for height := startHeight; height <= endHeight; height++ {
		headers[height-startHeight], err = chain.DefaultLedger.Store.GetHeaderByHeight(height)
		if err != nil {
			return replyBuf, false, err
		}
	}

	replyMsg, err = NewGetBlockHeadersReply(headers)
	if err != nil {
		return replyBuf, false, err
	}

	replyBuf, err = localNode.SerializeMessage(replyMsg, false)
	return replyBuf, false, err
}

// getBlocksMessageHandler handles a GET_BLOCKS message
func (localNode *LocalNode) getBlocksMessageHandler(remoteMessage *RemoteMessage) ([]byte, bool, error) {
	replyMsg, err := NewGetBlocksReply(nil)
	if err != nil {
		return nil, false, err
	}

	replyBuf, err := localNode.SerializeMessage(replyMsg, false)
	if err != nil {
		return nil, false, err
	}

	msgBody := &pb.GetBlocks{}
	err = proto.Unmarshal(remoteMessage.Message, msgBody)
	if err != nil {
		return replyBuf, false, err
	}

	startHeight, endHeight := msgBody.StartHeight, msgBody.EndHeight

	if endHeight < startHeight {
		return replyBuf, false, nil
	}

	if endHeight-startHeight > maxSyncBlocksBatchSize {
		return replyBuf, false, nil
	}

	if endHeight > chain.DefaultLedger.Store.GetHeight() {
		return replyBuf, false, nil
	}

	blocks := make([]*block.Block, endHeight-startHeight+1, endHeight-startHeight+1)
	for height := startHeight; height <= endHeight; height++ {
		blocks[height-startHeight], err = chain.DefaultLedger.Store.GetBlockByHeight(height)
		if err != nil {
			return replyBuf, false, err
		}
	}

	replyMsg, err = NewGetBlocksReply(blocks)
	if err != nil {
		return replyBuf, false, err
	}

	replyBuf, err = localNode.SerializeMessage(replyMsg, false)
	return replyBuf, false, err
}

// GetBlockHeaders requests a range of consecutive block headers from a neighbor
// using GET_BLOCK_HEADERS message
func (remoteNode *RemoteNode) GetBlockHeaders(startHeight, endHeight uint32) ([]*block.Header, error) {
	if startHeight > endHeight {
		return nil, fmt.Errorf("start height %d is higher than end height %d", startHeight, endHeight)
	}

	msg, err := NewGetBlockHeadersMessage(startHeight, endHeight)
	if err != nil {
		return nil, err
	}

	buf, err := remoteNode.localNode.SerializeMessage(msg, false)
	if err != nil {
		return nil, err
	}

	replyBytes, err := remoteNode.SendBytesSyncWithTimeout(buf, syncReplyTimeout)
	if err != nil {
		return nil, err
	}

	replyMsg := &pb.GetBlockHeadersReply{}
	err = proto.Unmarshal(replyBytes, replyMsg)
	if err != nil {
		return nil, err
	}

	if len(replyMsg.BlockHeaders) < int(endHeight-startHeight+1) {
		return nil, fmt.Errorf("result contains %d instead of %d headers", len(replyMsg.BlockHeaders), endHeight-startHeight+1)
	}

	headers := make([]*block.Header, len(replyMsg.BlockHeaders))
	for i, msgHeader := range replyMsg.BlockHeaders {
		headers[i] = &block.Header{Header: msgHeader}
	}

	return headers, nil
}

// GetBlocks requests a range of consecutive blocks from a neighbor using
// GET_BLOCKS message
func (remoteNode *RemoteNode) GetBlocks(startHeight, endHeight uint32) ([]*block.Block, error) {
	if startHeight > endHeight {
		return nil, fmt.Errorf("start height %d is higher than end height %d", startHeight, endHeight)
	}

	msg, err := NewGetBlocksMessage(startHeight, endHeight)
	if err != nil {
		return nil, err
	}

	buf, err := remoteNode.localNode.SerializeMessage(msg, false)
	if err != nil {
		return nil, err
	}

	replyBytes, err := remoteNode.SendBytesSyncWithTimeout(buf, syncReplyTimeout)
	if err != nil {
		return nil, err
	}

	replyMsg := &pb.GetBlocksReply{}
	err = proto.Unmarshal(replyBytes, replyMsg)
	if err != nil {
		return nil, err
	}

	if len(replyMsg.Blocks) < int(endHeight-startHeight+1) {
		return nil, fmt.Errorf("result contains %d instead of %d blocks", len(replyMsg.Blocks), endHeight-startHeight+1)
	}

	blocks := make([]*block.Block, len(replyMsg.Blocks))
	for i, msgBlock := range replyMsg.Blocks {
		blocks[i] = &block.Block{}
		blocks[i].FromMsgBlock(msgBlock)
	}

	return blocks, nil
}

// initSyncing initializes block syncing state and registers message handler
func (localNode *LocalNode) initSyncing() {
	localNode.AddMessageHandler(pb.GET_BLOCK_HEADERS, localNode.getBlockHeadersMessageHandler)
	localNode.AddMessageHandler(pb.GET_BLOCKS, localNode.getBlocksMessageHandler)
	localNode.ResetSyncing()
}

// StartSyncing starts block syncing from current local ledger until it gets to
// block height stopHeight with block hash stopHash from given neighbors
func (localNode *LocalNode) StartSyncing(stopHash common.Uint256, stopHeight uint32, neighbors []*RemoteNode) (bool, error) {
	var err error
	started := false

	localNode.RLock()
	syncOnce := localNode.syncOnce
	localNode.RUnlock()

	syncOnce.Do(func() {
		started = true
		localNode.SetSyncState(pb.SYNC_STARTED)

		currentHeight := chain.DefaultLedger.Store.GetHeight()
		currentHash := chain.DefaultLedger.Store.GetHeaderHashByHeight(currentHeight)
		if stopHeight <= currentHeight {
			err = fmt.Errorf("sync stop height %d is not higher than current height %d", stopHeight, currentHeight)
			return
		}

		if len(neighbors) == 0 {
			err = fmt.Errorf("no neighbors to sync from")
			return
		}

		var rollbacked bool
		rollbacked, err = localNode.maybeRollback(neighbors)
		if err != nil {
			panic(fmt.Errorf("Rollback error: %v", err))
		}
		if rollbacked {
			currentHeight = chain.DefaultLedger.Store.GetHeight()
			currentHash = chain.DefaultLedger.Store.GetHeaderHashByHeight(currentHeight)
		}

		startTime := time.Now()
		var headersHash []common.Uint256
		headersHash, err = localNode.syncBlockHeaders(currentHeight+1, stopHeight, currentHash, stopHash, neighbors)
		if err != nil {
			err = fmt.Errorf("sync block headers error: %v", err)
			return
		}

		log.Infof("Synced %d block headers in %s", stopHeight-currentHeight, time.Since(startTime))

		startTime = time.Now()
		err = localNode.syncBlocks(currentHeight+1, stopHeight, neighbors, headersHash)
		if err != nil {
			err = fmt.Errorf("sync blocks error: %v", err)
			return
		}

		log.Infof("Synced %d blocks in %s", stopHeight-currentHeight, time.Since(startTime))

		localNode.SetSyncState(pb.SYNC_FINISHED)
	})

	return started, err
}

// ResetSyncing resets syncOnce and allows for future block syncing
func (localNode *LocalNode) ResetSyncing() {
	localNode.Lock()
	defer localNode.Unlock()
	localNode.syncOnce = new(sync.Once)
}

func (localNode *LocalNode) syncBlockHeaders(startHeight, stopHeight uint32, startPrevHash, stopHash common.Uint256, neighbors []*RemoteNode) ([]common.Uint256, error) {
	var nextHeader *block.Header
	headersHash := make([]common.Uint256, stopHeight-startHeight+1, stopHeight-startHeight+1)
	numBatches := (stopHeight-startHeight)/config.Parameters.SyncBlockHeadersBatchSize + 1
	numWorkers := uint32(len(neighbors)) * concurrentSyncRequestPerNeighbor

	getBatchHeightRange := func(batchID uint32) (uint32, uint32) {
		batchStartHeight := startHeight + batchID*config.Parameters.SyncBlockHeadersBatchSize
		batchEndHeight := batchStartHeight + config.Parameters.SyncBlockHeadersBatchSize - 1
		if batchEndHeight > stopHeight {
			batchEndHeight = stopHeight
		}
		return batchStartHeight, batchEndHeight
	}

	getHeader := func(workerID, batchID uint32) (interface{}, bool) {
		neighbor := neighbors[workerID%uint32(len(neighbors))]
		batchStartHeight, batchEndHeight := getBatchHeightRange(batchID)

		batchHeaders, err := neighbor.GetBlockHeaders(batchStartHeight, batchEndHeight)
		if err != nil {
			log.Warningf("Get block headers error: %v", err)
			return nil, false
		}

		return batchHeaders, true
	}

	saveHeader := func(batchID uint32, result interface{}) bool {
		batchHeaders, ok := result.([]*block.Header)
		if !ok {
			log.Warningf("Convert batch headers error")
			return false
		}

		batchStartHeight, batchEndHeight := getBatchHeightRange(batchID)
		for height := batchEndHeight; height >= batchStartHeight; height-- {
			header := batchHeaders[height-batchStartHeight]
			headerHash := header.Hash()
			if height == stopHeight && headerHash != stopHash {
				log.Warningf("End header hash %s is different from stop hash %s", headerHash.ToHexString(), stopHash.ToHexString())
				return false
			}
			if height < stopHeight {
				nextPrevHash, _ := common.Uint256ParseFromBytes(nextHeader.UnsignedHeader.PrevBlockHash)
				if nextHeader == nil || headerHash != nextPrevHash {
					log.Warningf("Header hash %s is different from prev hash in next block %s", headerHash.ToHexString(), nextPrevHash.ToHexString())
					return false
				}
			}

			prevHash, _ := common.Uint256ParseFromBytes(header.UnsignedHeader.PrevBlockHash)
			if height == startHeight && prevHash != startPrevHash {
				log.Warningf("Start header prev hash %s is different from start prev hash %s", prevHash.ToHexString(), startPrevHash.ToHexString())
				return false
			}
			headersHash[height-startHeight] = headerHash
			nextHeader = header
		}
		return true
	}

	cs, err := consequential.NewConSequential(&consequential.Config{
		StartJobID:          numBatches - 1,
		EndJobID:            0,
		JobBufSize:          config.Parameters.SyncBatchWindowSize,
		WorkerPoolSize:      numWorkers,
		MaxWorkerFails:      maxSyncWorkerFails,
		WorkerStartInterval: syncWorkerStartInterval,
		RunJob:              getHeader,
		FinishJob:           saveHeader,
	})
	if err != nil {
		return nil, err
	}

	err = cs.Start()
	if err != nil {
		return nil, err
	}

	return headersHash, nil
}

func (localNode *LocalNode) syncBlocks(startHeight, stopHeight uint32, neighbors []*RemoteNode, headersHash []common.Uint256) error {
	numBatches := (stopHeight-startHeight)/config.Parameters.SyncBlocksBatchSize + 1
	numWorkers := uint32(len(neighbors)) * concurrentSyncRequestPerNeighbor

	getBatchHeightRange := func(batchID uint32) (uint32, uint32) {
		batchStartHeight := startHeight + batchID*config.Parameters.SyncBlocksBatchSize
		batchEndHeight := batchStartHeight + config.Parameters.SyncBlocksBatchSize - 1
		if batchEndHeight > stopHeight {
			batchEndHeight = stopHeight
		}
		return batchStartHeight, batchEndHeight
	}

	getBlock := func(workerID, batchID uint32) (interface{}, bool) {
		neighbor := neighbors[workerID%uint32(len(neighbors))]
		batchStartHeight, batchEndHeight := getBatchHeightRange(batchID)

		batchBlocks, err := neighbor.GetBlocks(batchStartHeight, batchEndHeight)
		if err != nil {
			log.Warningf("Get blocks error: %v", err)
			return nil, false
		}

		return batchBlocks, true
	}

	saveBlock := func(batchID uint32, result interface{}) bool {
		batchBlocks, ok := result.([]*block.Block)
		if !ok {
			log.Warningf("Convert batch blocks error")
			return false
		}

		batchStartHeight, batchEndHeight := getBatchHeightRange(batchID)
		for height := batchStartHeight; height <= batchEndHeight; height++ {
			block := batchBlocks[height-batchStartHeight]
			blockHash := block.Hash()
			headerHash := headersHash[height-startHeight]
			if blockHash != headerHash {
				log.Warningf("Block hash %s is different from header hash %s", (&blockHash).ToHexString(), (&headerHash).ToHexString())
				return false
			}

			err := chain.DefaultLedger.Blockchain.AddBlock(block, false)
			if err != nil {
				return false
			}
		}

		return true
	}

	cs, err := consequential.NewConSequential(&consequential.Config{
		StartJobID:          0,
		EndJobID:            numBatches - 1,
		JobBufSize:          config.Parameters.SyncBatchWindowSize,
		WorkerPoolSize:      numWorkers,
		MaxWorkerFails:      maxSyncWorkerFails,
		WorkerStartInterval: syncWorkerStartInterval,
		RunJob:              getBlock,
		FinishJob:           saveBlock,
	})
	if err != nil {
		return err
	}

	err = cs.Start()
	if err != nil {
		return err
	}

	return nil
}

package chain_domain

type ChainStorageRepo interface {
	GetLastBlock() *Block
	GetBlockByHash(hash string) *Block
	GetBlockByNumber(number uint64) *Block
	SaveBlock(block *Block) *Block
}
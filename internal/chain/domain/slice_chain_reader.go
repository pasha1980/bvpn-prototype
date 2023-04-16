package domain

import "bvpn-prototype/internal/protocol/entity"

type SliceChainReader struct {
	index int64
	slice []entity.Block
}

func (r *SliceChainReader) Start() {
	r.index = 0
}

func (r *SliceChainReader) Next() *entity.Block {
	r.index++
	if r.index > (int64(len(r.slice)) - 1) {
		return nil
	}

	block := r.slice[r.index]
	return &block
}

func (r *SliceChainReader) Last() *entity.Block {
	if len(r.slice) == 0 {
		return nil
	}

	block := r.slice[len(r.slice)-1]
	return &block
}

func (r *SliceChainReader) Len() int64 {
	return int64(len(r.slice))
}

func NewSliceChainReader(chain []entity.Block) *SliceChainReader {
	return &SliceChainReader{
		index: 0,
		slice: chain,
	}
}

// Copyright 2017 Delos Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0
package types

type BlockMeta struct {
	Hash        []byte        `json:"hash"`         // The block hash
	Header      *Header       `json:"header"`       // The block's Header
	PartsHeader PartSetHeader `json:"parts_header"` // The PartSetHeader, for transfer
}

func NewBlockMeta(block *Block, blockParts *PartSet) *BlockMeta {
	return &BlockMeta{
		Hash:        block.Hash(),
		Header:      block.Header,
		PartsHeader: blockParts.Header(),
	}
}

package rs

const (
	DataShards 		= 4
	ParityShards 	= 2
	AllShares		= DataShards + ParityShards
	BlockPerShare	= 8000
	BlockSize		= BlockPerShare * DataShards
)

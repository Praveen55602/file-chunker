package manifest

type Manifest struct {
	FileID      string
	ChunkSize   int64
	Filename    string
	TotalSize   int64
	ChunkHashes []string
}

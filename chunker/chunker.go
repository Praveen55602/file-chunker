package chunker

import "fmt"

type Manifest struct {
	Filename    string
	TotalSize   int64
	ChunkHashes []string
}

// Split contains the actual logic to chunk the file
func Split(filepath string) (*Manifest, error) {
	fmt.Printf("Splitting file: %s", filepath)
	// We will write the real logic here later.
	// For now, let's just return a dummy Manifest to prove the package works.
	return &Manifest{
		Filename:  filepath,
		TotalSize: 1024,
	}, nil
}

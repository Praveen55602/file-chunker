package assembler

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/praveensingh/file-chunker/manifest"
)

var (
	ErrHashMismatch = errors.New("chunk hash does not match expected value")
	ErrUnknownChunk = errors.New("hash not found in manifest")
)

// Assembler handles out-of-order, on-the-fly file reconstruction
type Assembler struct {
	manifest *manifest.Manifest
	file     *os.File
}

// NewAssembler initializes the target file and pre-allocates the disk space
func NewAssembler(m *manifest.Manifest, targetPath string) (*Assembler, error) {
	// Open the file for read/write, create it if it doesn't exist
	file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open target file: %w", err)
	}

	// Pre-allocate the exact file size on disk instantly
	if err := file.Truncate(m.TotalSize); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to pre-allocate file space: %w", err)
	}

	return &Assembler{
		m,
		file,
	}, nil
}

// Close ensures the file handle is released when the download finishes
func (a *Assembler) Close() error {
	return a.file.Close()
}

// WriteChunk accepts raw bytes, verifies the hash, and writes it to the correct disk offset.
// This is completely decoupled from how the data was downloaded.
func (a *Assembler) WriteChunk(expectedHash string, data []byte) error {
	// 1. Instant Verification (The Bouncer)
	actualHashBytes := sha256.Sum256(data)
	actualHash := hex.EncodeToString(actualHashBytes[:])

	if actualHash != expectedHash {
		return ErrHashMismatch
	}

	// 2. Find the chunk's index in the Manifest to calculate its offset
	chunkIndex := -1
	for i, hash := range a.manifest.ChunkHashes {
		if hash == expectedHash {
			chunkIndex = i
			break
		}
	}

	if chunkIndex == -1 {
		return ErrUnknownChunk
	}

	// 3. Calculate the byte offset
	// e.g., If ChunkSize is 512KB, Chunk 0 goes to offset 0, Chunk 1 goes to 524,288.
	offset := int64(chunkIndex) * a.manifest.ChunkSize

	// 4. Write exactly to that position in the pre-allocated file
	_, err := a.file.WriteAt(data, offset)
	if err != nil {
		return fmt.Errorf("failed to write chunk at offset %d: %w", offset, err)
	}

	return nil
}

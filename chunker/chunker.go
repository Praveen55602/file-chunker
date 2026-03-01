package chunker

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Praveensingh55602/file-chunker/manifest"
)

// Split reads a file, chunks it into fixed sizes, saves them to an output directory,
// and returns the Manifest required to reassemble them.
func Split(sourcePath string, outputDir string, chunkSize int64) (*manifest.Manifest, error) {
	// open the file for reading(does not load the file into ram only gives a pointer to file location from there we can read it)
	file, err := os.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	// get the file info
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	manifest := &manifest.Manifest{
		FileID:    info.Name(), // Keeping it simple for the MVP
		Filename:  info.Name(),
		TotalSize: info.Size(),
		ChunkSize: chunkSize,
	}

	// Ensure the output directory exists
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// create our fixed-size memory buffer
	buffer := make([]byte, chunkSize)

	// the streaming loop
	for {
		// Read up to chunkSize bytes from the file into our buffer
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		// we've reached end of file
		if bytesRead == 0 {
			break
		}

		// CRITICAL: Slice the buffer to exactly the number of bytes read.
		// The last chunk of a file is almost never exactly 'chunkSize' bytes.
		chunkData := buffer[:bytesRead]

		// Hash the chunks
		hash := sha256.Sum256(chunkData)
		hashedString := hex.EncodeToString(hash[:])

		// Add the hash to our manifest in the correct order
		manifest.ChunkHashes = append(manifest.ChunkHashes, hashedString)

		// 7. Write the chunk to disk.
		// Naming the chunk file by its hash is a pattern called "Content-Addressable Storage"
		chunkPath := filepath.Join(outputDir, hashedString)
		if err := os.WriteFile(chunkPath, chunkData, 0755); err != nil {
			return nil, fmt.Errorf("failed to write chunk %s: %w", hashedString, err)
		}
	}
	return manifest, nil
}

package test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/Praveen55602/file-chunker/assembler"
	"github.com/Praveen55602/file-chunker/manifest"
)

func TestAssembler_WriteChunk(t *testing.T) {
	manifest, err := manifest.LoadManifest("test-chunks/testimage.jpg-manifest.json")
	if err != nil {
		t.Skip("chunks and manifest not found please run the split test first")
	}
	testAssembler, err := assembler.NewAssembler(manifest, "./assembleresult")
	if err != nil {
		t.Fatal(err)
	}
	defer testAssembler.Close()

	// Loop starting from the LAST chunk down to the FIRST chunk
	for i := len(manifest.ChunkHashes) - 1; i >= 0; i-- {
		hash := manifest.ChunkHashes[i]

		// Simulate receiving the payload by reading the chunk file from disk
		chunkPath := filepath.Join("./test-chunks", hash)
		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			log.Fatalf("Failed to read chunk from disk: %v", err)
		}

		fmt.Printf("-> Writing Chunk %d (Hash: %s...) to offset %d\n", i, hash[:8], int64(i)*manifest.ChunkSize)

		// The Assembler doesn't care what order we call this in!
		err = testAssembler.WriteChunk(hash, chunkData)
		if err != nil {
			t.Errorf("Assembler rejected chunk: %v", err)
		}
	}

	if equal, err := simpleCompare("testimage.jpg", "assembleresult.jpg"); err != nil || !equal {
		t.Errorf("file not successfullly assembled")
	}
}

func simpleCompare(file1, file2 string) (bool, error) {
	f1, err := os.ReadFile(file1)
	if err != nil {
		return false, err
	}

	f2, err := os.ReadFile(file2)
	if err != nil {
		return false, err
	}

	return bytes.Equal(f1, f2), nil
}

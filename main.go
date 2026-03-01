// main.go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/praveen/file-chunker/assembler"
	"github.com/praveen/file-chunker/chunker"
)

func main() {
	// 1. Create a test file with enough text to generate multiple chunks
	sourceFile := "test.txt"
	originalText := "Line 1: Hello P2P World!\nLine 2: This is a test of out-of-order chunking.\nLine 3: If this works, the file offset math is perfect.\nLine 4: End of transmission."
	os.WriteFile(sourceFile, []byte(originalText), 0644)

	var chunkSize int64 = 15 // Tiny chunk size to force ~10 chunks

	// 2. SPLIT the file
	fmt.Println("Splitting file...")
	manifest, err := chunker.Split(sourceFile, "./chunks", chunkSize)
	if err != nil {
		log.Fatalf("Split failed: %v", err)
	}
	fmt.Printf("Created %d chunks.\n\n", len(manifest.ChunkHashes))

	// 3. INIT the Assembler
	targetFile := "restored.txt"
	fmt.Printf("Initializing Assembler. Pre-allocating %d bytes for %s...\n", manifest.TotalSize, targetFile)
	assembler, err := assembler.NewAssembler(manifest, targetFile)
	if err != nil {
		log.Fatalf("Failed to init assembler: %v", err)
	}

	// 4. ASSEMBLE OUT OF ORDER (Backwards)
	fmt.Println("\nSimulating out-of-order P2P downloads (writing backwards)...")

	// Loop starting from the LAST chunk down to the FIRST chunk
	for i := len(manifest.ChunkHashes) - 1; i >= 0; i-- {
		hash := manifest.ChunkHashes[i]

		// Simulate receiving the payload by reading the chunk file from disk
		chunkPath := filepath.Join("./chunks", hash)
		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			log.Fatalf("Failed to read chunk from disk: %v", err)
		}

		fmt.Printf("-> Writing Chunk %d (Hash: %s...) to offset %d\n", i, hash[:8], int64(i)*chunkSize)

		// The Assembler doesn't care what order we call this in!
		err = assembler.WriteChunk(hash, chunkData)
		if err != nil {
			log.Fatalf("Assembler rejected chunk: %v", err)
		}
	}

	// Close the file handle so the OS flushes the data and lets us read it
	assembler.Close()

	// 5. VERIFY the final file
	restoredData, err := os.ReadFile(targetFile)
	if err != nil {
		log.Fatalf("Failed to read restored file: %v", err)
	}

	fmt.Println("\n--- RESTORED FILE CONTENTS ---")
	fmt.Println(string(restoredData))
	fmt.Println("------------------------------")

	if string(restoredData) == originalText {
		fmt.Println("\nSUCCESS: The out-of-order assembly worked perfectly!")
	} else {
		fmt.Println("\nFAILED: The restored file does not match the original.")
	}
}

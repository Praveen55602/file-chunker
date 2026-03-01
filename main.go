// main.go
package main

import (
	"fmt"
	"log"

	// Here is where you import your custom package!
	// It is the module name + the folder path
	"github.com/praveen/file-chunker/chunker"
)

func main() {
	fmt.Println("Starting the chunking system...")

	// We'll use a tiny chunk size of 15 bytes to force multiple chunks for testing
	var chunkSize int64 = 15

	manifest, err := chunker.Split("testfile.txt", "./chunks", chunkSize)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully chunked: %s\n", manifest.Filename)
	fmt.Printf("Total Size: %d bytes\n", manifest.TotalSize)
	fmt.Printf("Number of chunks created: %d\n", len(manifest.ChunkHashes))

	fmt.Println("Chunk Hashes in order:")
	for i, hash := range manifest.ChunkHashes {
		fmt.Printf(" %d: %s\n", i, hash)
	}
}

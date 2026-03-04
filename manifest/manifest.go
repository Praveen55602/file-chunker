package manifest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Manifest struct {
	FileID        string
	FileExtention string
	ChunkSize     int64
	Filename      string
	TotalSize     int64
	ChunkHashes   []string
}

func LoadManifest(manifestPath string) (*Manifest, error) {
	//manifest path will have a json file, which we'll deserialize
	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open manifest file: %w", err)
	}
	defer manifestFile.Close()

	data, err := io.ReadAll(manifestFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file: %w", err)
	}

	//deserialize
	manifest := &Manifest{}
	err = json.Unmarshal(data, manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize manifest: %w", err)
	}

	return manifest, nil
}

func SaveManifest(manifest *Manifest, manifestDirectory string) error {
	log.Printf("Saving manifest to %s\n", manifestDirectory)
	//serialize
	data, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to serialize manifest: %w", err)
	}

	fileNameWithoutExtention := manifest.Filename[:len(manifest.Filename)-len(manifest.FileExtention)]
	manifestPath := filepath.Join(manifestDirectory, fileNameWithoutExtention+"-manifest.json")
	if err := os.WriteFile(manifestPath, data, 0755); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}
	log.Printf("manifest saved to %s\n", manifestDirectory)
	return nil
}

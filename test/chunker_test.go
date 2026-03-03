package test

import (
	"os"
	"reflect"
	"testing"

	"github.com/Praveen55602/file-chunker/chunker"
	"github.com/Praveen55602/file-chunker/manifest"
)

func TestLoadManifest(t *testing.T) {
	_, err := os.Stat("./test-chunks")
	if os.IsNotExist(err) {
		t.Skip("chunks not found please run the split test first")
	}
	type args struct {
		manifestPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *manifest.Manifest
		wantErr bool
	}{
		{
			name: "load manifest success test",
			args: args{
				manifestPath: "./test-chunks/testimage.jpg-manifest.json",
			},
			want: &manifest.Manifest{
				FileID:        "testimage.jpg",
				FileExtention: ".jpg",
				ChunkSize:     2048,
				Filename:      "testimage.jpg",
				TotalSize:     16130,
				ChunkHashes:   []string{"3102c8562326d65d632ad73ebd4c5fa743e6c480855ff99cf7fbfd2ef2fa7ca3", "f77a153ac9cfcf872343a03f2573d4edbfaf5f520c9e02fc92e937b76822ad51", "5dc5d60971fb9209b7364466fc5a5a8e2b617814acc119c0b32257c08ea118f5", "479fe05bf73264a2c26825a449ec79cbc00ed8757124ae2d11bb6f88a77885c5", "d1d5ef48838473da1f90e26435ba5d825290bc52c4991fff5007f4ccb9958cb9", "67274b09a16bf5b5deb5a066c9ac2c040a310122f597cc90562182bcb4db4e59", "2d23dc9f5f3456cd5e5beeff6d879153fb84e4342f76b8c685ffa13c886646ac", "6c8f04260b682fbfbfd022b3b4945a1dfe0f3ed903d313813ee2a300e44a7bc2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := manifest.LoadManifest(tt.args.manifestPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadManifest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	type args struct {
		sourcePath string
		outputDir  string
		chunkSize  int64
	}
	tests := []struct {
		name    string
		args    args
		want    *manifest.Manifest
		wantErr bool
	}{
		{
			name: "split success test",
			args: args{
				sourcePath: "testimage.jpg",
				outputDir:  "./test-chunks",
				chunkSize:  1024 * 2,
			},
			want: &manifest.Manifest{
				FileID:        "testimage.jpg",
				FileExtention: ".jpg",
				ChunkSize:     2048,
				Filename:      "testimage.jpg",
				TotalSize:     16130,
				ChunkHashes:   []string{"3102c8562326d65d632ad73ebd4c5fa743e6c480855ff99cf7fbfd2ef2fa7ca3", "f77a153ac9cfcf872343a03f2573d4edbfaf5f520c9e02fc92e937b76822ad51", "5dc5d60971fb9209b7364466fc5a5a8e2b617814acc119c0b32257c08ea118f5", "479fe05bf73264a2c26825a449ec79cbc00ed8757124ae2d11bb6f88a77885c5", "d1d5ef48838473da1f90e26435ba5d825290bc52c4991fff5007f4ccb9958cb9", "67274b09a16bf5b5deb5a066c9ac2c040a310122f597cc90562182bcb4db4e59", "2d23dc9f5f3456cd5e5beeff6d879153fb84e4342f76b8c685ffa13c886646ac", "6c8f04260b682fbfbfd022b3b4945a1dfe0f3ed903d313813ee2a300e44a7bc2"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chunker.Split(tt.args.sourcePath, tt.args.outputDir, tt.args.chunkSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("Split() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split() = %v, want %v", got, tt.want)
			}
		})
	}
}

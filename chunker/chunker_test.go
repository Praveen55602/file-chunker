package chunker

import (
	"reflect"
	"testing"

	"github.com/Praveen55602/file-chunker/manifest"
)

func TestLoadManifest(t *testing.T) {
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
				manifestPath: "../chunks/test.txt-manifest.json",
			},
			want: &manifest.Manifest{
				FileID:    "test.txt",
				ChunkSize: 15,
				Filename:  "test.txt",
				TotalSize: 158, ChunkHashes: []string{
					"da65daa484de96cd88a88eee1bc8dec35b9038caf9624597aee38ad65de7c581",
					"dc26b5a09056f02b81a1c85c2bed92435904f4a951937db982da2f15e9d4f295",
					"3f7483716edfff5ad0d837d72f105d7600b492c4355ab7fdd1223c01437c80f5",
					"fcb8851d97b1d88bbd1e4e5631d795de39881b9c13bd4d51db35dae14219ce27",
					"b5781e8a119b41c1b9ff5b76df3d647304e62150970da5d1aed6cbc1304fad08",
					"19a3569d2e5e2a9ffaf234f0149a2383cc1c8070199e062774b91e26917a4307",
					"9c1ca7d441e35cf9a38af8bddd968a707544508773f0e205bac7e42ac1ea71e2",
					"d1ad85732c0beed39f20b021d20b2011866a56d421c4d2c458f17f27bc4b1997",
					"723f5d1377a4147be143be3b101cf2962020d5be6cf5de933c14bcb51bb3dd32",
					"8456fabadd4d480c30d7b16260de73cec4a8abbb738cb31fb6e2995d70912137",
					"8edc42bed1c50247f58ad07bfc3adae09949e429d68250933280966fd6393cef",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadManifest(tt.args.manifestPath)
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

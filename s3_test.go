package blobcat_test

import (
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/midorigreen/blobcat"
	"github.com/pkg/errors"
)

func TestWriteExt(t *testing.T) {
	tests := []struct {
		name      string
		ext       string
		input     string
		want      string
		wantError bool
		err       error
	}{
		{
			name:  "normal scenario",
			ext:   "gz",
			input: "hogehogehoge",
			want:  "hogehogehoge",
		},
		{
			name:      "unexpected scenario",
			ext:       "txt",
			input:     "hogehogehoge",
			wantError: true,
			err:       errors.New("not implements ext"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inb := &bytes.Buffer{}
			gw := gzip.NewWriter(inb)
			_, gerr := gw.Write([]byte(tt.input))
			if gerr != nil {
				t.Fatalf("failed gzip write: %v", gerr)
			}
			gw.Flush()
			gw.Close()

			got := &bytes.Buffer{}
			err := blobcat.WriteExt(tt.ext, inb, got)
			if !tt.wantError && err != nil {
				t.Fatalf("failed exec test func: %v", err)
			}

			if tt.wantError && err.Error() != tt.err.Error() {
				t.Fatalf("want: %#v, got: %#v", tt.err, err)
			}

			if tt.want != got.String() {
				t.Fatalf("unexpected value want: %v, got: %v", tt.want, got.String())
			}
		})
	}
}

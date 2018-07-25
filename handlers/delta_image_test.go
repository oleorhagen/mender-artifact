// Copyright 2017 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package handlers

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeltaCompose(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	tw := tar.NewWriter(buf)
	defer tw.Close()
	f, err := ioutil.TempFile("", "update")
	assert.NoError(t, err)
	defer os.Remove(f.Name())
	// Test the composed delta header
	tests := map[string]struct {
		args     *ComposeHeaderArgs
		expected string
	}{
		"tc1": {
			args: &ComposeHeaderArgs{
				TarWriter: tw,
				No:        3,
			},
			expected: "",
		},
	}

	for name, test := range tests {
		delta := NewDelta(f.Name(), "foo", "bar")
		err := delta.ComposeHeader(test.args)
		require.Nil(t, err, "%s failed", name)
		// Iterate over the files in the archive.
		tr := tar.NewReader(buf)
		fileHdr, err := tr.Next()
		require.Nil(t, err)
		if filepath.Base(fileHdr.Name) != "files" {
			t.Fatalf("First file header is not files: FileHeader: %s", fileHdr.Name)
		}
		typeInfoHdr, err := tr.Next()
		require.Nil(t, err)
		if filepath.Base(typeInfoHdr.Name) != "type-info" {
			t.Fatal("Second file header is not type-info: FileHeader: %s", typeInfoHdr.Name)
		}
		metaData, err := tr.Next()
		require.Nil(t, err)
		if filepath.Base(metaData.Name) != "meta-data" {
			t.Fatal("Third file header is not meta-data: FileHeader: %s", metaData.Name)
		}
		eof, err := tr.Next()
		require.Equal(t, io.EOF, err, "TarReader reads more than the three composed header files, FileHeader: %v", eof)
	}
}

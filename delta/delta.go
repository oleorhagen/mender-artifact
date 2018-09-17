package delta

import (
	"bytes"
	"errors"
	"io"

	"github.com/mendersoftware/mender-artifact/areader"
	"github.com/mendersoftware/mender-artifact/handlers"
)

type delta bytes.Buffer

func newDelta() delta {
	var d delta
	return d
}

// Write gets the first 20 bytes of an update to simulate a delta patch.
func (d delta) Write(b []byte) (int, error) {
	if len(b) < 20 {
		return 0, errors.New("Update too small")
	}
	return d.Write(b[0:20])
}

func createDeltaFromArtifactUpdate(artifact io.Reader) (*delta, error) {

	ar := areader.NewReader(artifact)
	deltaGenerator := newDelta()
	rootfs := handlers.NewRootfsInstaller(3)
	rootfs.InstallHandler = func(r io.Reader, df *handlers.DataFile) error {
		_, err := io.Copy(deltaGenerator, r)
		return err
	}
	ar.RegisterHandler(rootfs)
	ar.ReadArtifact()
	return &deltaGenerator, nil
}

func addDeltaToArtifact(artifact io.Reader) (io.Reader, error) {

	// ar := awriter.NewWriter(artifact)
	// TODO - needs to get the arguments (WriteArtifactArgs) to be passed to augment-manifest-header.
	return nil, nil
}

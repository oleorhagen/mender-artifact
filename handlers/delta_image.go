package handlers

import (
	"path/filepath"

	"github.com/mendersoftware/mender-artifact/artifact"
	"github.com/pkg/errors"
)

type DeltaImage struct {
	*Rootfs
}

func NewDeltaImage(updateFile string) *DeltaImage {
	return &DeltaImage{
		&Rootfs{
			update:  &DataFile{Name: updateFile},
			version: 3,
		},
	}
}

func (d *DeltaImage) GetType() string {
	return "delta-image"
}

func (d *DeltaImage) ComposeHeader(args *ComposeHeaderArgs) error {
	updFiles := filepath.Base(d.update.Name)
	path := artifact.UpdateHeaderPath(args.No)
	if args.Augmented {
		if err := writeFiles(args.TarWriter, []string{updFiles},
			path); err != nil {
			return err
		}
	} else {
		// The header in a version 3 artifact will not contain the update,
		// and hence there is no files in the files list.
		if err := writeEmptyFiles(args.TarWriter, []string{updFiles},
			path); err != nil {
			return err
		}
	}
	if err := writeTypeInfoV3(&WriteInfoArgs{
		tarWriter:  args.TarWriter,
		updateType: "delta-image",
		dir:        path,
		depends:    args.TypeInfoDepends,
		provides:   args.TypeInfoProvides,
	}); err != nil {
		return errors.Wrap(err, "ComposeHeader: ")
	}
	// store empty meta-data
	// the file needs to be a part of artifact even if this one is empty
	sw := artifact.NewTarWriterStream(args.TarWriter)
	if err := sw.Write(nil, filepath.Join(path, "meta-data")); err != nil {
		return errors.Wrap(err, "update: can not store meta-data")
	}
	return nil
}

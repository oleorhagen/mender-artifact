// Copyright 2018 Northern.tech AS
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
	"path/filepath"

	"github.com/mendersoftware/mender-artifact/artifact"
	"github.com/pkg/errors"
)

type Delta struct {
	*Rootfs
	OldRootfsChecksum string // Old rootfs checksum depend.
	NewRootfsChecksum string // New rootfs checksum provide.
}

func NewDelta(updFile, newRootfsChecksum, oldRootfsChecksum string) *Delta {
	rootfs := NewRootfsV3(updFile)
	return &Delta{
		Rootfs:            rootfs,
		NewRootfsChecksum: newRootfsChecksum,
		OldRootfsChecksum: oldRootfsChecksum,
	}
}

func (d *Delta) GetType() string {
	return "delta-update"
}

func (d *Delta) ComposeHeader(args *ComposeHeaderArgs) error {
	path := artifact.UpdateHeaderPath(args.No)
	if err := writeFiles(args.TarWriter, []string{filepath.Base(d.update.Name)},
		path); err != nil {
		return err
	}
	if err := writeTypeInfoV3(&WriteInfoArgs{
		tarWriter:  args.TarWriter,
		updateType: d.GetType(),
		dir:        path,
		depends:    []artifact.TypeInfoDepends{{d.OldRootfsChecksum}},
		provides:   []artifact.TypeInfoProvides{{d.NewRootfsChecksum}},
	}); err != nil {
		return err
	}
	// store empty meta-data
	// the file needs to be a part of artifact even if this one is empty
	sw := artifact.NewTarWriterStream(args.TarWriter)
	if err := sw.Write(nil, filepath.Join(path, "meta-data")); err != nil {
		return errors.Wrap(err, "update: can not store meta-data")
	}
	return nil
}

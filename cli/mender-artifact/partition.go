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

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// TODO - add all the methods needed by the partitionFile struct.
type Imager interface {
	GetImagePath() string
	GetImageKey() string
	GetImageName() string
	Repack() error
}

// TODO - create a partitionFile error type

// PartitionReadWriteClosePacker wraps io.ReadWriteCloser with a Repack method
type PartitionReadWriteClosePacker interface {
	io.ReadWriteCloser
	Repack() error
}

type partition struct {
	offset string
	size   string
	path   string
	name   string
}

// OpenFile returns a partitionFile if the file is found on the image.
// TODO - embed this in the sd and artifact images
func (si *partition) OpenFile(fpath string) (partitionFile, error) {
	// Check if the directory exists
	// Create the partitionFile
	// pf :=
	return nil, nil
}

type sdImage struct {
	abprts  []partition // A/B partitions - these should be similar, so mirror operations on both.
	dp      partition   // Data partition
	key     string
	imgpath string
}

func NewSDImage(imgpath, key string) (*sdImage, error) {
	modcands, isArtifact, err := getCandidatesForModify(si.imgpath, si.key)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unpack sdimg")
	}
	if isArtifact {
		return nil, errors.New("wrong imagetype: artifact")
	}
	return sdImage{
		imgpath: imgpath,
		key:     key,
		abprts:  modcands[0:1],
		dp:      modcands[2],
	}, nil
}

func (si *sdImage) Repack() error {
	return nil
}

// TODO - decode IMAGE-type?
func (si *sdImage) Decode() error {
	modcands, isArtifact, err := getCandidatesForModify(si.imgpath, si.key)
	return nil
}

func (si *sdImage) MarshalSDImage() error {
	return nil
}

func (si *sdImage) UnmarshalSDImage(data []byte, v interface{}) error {
	return nil
}

// OpenFile returns Partitions, which is an array of partitionFiles,
// if the file is found on the image.
func (si *sdImage) OpenFile(fpath string) (Partitions, error) {
	if filepath.HasPrefix(fpath, "data") {
		pf, err := NewPartitionFile(si.dp.name, fpath, si.key)
		if err != nil {
			return nil, errors.Wrapf(err, "sdimg:")
		}
		return []partitionFile{pf}, nil
	}
	// File is supposed to be on the A/B partition
	pfa, err := NewPartitionFile(si.abprts[0].name, fpath, key)
	if err != nil {
		return nil, errors.Wrapf(err, "sdimg:")
	}
	pfb, err := NewPartitionFile(si.abprts[1].name, fpath, key)
	if err != nil {
		return nil, errors.Wrapf(err, "sdimg:")
	}
	return []partitionFile{pfa, pfb}, nil
}

type artifactImage struct {
	p       partition
	imgpath string
	key     string
}

func NewArtifactImage(imgpath, key) (*artifactImage, error) {
	modcands, isArt, err := getCandidatesForModify(imgpath, key)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unpack artifact")
	}
	if !isArt {
		return nil, errors.New("wrong imagetype: sdimg")
	}
	return &artifactImage{
		p:       modcands[0],
		imgpath: imgpath,
		key:     key,
	}, nil
}

// OpenFile returns a partitionFile around the file on the mender-image.
func (ai *artifactImage) OpenFile(fpath string) (Partitions, error) {
	pf, err := NewPartitionFile(ai.p.name, fpath, ai.key)
	if err != nil {
		return nil, errors.Wrapf(err, "artifact-image: ")
	}
	return []partitionFile{pf}, nil
}

// Close closes the underlying tempfiles, and repacks the artifact.
func (ai *artifactImage) Repack() error {
}

// Repack repacks the sdimg.
// func (p *sdImage) Repack() error {
// 	// make modified images part of sdimg again
// 	var ps []partition
// 	for _, pf := range p {
// 		ps = append(ps, pf.partition)
// 	}
// 	return repackSdimg(ps, ps[0].name)
// }

// parseImgPath parses cli input of the form
// path/to/[sdimg,mender]:/path/inside/img/file
// into path/to/[sdimg,mender] and path/inside/img/file
func parseImgPath(imgpath string) (imgname, fpath string, err error) {
	paths := strings.SplitN(imgpath, ":", 2)

	if len(paths) != 2 {
		return "", "", errors.New("failed to parse image path")
	}

	if len(paths[0]) < 2 {
		return "", "", errors.New("invalid image or artifact path given")
	}

	if len(paths[1]) == 0 {
		return "", "", errors.New("please enter a path into the image")
	}

	return paths[0], paths[1], nil
}

// PartitionPacker has the functionality to repack an image or artifact.
type PartitionPacker interface {
	Repack() error
}

// partitionFile wraps partition and implements ReadWriteCloser
type partitionFile struct {
	PartitionPacker // Embed the underlying type (artifact, sdimg).
	key             string
	imagefilepath   string
}

// NewPartitionFile wraps one of the image struct in a partitionFile, which
// is a convenience type for writing to and reading from and sdimg, or artifact-image.
// Implements io.ReadWriteCloser.
func NewPartitionFile(img Imager, fpath string) (*partitionFile, error) {
	// Check that the directory path exists
	cmd := exec.Command("debugfs", "-w", "-R", fmt.Sprintf("cd %s", filepath.Dir(fpath)), img.GetPartitionName())
	ep, _ := cmd.StderrPipe()
	if err = cmd.Start(); err != nil {
		return errors.Wrap(err, "debugfs: run debugfs script")
	}
	data, err := ioutil.ReadAll(ep)
	if err != nil {
		return nil, err
	}
	if len(data) != 0 && strings.Contains(string(data), "File not found") {
		return nil, errors.New("directory does not exist")
	}

	return nil, nil
}

// Write reads all bytes from b into the partitionFile using debugfs.
func (p *partitionFile) Write(b []byte) (int, error) {
	f, err := ioutil.TempFile("", "mendertmp")

	// ignore tempfile os-cleanup errors
	defer f.Close()
	defer os.Remove(f.Name())

	if err != nil {
		return 0, err
	}
	if _, err := f.WriteAt(b, 0); err != nil {
		return 0, err
	}

	err = debugfsReplaceFile(p.imagefilepath, f.Name(), p.path)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

// Read reads all bytes from the filepath on the partition image into b
func (p *partitionFile) Read(b []byte) (int, error) {
	str, err := debugfsCopyFile(p.imagefilepath, p.path)
	if err != nil {
		return 0, errors.Wrap(err, "ReadError: debugfsCopyFile failed")
	}
	data, err := ioutil.ReadFile(filepath.Join(str, filepath.Base(p.imagefilepath)))
	if err != nil {
		return 0, errors.Wrapf(err, "ReadError: ioutil.Readfile failed to read file: %s", filepath.Join(str, filepath.Base(p.imagefilepath)))
	}
	defer os.RemoveAll(str) // ignore error removing tmp-dir
	return copy(b, data), io.EOF
}

// Close removes the temporary file held by partitionFile path.
func (p *partitionFile) Close() error {
	if p != nil {
		os.Remove(p.path) // ignore error for tmp-dir
	}
	return nil
}

// Repack repacks the artifact.
func (p *partitionFile) Repack() error {
	err := repackArtifact(p.name, p.path,
		p.key, filepath.Base(p.name))
	os.Remove(p.path) // ignore error, file exists in /tmp only
	return err
}

// Partitions is a wrapper around partitionFile, so that
// a write is duplicated to both Partitions' files during a write
// TODO - what is a good name?
type Partitions []partitionFile

// Write writes a file to both sdimg Partitions.
func (p Partitions) Write(b []byte) (int, error) {
	for _, part := range p {
		n, err := part.Write(b)
		if err != nil {
			return n, err
		}
		if n != len(b) {
			return n, io.ErrShortWrite
		}
	}
	return len(b), nil
}

// Read reads a file from an image.
func (p Partitions) Read(b []byte) (int, error) {
	return p[0].Read(b)
}

// Close closes all the underlying partition-files. e.g one close for each partition.
func (p Partitions) Close() (err error) {
	for part := range p {
		err = part.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

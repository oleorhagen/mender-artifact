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
	"os"
	"regexp"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var isimg = regexp.MustCompile(`\.(mender|sdimg|uefiimg)`)

func Cat(c *cli.Context) (err error) {
	if c.NArg() != 1 {
		return cli.NewExitError(fmt.Sprintf("Got %d arguments, wants one", c.NArg()), 1)
	}
	if !isimg.MatchString(c.Args().First()) {
		return cli.NewExitError("The input image does not seem to be a valid image", 1)
	}
	r, err := NewPartitionFile(c.Args().First(), c.String("key"))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("failed to open the partition reader: err: %v", err), 1)
	}
	var w io.WriteCloser = os.Stdout
	if _, err = io.Copy(w, r); err != nil {
		return cli.NewExitError(fmt.Sprintf("failed to copy from: %s to stdout: err: %v", c.Args().First(), err), 1)
	}
	return r.Close()
}

func Copy(c *cli.Context) (err error) {
	var r io.ReadCloser
	var w io.WriteCloser

	switch parseCLIOptions(c) {
	case copyin:
		r, err = os.OpenFile(c.Args().First(), os.O_CREATE|os.O_RDWR, 0655)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%v", err), 1)
		}
		w, err = NewPartitionFile(c.Args().Get(1), c.String("key"))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%v", err), 1)
		}
	case copyinstdin:
		r = os.Stdin
		w, err = NewPartitionFile(c.Args().First(), c.String("key"))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%v", err), 1)
		}
	case copyout:
		r, err = NewPartitionFile(c.Args().First(), c.String("key"))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%v", err), 1)
		}
		w, err = os.OpenFile(c.Args().Get(1), os.O_CREATE|os.O_RDWR, 0655)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%v", err), 1)
		}
	case parseError:
		return cli.NewExitError(fmt.Sprintln("no artifact or sdimage provided"), 1)
	case argerror:
		return cli.NewExitError(fmt.Sprintf("got %d arguments, wants two", c.NArg()), 1)
	default:
		return cli.NewExitError("critical error", 1)
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err = w.Close(); err != nil {
		return err
	}
	return r.Close()
}

// ListFiles lists all the files in a current directory using the command:
// `mender-artifact ls <img>:/path/`
func ListFiles(c *cli.Context) (err error) {
	imgname, dpath, err := parseImgPath(c.Args().First())
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "ListFiles: "), 1)
	}
	modcands, isArtifact, err := getCandidatesForModify(imgname, []byte(c.String("key")))
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "ListFiles: "), 1)
	}
	for _, tmpPartition := range modcands {
		defer os.Remove(tmpPartition.path)
	}
	var img partition
	if isArtifact {
		img = modcands[0]
	} else {
		img = modcands[1] // RootfsA
	}
	s, err := debugfsRunls(dpath, img.path)
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "ListFiles: "), 1)
	}
	fmt.Println(s)
	return nil
}

// Install installs a file from the host filesystem onto either
// a mender artifact, or an sdimg.
func Install(c *cli.Context) (err error) {
	var r io.ReadCloser
	var w io.WriteCloser
	switch parseCLIOptions(c) {
	case copyin:
		var perm os.FileMode
		if c.Int("mode") == 0 {
			return cli.NewExitError("File permissions needs to be set, if you are simply copying, the cp command should fit your needs", 1)
		}
		perm = os.FileMode(c.Int("mode"))
		r, err = os.OpenFile(c.Args().First(), os.O_RDWR, perm)
		defer r.Close()
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%v", err), 1)
		}
		w, err = NewPartitionFile(c.Args().Get(1), c.String("key"))
		defer w.Close()
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%v", err), 1)
		}
		if _, err = io.Copy(w, r); err != nil {
			return cli.NewExitError(fmt.Sprintf("%v", err), 1)
		}
		return nil
	case parseError:
		return cli.NewExitError("No artifact or sdimg provided", 1)
	case argerror:
		return cli.NewExitError(fmt.Sprintf("got %d arguments, wants two", c.NArg()), 1)
	default:
		return cli.NewExitError("Unrecognized parse error", 1)
	}
}

// enumerate cli-options
const (
	copyin = iota
	copyinstdin
	copyout
	parseError
	argerror
	criterror
)

func parseCLIOptions(c *cli.Context) int {

	finfo, err := os.Stdin.Stat()
	if err != nil {
		return criterror
	}

	// no data on stdin
	if finfo.Mode()&os.ModeNamedPipe == 0 {
		if c.NArg() != 2 {
			return argerror
		}
		switch {

		case isimg.MatchString(c.Args().First()):
			return copyout

		case isimg.MatchString(c.Args().Get(1)):
			return copyin

		default:
			return parseError

		}
	} else {
		// data on stdin
		if isimg.MatchString(c.Args().First()) {
			return copyinstdin
		}
		return parseError
	}
}

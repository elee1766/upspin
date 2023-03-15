// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package disk provides a storage.Storage that stores data on local disk.
package disk // import "upspin.io/cloud/storage/disk"

import (
	"github.com/spf13/afero"
	"upspin.io/cloud/storage"
	"upspin.io/cloud/storage/aferofs"
	"upspin.io/errors"
)

// New initializes and returns a disk-backed storage.Storage with the given
// options. The single, required option is "basePath" that must be an absolute
// path under which all objects should be stored.
func New(opts *storage.Opts) (storage.Storage, error) {
	const op errors.Op = "cloud/storage/disk.New"

	base, ok := opts.Opts["basePath"]
	if !ok {
		return nil, errors.E(op, "the basePath option must be specified")
	}

	return aferofs.NewDialer(afero.NewBasePathFs(afero.NewOsFs(), base))(opts)
}

func init() {
	storage.Register("Disk", New)
}

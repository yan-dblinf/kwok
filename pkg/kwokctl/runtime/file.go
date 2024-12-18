/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runtime

import (
	"context"
	"io"
	"io/fs"
	"os"
	"strings"

	"sigs.k8s.io/kwok/pkg/kwokctl/dryrun"
	"sigs.k8s.io/kwok/pkg/kwokctl/pki"
	"sigs.k8s.io/kwok/pkg/utils/file"
)

// DownloadWithCache downloads the src file to the dest file.
func (c *Cluster) DownloadWithCache(ctx context.Context, cacheDir, src, dest string, mode fs.FileMode, quiet bool) error {
	if s := strings.SplitN(src, "#", 2); len(s) == 2 {
		if c.IsDryRun() && !c.ShouldDownload() {
			dryrun.PrintMessage("# Download %s and extract %s to %s", s[0], s[1], dest)
			return nil
		}
		return file.DownloadWithCacheAndExtract(ctx, cacheDir, s[0], dest, s[1], mode, quiet, true, c.IsDryRun())
	}

	if c.IsDryRun() && !c.ShouldDownload() {
		dryrun.PrintMessage("# Download %s to %s", src, dest)
		return nil
	}
	return file.DownloadWithCache(ctx, cacheDir, src, dest, mode, quiet, c.IsDryRun())
}

// GeneratePki generates the pki for kwokctl
func (c *Cluster) GeneratePki(pkiPath string, sans ...string) error {
	if c.IsDryRun() {
		dryrun.PrintMessage("# Generate PKI to %s", pkiPath)
		return nil
	}

	return pki.GeneratePki(pkiPath, sans...)
}

// CreateFile creates a file.
func (c *Cluster) CreateFile(name string) error {
	if c.IsDryRun() {
		dryrun.PrintMessage("touch %s", name)
		return nil
	}

	return file.Create(name)
}

// CopyFile copies a file from src to dst.
func (c *Cluster) CopyFile(oldpath, newpath string) error {
	if c.IsDryRun() {
		dryrun.PrintMessage("cp %s %s", oldpath, newpath)
		return nil
	}

	return file.Copy(oldpath, newpath)
}

// RenameFile renames a file.
func (c *Cluster) RenameFile(oldpath, newpath string) error {
	if c.IsDryRun() {
		dryrun.PrintMessage("mv %s %s", oldpath, newpath)
		return nil
	}

	return file.Rename(oldpath, newpath)
}

// AppendToFile appends content to a file.
func (c *Cluster) AppendToFile(name string, content []byte) error {
	if c.IsDryRun() {
		dryrun.PrintMessage("cat <<EOF >>%s\n%s\nEOF", name, string(content))
		return nil
	}

	return file.Append(name, content)
}

// Remove removes a file.
func (c *Cluster) Remove(name string) error {
	if c.IsDryRun() {
		dryrun.PrintMessage("rm %s", name)
		return nil
	}

	return file.Remove(name)
}

// RemoveAll removes a directory and all its contents.
func (c *Cluster) RemoveAll(name string) error {
	if c.IsDryRun() {
		dryrun.PrintMessage("rm -rf %s", name)
		return nil
	}

	return file.RemoveAll(name)
}

// OpenFile opens/creates a file for writing.
func (c *Cluster) OpenFile(name string) (io.WriteCloser, error) {
	if c.IsDryRun() {
		return dryrun.NewCatToFileWriter(name), nil
	}

	return file.Open(name)
}

// WriteFile writes content to a file.
func (c *Cluster) WriteFile(name string, content []byte) error {
	if c.IsDryRun() {
		dryrun.PrintMessage("cat <<EOF >%s\n%s\nEOF", name, string(content))
		return nil
	}

	return file.Write(name, content)
}

// WriteFileWithMode writes content to a file with the given mode.
func (c *Cluster) WriteFileWithMode(name string, content []byte, mode os.FileMode) error {
	if c.IsDryRun() {
		dryrun.PrintMessage("cat <<EOF >%s\n%s\nEOF", name, string(content))
		dryrun.PrintMessage("chmod 0%o %s", mode, name)
		return nil
	}

	return file.WriteWithMode(name, content, mode)
}

// MkdirAll creates a directory.
func (c *Cluster) MkdirAll(name string) error {
	if c.IsDryRun() {
		dryrun.PrintMessage("mkdir -p %s", name)
		return nil
	}

	return file.MkdirAll(name)
}

// EnsureBinary ensures the binary exists.
func (c *Cluster) EnsureBinary(ctx context.Context, name, binary string) (string, error) {
	config, err := c.Config(ctx)
	if err != nil {
		return "", err
	}
	conf := config.Options

	binaryPath := c.GetBinPath(name + conf.BinSuffix)
	err = c.DownloadWithCache(ctx, conf.CacheDir, binary, binaryPath, 0750, conf.QuietPull)
	if err != nil {
		return "", err
	}

	return binaryPath, nil
}

/*
 *     Copyright 2024 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

type FileSize = uint64

const (
	// FileSize1MiB represents the size of 1MiB.
	FileSize1MiB FileSize = 1024 * 1024
	// FileSize10MiB represents the size of 10MiB.
	FileSize10MiB FileSize = 10 * FileSize1MiB
	// FileSize100MiB represents the size of 100MiB.
	FileSize100MiB FileSize = 100 * FileSize1MiB
)

const (
	// defaultFileServerEndpoint is the default endpoint of the file server.
	defaultFileServerEndpoint = "http://file-server-v2.dragonfly-e2e.svc"
	// defaultFileServerLocalDir is the default local directory of the file server.
	defaultFileServerLocalDir = "/tmp/e2e-file-server"
)

// FileServer represents a file server.
type FileServer struct {
	// endpoint is the endpoint of the file server.
	endpoint *url.URL
	// localDir is the local directory of the file server.
	localDir string
}

// NewFileServer creates a new file server.
func NewFileServer(opts ...FileServerOption) (*FileServer, error) {
	fileServerOptions := &FileServerOptions{
		endpoint: defaultFileServerEndpoint,
		localDir: defaultFileServerLocalDir,
	}

	for _, opt := range opts {
		opt(fileServerOptions)
	}

	// parse the endpoint
	var u *url.URL
	u, err := url.Parse(fileServerOptions.endpoint)
	if err != nil {
		return nil, err
	}

	// ensure the localDir exists
	if err = os.MkdirAll(fileServerOptions.localDir, 0755); err != nil {
		return nil, err
	}

	fileServer := &FileServer{
		endpoint: u,
		localDir: fileServerOptions.localDir,
	}

	return fileServer, nil
}

// PrepareFile prepares a file by the size and returns File instance.
func (fs *FileServer) PrepareFile(size FileSize) (*File, error) {
	// create file in the localDir
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(data)
	fileName := hex.EncodeToString(hash[:])
	filePath := filepath.Join(fs.localDir, fileName)
	if err = os.WriteFile(filePath, data, 0644); err != nil {
		return nil, err
	}

	// get file info
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	file := &File{
		info:        info,
		downloadURL: fmt.Sprintf("%s/%s", fs.endpoint.String(), fileName),
	}

	return file, nil
}

// CleanFile removes the file by the file info.
func (fs *FileServer) CleanFile(info os.FileInfo) error {
	// remove the file in the local dir
	filepath := filepath.Join(fs.localDir, info.Name())
	return os.RemoveAll(filepath)
}

// CleanAll removes all files in the local dir.
func (fs *FileServer) CleanAll() error {
	return os.RemoveAll(fs.localDir)
}

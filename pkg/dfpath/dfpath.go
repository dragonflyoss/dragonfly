/*
 *     Copyright 2020 The Dragonfly Authors
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

//go:generate mockgen -destination mocks/dfpath_mock.go -source dfpath.go -package mocks

package dfpath

import (
	"io/fs"
	"os"
	"sync"

	"github.com/hashicorp/go-multierror"
)

// Dfpath is the interface used for init project path.
type Dfpath interface {
	CacheDir() string
	CacheDirMode() fs.FileMode
	LogDir() string
	PluginDir() string
}

// Dfpath provides init project path function.
type dfpath struct {
	cacheDir     string
	cacheDirMode fs.FileMode
	logDir       string
	pluginDir    string
}

// Cache of the dfpath.
var cache struct {
	sync.Once
	d   *dfpath
	err *multierror.Error
}

// Option is a functional option for configuring the dfpath.
type Option func(d *dfpath)

// WithCacheDir set the cache directory.
func WithCacheDir(dir string) Option {
	return func(d *dfpath) {
		d.cacheDir = dir
	}
}

// WithCacheDirMode sets the cacheDir mode
func WithCacheDirMode(mode fs.FileMode) Option {
	return func(d *dfpath) {
		d.cacheDirMode = mode
	}
}

// WithLogDir set the log directory.
func WithLogDir(dir string) Option {
	return func(d *dfpath) {
		d.logDir = dir
	}
}

// WithPluginDir set plugin directory.
func WithPluginDir(dir string) Option {
	return func(d *dfpath) {
		d.pluginDir = dir
	}
}

// New returns a new dfpath interface.
func New(options ...Option) (Dfpath, error) {
	cache.Do(func() {
		d := &dfpath{
			logDir:       DefaultLogDir,
			pluginDir:    DefaultPluginDir,
			cacheDir:     DefaultCacheDir,
			cacheDirMode: DefaultCacheDirMode,
		}

		for _, opt := range options {
			opt(d)
		}

		// Create log directory.
		if err := os.MkdirAll(d.logDir, fs.FileMode(0700)); err != nil {
			cache.err = multierror.Append(cache.err, err)
		}

		// Create plugin directory.
		if err := os.MkdirAll(d.pluginDir, fs.FileMode(0700)); err != nil {
			cache.err = multierror.Append(cache.err, err)
		}

		// Create cache directory.
		if err := os.MkdirAll(d.cacheDir, d.cacheDirMode); err != nil {
			cache.err = multierror.Append(cache.err, err)
		}

		cache.d = d
	})

	if cache.err.ErrorOrNil() != nil {
		return nil, cache.err
	}

	d := *cache.d
	return &d, nil
}

func (d *dfpath) CacheDir() string {
	return d.cacheDir
}

func (d *dfpath) CacheDirMode() fs.FileMode {
	return d.cacheDirMode
}

func (d *dfpath) LogDir() string {
	return d.logDir
}

func (d *dfpath) PluginDir() string {
	return d.pluginDir
}

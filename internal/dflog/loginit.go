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

package logger

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"d7y.io/dragonfly/v2/pkg/types"
)

type logInitMeta struct {
	fileName             string
	setSugaredLoggerFunc func(*zap.SugaredLogger)
	setLoggerFunc        func(log *zap.Logger)
}

func InitManager(logLevel string, console bool, dir string, rotateConfig LogRotateConfig) error {
	if console {
		return createConsoleLogger(logLevel)
	}

	logDir := filepath.Join(dir, types.ManagerName)
	var meta = []logInitMeta{
		{
			fileName:             CoreLogFileName,
			setSugaredLoggerFunc: SetCoreLogger,
		},
		{
			fileName:             GrpcLogFileName,
			setSugaredLoggerFunc: SetGrpcLogger,
		},
		{
			fileName:             GinLogFileName,
			setSugaredLoggerFunc: SetGinLogger,
		},
		{
			fileName:             GCLogFileName,
			setSugaredLoggerFunc: SetGCLogger,
		},
		{
			fileName:             JobLogFileName,
			setSugaredLoggerFunc: SetJobLogger,
		},
	}

	return createFileLogger(logLevel, meta, logDir, rotateConfig)
}

func InitScheduler(logLevel string, console bool, dir string, rotateConfig LogRotateConfig) error {
	if console {
		return createConsoleLogger(logLevel)
	}

	logDir := filepath.Join(dir, types.SchedulerName)
	var meta = []logInitMeta{
		{
			fileName:             CoreLogFileName,
			setSugaredLoggerFunc: SetCoreLogger,
		},
		{
			fileName:             GrpcLogFileName,
			setSugaredLoggerFunc: SetGrpcLogger,
		},
		{
			fileName:             GCLogFileName,
			setSugaredLoggerFunc: SetGCLogger,
		},
		{
			fileName:             JobLogFileName,
			setSugaredLoggerFunc: SetJobLogger,
		},
	}

	return createFileLogger(logLevel, meta, logDir, rotateConfig)
}

func InitDaemon(logLevel string, console bool, dir string, rotateConfig LogRotateConfig) error {
	if console {
		return createConsoleLogger(logLevel)
	}

	logDir := filepath.Join(dir, types.DaemonName)
	var meta = []logInitMeta{
		{
			fileName:             CoreLogFileName,
			setSugaredLoggerFunc: SetCoreLogger,
		},
		{
			fileName:             GrpcLogFileName,
			setSugaredLoggerFunc: SetGrpcLogger,
		},
		{
			fileName:             GinLogFileName,
			setSugaredLoggerFunc: SetGinLogger,
		},
		{
			fileName:             GCLogFileName,
			setSugaredLoggerFunc: SetGCLogger,
		},
	}

	return createFileLogger(logLevel, meta, logDir, rotateConfig)
}

func InitDfget(logLevel string, console bool, dir string, rotateConfig LogRotateConfig) error {
	if console {
		return createConsoleLogger(logLevel)
	}

	logDir := filepath.Join(dir, types.DfgetName)
	var meta = []logInitMeta{
		{
			fileName:             CoreLogFileName,
			setSugaredLoggerFunc: SetCoreLogger,
		},
		{
			fileName:             GrpcLogFileName,
			setSugaredLoggerFunc: SetGrpcLogger,
		},
	}

	return createFileLogger(logLevel, meta, logDir, rotateConfig)
}

func InitDfcache(logLevel string, dir string, rotateConfig LogRotateConfig) error {
	logDir := filepath.Join(dir, types.DfcacheName)
	var meta = []logInitMeta{
		{
			fileName:             CoreLogFileName,
			setSugaredLoggerFunc: SetCoreLogger,
		},
		{
			fileName:             GrpcLogFileName,
			setSugaredLoggerFunc: SetGrpcLogger,
		},
	}

	return createFileLogger(logLevel, meta, logDir, rotateConfig)
}

func createConsoleLogger(logLevel string) error {
	levels = nil
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if logLevel != "" {
		switch strings.ToLower(logLevel) {
		case "debug":
			config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		case "info":
			config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		case "warn":
			config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		case "error":
			config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		default:
			fmt.Printf("Warning: invalid log level '%s', using 'info' instead\n", logLevel)
		}
	}

	log, err := config.Build(zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel), zap.AddCallerSkip(1))
	if err == nil {
		sugar := log.Sugar()
		SetCoreLogger(sugar)
		SetGrpcLogger(sugar)
		SetGinLogger(sugar)
		SetGCLogger(sugar)
		SetStorageGCLogger(sugar)
		SetKeepAliveLogger(sugar)
		SetStatSeedLogger(log)
		SetDownloadLogger(log)
		SetJobLogger(sugar)
	}
	levels = append(levels, config.Level)
	startLoggerSignalHandler()
	return nil
}

func createFileLogger(logLevel string, meta []logInitMeta, logDir string, rotateConfig LogRotateConfig) error {
	levels = nil
	// create parent dir first
	_ = os.MkdirAll(logDir, fs.FileMode(0700))

	for _, m := range meta {
		log, level, err := CreateLogger(path.Join(logDir, m.fileName), false, false, logLevel, rotateConfig)
		if err != nil {
			return err
		}
		if m.setSugaredLoggerFunc != nil {
			m.setSugaredLoggerFunc(log.Sugar())
		} else {
			m.setLoggerFunc(log)
		}

		levels = append(levels, level)
	}

	startLoggerSignalHandler()
	return nil
}

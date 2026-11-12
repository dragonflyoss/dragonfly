/*
 *     Copyright 2026 The Dragonfly Authors
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

package injector

import (
	"time"
)

const (
	// ConfigReloadInterval is the wait time for config reload.
	ConfigReloadInterval time.Duration = 30 * time.Second

	// Domain is the domain name of the injector.
	Domain string = "dragonfly.io"

	// ConfigMapPath is mounted by config/injector/deployment.yaml.
	ConfigMapPath string = "/etc/dragonfly"

	// Namespace labels for injection control.
	NamespaceInjectLabelName  string = Domain + "/" + "inject"
	NamespaceInjectLabelValue string = "true"

	// Pod annotation for injection control.
	PodInjectAnnotationName  string = Domain + "/" + "inject"
	PodInjectAnnotationValue string = "true"

	// Dfdaemon unix sock config.
	SkipUnixSockInjectAnnotationName  string = Domain + "/" + "skip-unix-sock-inject"
	SkipUnixSockInjectAnnotationValue string = "true"
	DfdaemonUnixSockVolumeName        string = "dfdaemon-unix-sock"
	DfdaemonUnixSockPath              string = "/var/run/dragonfly/dfdaemon.sock"

	InitContainerImageAnnotation string = Domain + "/" + "init-container-image"
	InitContainerName            string = "dragonfly-binaries"
	BinaryVolumeName             string = InitContainerName + "-" + "volume"
	BinaryVolumeMountPath        string = "/dragonfly/bin"
	BinaryDirPath                string = "/usr/local/bin"
	BinaryMountDirPath           string = "/usr/local/bin"

	DfgetBinaryName    string = "dfget"
	DfcacheBinaryName  string = "dfcache"
	DfstoreBinaryName  string = "dfstore"
	DfctlBinaryName    string = "dfctl"
	DfdaemonBinaryName string = "dfdaemon"
)

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

package cmd

import (
	"os"

	"d7y.io/dragonfly/v2/cdnsystem/config"
	"d7y.io/dragonfly/v2/cdnsystem/server"
	"d7y.io/dragonfly/v2/cmd/common"
	logger "d7y.io/dragonfly/v2/pkg/dflog"
	"d7y.io/dragonfly/v2/pkg/dflog/logcore"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	cfgFile string
	cfg     *config.Config
)

const (
	cdnSystemEnvPrefix = "cdn"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cdn",
	Short: "the cdn system of dragonfly",
	Long: `cdn system is a long-running process and is mainly responsible
for caching downloaded data to avoid downloading the same files
from remote source repeatedly.`,
	Args:              cobra.NoArgs,
	DisableAutoGenTag: true,
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := logcore.InitCdnSystem(cfg.Console); err != nil {
			return errors.Wrap(err, "init cdn system logger")
		}

		return runCdnSystem()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func init() {
	// Initialize default cdn system config
	cfg = config.New()

	// Initialize cobra
	common.InitCobra(rootCmd, cfgFile, cdnSystemEnvPrefix, cfg)
}

func runCdnSystem() error {
	// cdn system config values
	s, _ := yaml.Marshal(cfg)
	logger.Infof("cdn system configuration:\n%s", string(s))

	// initialize verbose mode
	common.InitVerboseMode(cfg.Verbose, cfg.PProfPort)

	if svr, err := server.New(cfg); err != nil {
		return err
	} else {
		return svr.Serve()
	}
}

/*
 *     Copyright 2022 The Dragonfly Authors
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
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"d7y.io/dragonfly/v2/client/config"
	"d7y.io/dragonfly/v2/client/dfcache"
	"d7y.io/dragonfly/v2/pkg/rpc/dfdaemon/client"
)

const statDesc = "stat checks if a file exists in P2P cache system"

// statCmd represents the cache stat command
var statCmd = &cobra.Command{
	Use:                "stat <-i cid>",
	Short:              statDesc,
	Long:               statDesc,
	Args:               cobra.NoArgs,
	DisableAutoGenTag:  true,
	SilenceUsage:       true,
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		return runDfcacheSubcmd(ctx, config.CmdStat, args)
	},
}

func initStat() {
	// Add the command to parent
	rootCmd.AddCommand(statCmd)

	flags := statCmd.Flags()
	flags.BoolVarP(&dfcacheConfig.LocalOnly, "local", "l", dfcacheConfig.LocalOnly, "only check task exists locally, and don't check other peers in P2P network")
	if err := viper.BindPFlags(flags); err != nil {
		panic(fmt.Errorf("bind cache stat flags to viper: %w", err))
	}
}

func runStat(cfg *config.DfcacheConfig, client client.V1) error {
	return dfcache.Stat(cfg, client)
}

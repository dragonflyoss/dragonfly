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

	"github.com/spf13/cobra"

	"d7y.io/dragonfly/v2/client/config"
	"d7y.io/dragonfly/v2/client/dfcache"
	"d7y.io/dragonfly/v2/pkg/rpc/dfdaemon/client"
)

const deleteDesc = "delete file from P2P cache system"

// deleteCmd represents the cache delete command
var deleteCmd = &cobra.Command{
	Use:                "delete <-i cid> [flags]",
	Short:              deleteDesc,
	Long:               deleteDesc,
	Args:               cobra.NoArgs,
	DisableAutoGenTag:  true,
	SilenceUsage:       true,
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		return runDfcacheSubcmd(ctx, config.CmdDelete, args)
	},
}

func initDelete() {
	// Add the command to parent
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(cfg *config.DfcacheConfig, client client.V1) error {
	return dfcache.Delete(cfg, client)
}

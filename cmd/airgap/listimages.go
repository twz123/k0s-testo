/*
Copyright 2021 k0s authors

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

package airgap

import (
	"fmt"

	"github.com/k0sproject/k0s/cmd/internal"
	"github.com/k0sproject/k0s/pkg/airgap"
	"github.com/k0sproject/k0s/pkg/config"

	"github.com/spf13/cobra"
)

func newAirgapListImagesCmd() *cobra.Command {
	var (
		debugFlags internal.DebugFlags
		all        bool
	)

	cmd := &cobra.Command{
		Use:              "list-images",
		Short:            "List image names and versions needed for airgapped installations",
		Example:          `k0s airgap list-images`,
		Args:             cobra.NoArgs,
		PersistentPreRun: debugFlags.Run,
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts, err := config.GetCmdOpts(cmd)
			if err != nil {
				return err
			}

			clusterConfig, err := opts.K0sVars.NodeConfig()
			if err != nil {
				return fmt.Errorf("failed to get config: %w", err)
			}

			out := cmd.OutOrStdout()
			for _, uri := range airgap.GetImageURIs(clusterConfig.Spec, all) {
				if _, err := fmt.Fprintln(out, uri); err != nil {
					return err
				}
			}
			return nil
		},
	}

	debugFlags.AddToFlagSet(cmd.PersistentFlags())

	flags := cmd.Flags()
	flags.AddFlagSet(config.GetPersistentFlagSet())
	flags.AddFlagSet(config.FileInputFlag())
	flags.BoolVar(&all, "all", false, "include all images, even if they are not used in the current configuration")

	return cmd
}

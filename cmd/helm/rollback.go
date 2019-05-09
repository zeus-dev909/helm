/*
Copyright The Helm Authors.

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

package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"helm.sh/helm/cmd/helm/require"
	"helm.sh/helm/pkg/action"
)

const rollbackDesc = `
This command rolls back a release to a previous revision.

The first argument of the rollback command is the name of a release, and the
second is a revision (version) number. To see revision numbers, run
'helm history RELEASE'.
`

func newRollbackCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	client := action.NewRollback(cfg)

	cmd := &cobra.Command{
		Use:   "rollback [RELEASE] [REVISION]",
		Short: "roll back a release to a previous revision",
		Long:  rollbackDesc,
		Args:  require.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := client.Run(args[0])
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "Rollback was a success! Happy Helming!\n")

			return nil
		},
	}

	f := cmd.Flags()
	f.IntVarP(&client.Version, "version", "v", 0, "revision number to rollback to (default: rollback to previous release)")
	f.BoolVar(&client.DryRun, "dry-run", false, "simulate a rollback")
	f.BoolVar(&client.Recreate, "recreate-pods", false, "performs pods restart for the resource if applicable")
	f.BoolVar(&client.Force, "force", false, "force resource update through delete/recreate if needed")
	f.BoolVar(&client.DisableHooks, "no-hooks", false, "prevent hooks from running during rollback")
	f.DurationVar(&client.Timeout, "timeout", 300, "time to wait for any individual Kubernetes operation (like Jobs for hooks)")
	f.BoolVar(&client.Wait, "wait", false, "if set, will wait until all Pods, PVCs, Services, and minimum number of Pods of a Deployment are in a ready state before marking the release as successful. It will wait for as long as --timeout")

	return cmd
}

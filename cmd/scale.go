/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/zeroryuki/k6/api/v1"
	"github.com/zeroryuki/k6/api/v1/client"
	"github.com/zeroryuki/k6/ui"
)

// scaleCmd represents the scale command
var scaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale a running test",
	Long: `Scale a running test.

  Use the global --address flag to specify the URL to the API server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vus := getNullInt64(cmd.Flags(), "vus")
		max := getNullInt64(cmd.Flags(), "max")
		if !vus.Valid && !max.Valid {
			return errors.New("Specify either -u/--vus or -m/--max")
		}

		c, err := client.New(address)
		if err != nil {
			return err
		}
		status, err := c.SetStatus(context.Background(), v1.Status{VUs: vus, VUsMax: max})
		if err != nil {
			return err
		}
		ui.Dump(stdout, status)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(scaleCmd)

	scaleCmd.Flags().Int64P("vus", "u", 1, "number of virtual users")
	scaleCmd.Flags().Int64P("max", "m", 0, "max available virtual users")
}

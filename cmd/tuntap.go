/*
Copyright © 2023 weak_ptr <weak_ptr@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"time"

	"github.com/nnnewb/nf/internal/tuntap"
	"github.com/spf13/cobra"
)

// tuntapCmd represents the tuntap command
var tuntapCmd = &cobra.Command{
	Use:   "tuntap",
	Short: "tuntap create tun/tap virtual device",
	Long:  `tuntap create tun/tap virtual device`,
	Run: func(cmd *cobra.Command, args []string) {
		f, err := tuntap.CAllocateTun("vtun1")
		cobra.CheckErr(err)
		defer f.Close()
		time.Sleep(time.Minute)
	},
}

func init() {
	rootCmd.AddCommand(tuntapCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tuntapCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tuntapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

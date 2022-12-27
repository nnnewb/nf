/*
Copyright © 2022 weak_ptr <weak_ptr@outlook.com>

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
	"fmt"

	"github.com/google/gopacket/pcap"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list network interface card.",
	Long:  `list network interface card.`,
	Run: func(cmd *cobra.Command, args []string) {
		ifs, err := pcap.FindAllDevs()
		if err != nil {
			panic(err)
		}
		for _, i := range ifs {
			println("---")
			fmt.Printf("Name: %s:\n", i.Name)
			fmt.Printf("Description: %s\n", i.Description)
			print("Address: ")
			for idx, addr := range i.Addresses {
				fmt.Printf("%s", addr.IP)
				if idx < len(i.Addresses)-1 {
					print(", ")
				}
			}
			println()
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

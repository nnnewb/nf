package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list network interface card.",
	Long:  `list network interface card.`,
	Run: func(cmd *cobra.Command, args []string) {

		ifs, err := net.Interfaces()
		cobra.CheckErr(err)

		for _, i := range ifs {
			println("---")
			fmt.Printf("Name: %s:\n", i.Name)
			print("Address: ")
			addrs, err := i.Addrs()
			cobra.CheckErr(err)

			for idx, addr := range addrs {
				fmt.Printf("%s", addr.String())
				if idx < len(addrs)-1 {
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

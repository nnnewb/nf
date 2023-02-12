/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"

	"github.com/songgao/water"
	"github.com/spf13/cobra"
)

// mktunCmd represents the mktun command
var mktunCmd = &cobra.Command{
	Use:   "mktun",
	Short: "make a tun device",
	Long:  `make a tun device`,
	Run: func(cmd *cobra.Command, args []string) {
		iface, err := water.New(water.Config{
			DeviceType: water.TUN,
			PlatformSpecificParams: water.PlatformSpecificParams{
				InterfaceName: "tun99",
				Network:       "192.168.19.2/24",
			},
		})
		cobra.CheckErr(err)
		var buffer [2048]byte
		for {
			n, err := iface.Read(buffer[:])
			cobra.CheckErr(err)
			s := hex.EncodeToString(buffer[:n])
			println(s)
		}
	},
}

func init() {
	rootCmd.AddCommand(mktunCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mktunCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mktunCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/songgao/water"
	"github.com/spf13/cobra"
)

// mktunCmd represents the mktun command
var mktunCmd = &cobra.Command{
	Use:   "mktun",
	Short: "make a tun device",
	Long:  `make a tun device`,
	Run: func(cmd *cobra.Command, args []string) {
		dest, err := cmd.Flags().GetIP("dest")
		cobra.CheckErr(err)

		port, err := cmd.Flags().GetInt16("port")
		cobra.CheckErr(err)

		conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", dest, port))
		cobra.CheckErr(err)

		iface, err := water.New(water.Config{
			DeviceType: water.TUN,
			PlatformSpecificParams: water.PlatformSpecificParams{
				Name: "tun99",
			},
		})
		cobra.CheckErr(err)

		file, err := os.OpenFile("traffic.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
		cobra.CheckErr(err)
		defer file.Close()

		var buffer [2048]byte
		for {
			n, err := iface.Read(buffer[:])
			cobra.CheckErr(err)
			pkt := gopacket.NewPacket(buffer[:n], layers.LayerTypeIPv4, gopacket.NoCopy)
			ipv4Layer := pkt.Layer(layers.LayerTypeIPv4)
			if ipv4Layer != nil {
				if ipv4Packet, ok := ipv4Layer.(*layers.IPv4); ok {
					log.Println(ipv4Packet.SrcIP, "->", ipv4Packet.DstIP, "carrying", ipv4Packet.NextLayerType(), "length", ipv4Packet.Length)
					_, err = conn.Write(buffer[:n])
					cobra.CheckErr(err)
					continue
				}
				log.Printf("layer to packet conversion failed, %+v\n", ipv4Layer)
			}

			log.Println("packet does not contains ipv4 layer")
		}
	},
}

func init() {
	rootCmd.AddCommand(mktunCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mktunCmd.PersistentFlags().String("foo", "", "A help for foo")
	mktunCmd.Flags().IP("dest", net.IPv4(0, 0, 0, 0), "tunnel server")
	mktunCmd.Flags().Int16("port", 0, "tunnel server port")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mktunCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

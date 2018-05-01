package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// killCmd represents the new command
var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill Your Sentinel SOCKS5 Node",
	Run: func(cmd *cobra.Command, args []string) {
		KillSocks5Node()
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
}

func KillSocks5Node() {
	cmd := "sudo killall ssserver"
	cmdParts := strings.Fields(cmd)

	killSocks := exec.Command(cmdParts[0], cmdParts[1:]...)

	if err := killSocks.Start(); err != nil {
		// fmt.Errorf("Could Not Start the Shadowsocks server: %v", err)
		panic(err)
	}

	fmt.Println("killed successfully")
}

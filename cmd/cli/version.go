package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print detailed version and build information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\033[36m%s\033[0m\n", banner)
		fmt.Println("Version Information:")
		fmt.Println("────────────────────────────────────")
		fmt.Printf("  Version:    %s\n", Version)
		fmt.Printf("  Git Commit: %s\n", GitCommit)
		fmt.Printf("  Build Time: %s\n", BuildTime)
		fmt.Printf("  Go Version: %s\n", runtime.Version())
		fmt.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println("────────────────────────────────────")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

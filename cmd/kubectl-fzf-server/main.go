package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"kubectlfzf/pkg/daemon"
	"kubectlfzf/pkg/httpserver"
	"kubectlfzf/pkg/k8s/resourcewatcher"
	"kubectlfzf/pkg/k8s/store"
	"kubectlfzf/pkg/kubectlfzfserver"
	"kubectlfzf/pkg/util"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	version   = "dev"
	gitCommit = "none"
	gitBranch = "unknown"
	goVersion = "unknown"
	buildDate = "unknown"
)

func versionFun(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Git hash: %s\n", gitCommit)
	fmt.Printf("Git branch: %s\n", gitBranch)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Go Version: %s\n", goVersion)
	os.Exit(0)
}

func startDaemonFun(cmd *cobra.Command, args []string) {
	daemon.StartDaemon()
}

func kubectlFzfServerFun(cmd *cobra.Command, args []string) {
	kubectlfzfserver.StartKubectlFzfServer()
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "kubectl_fzf_server",
		Run: kubectlFzfServerFun,
	}
	rootFlags := rootCmd.PersistentFlags()
	store.SetStoreConfigCli(rootFlags)
	httpserver.SetHttpServerConfigFlags(rootFlags)
	resourcewatcher.SetResourceWatcherCli(rootFlags)
	util.SetCommonCliFlags(rootFlags, "info")
	err := viper.BindPFlags(rootFlags)
	util.FatalIf(err)

	versionCmd := &cobra.Command{
		Use:   "version",
		Run:   versionFun,
		Short: "Print command version",
	}
	rootCmd.AddCommand(versionCmd)

	daemonCmd := &cobra.Command{
		Use: "daemon",
		Run: startDaemonFun,
	}
	daemonFlags := daemonCmd.Flags()
	daemon.SetDaemonFlags(daemonFlags)
	rootCmd.AddCommand(daemonCmd)
	err = viper.BindPFlags(daemonFlags)
	util.FatalIf(err)

	util.ConfigureViper()
	cobra.OnInitialize(util.CommonInitialization)
	defer pprof.StopCPUProfile()
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Root command failed: %v", err)
	}
}

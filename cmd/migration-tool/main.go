package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	interfaces "github.com/BeesNestInc/CassetteOS-Common"
	"github.com/BeesNestInc/CassetteOS-Common/utils/systemctl"
	"github.com/BeesNestInc/CassetteOS-MessageBus/common"
)

const (
	messageBusConfigDirPath  = "/etc/cassetteos"
	messageBusConfigFilePath = "/etc/cassetteos/message-bus.conf"
	messageBusName           = "cassetteos-message-bus.service"
	messageBusNameShort      = "message-bus"
)

//go:embedded ../../build/sysroot/etc/cassetteos/message-bus.conf.sample
// var _messageBusConfigFileSample string

var (
	commit = "private build"
	date   = "private build"

	_logger *Logger
)

// var _status *version.GlobalMigrationStatus

func main() {
	versionFlag := flag.Bool("v", false, "version")
	debugFlag := flag.Bool("d", true, "debug")
	forceFlag := flag.Bool("f", false, "force")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("v%s\n", common.MessageBusVersion)
		os.Exit(0)
	}

	println("git commit:", commit)
	println("build date:", date)

	_logger = NewLogger()

	if os.Getuid() != 0 {
		_logger.Info("Root privileges are required to run this program.")
		os.Exit(1)
	}

	if *debugFlag {
		_logger.DebugMode = true
	}

	if !*forceFlag {
		isRunning, err := systemctl.IsServiceRunning(messageBusName)
		if err != nil {
			_logger.Error("Failed to check if %s is running", messageBusName)
			panic(err)
		}

		if isRunning {
			_logger.Info("%s is running. If migration is still needed, try with -f.", messageBusName)
			os.Exit(1)
		}
	}

	migrationTools := []interfaces.MigrationTool{
		NewMigrationDummy(),
	}

	var selectedMigrationTool interfaces.MigrationTool

	// look for the right migration tool matching current version
	for _, tool := range migrationTools {
		migrationNeeded, err := tool.IsMigrationNeeded()
		if err != nil {
			panic(err)
		}

		if migrationNeeded {
			selectedMigrationTool = tool
			break
		}
	}

	if selectedMigrationTool == nil {
		_logger.Info("No migration to proceed.")
		return
	}

	if err := selectedMigrationTool.PreMigrate(); err != nil {
		panic(err)
	}

	if err := selectedMigrationTool.Migrate(); err != nil {
		panic(err)
	}

	if err := selectedMigrationTool.PostMigrate(); err != nil {
		_logger.Error("Migration succeeded, but post-migration failed: %s", err)
	}
}

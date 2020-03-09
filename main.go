package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/morfien101/asg-lifecyle-hook-manager/ec2metadatareader"
	"github.com/morfien101/asg-lifecyle-hook-manager/hookmanager"
)

var (
	version         = "0.0.4"
	actionAbandon   = "ABANDON"
	actionHeartBeat = "HEARTBEAT"
	actionContinue  = "CONTINUE"

	helpBlurb = `
	This application is used to interact with AWS Autoscaling Lifecycle hooks.
	It can set the hooks to Abandon, Continue or send a heartbeat.
	Only a single action can be invoked in a single run.
	It will consume credentials from instance roles or ENV vars.
	There is no provision for manually feeding in credentials and never will be.
	`

	versionFlag    = flag.Bool("v", false, "Show the version")
	helpFlag       = flag.Bool("h", false, "Show the help menu")
	verboseFlag    = flag.Bool("verbose", false, "Will log success statements as well as errors")
	asgNameFlag    = flag.String("n", "", "Name of the autoscaling group")
	hookNameFlag   = flag.String("l", "", "Name of the Lifecycle hook")
	instanceIDFlag = flag.String("i", "", "instance_id for the EC2 instance. If - is passed the instance ID is determined automatically from the metadata if available")
	hookActionFlag = flag.String("a", "", fmt.Sprintf("Set the lifecycle hook action. Valid values: %s, %s, %s", actionAbandon, actionContinue, actionHeartBeat))
)

func main() {
	flag.Parse()
	// These 2 functions have the ability to exit the app
	showStopperFlags()
	validateActions()

	// Do the work
	run()
}

func showStopperFlags() {
	if *helpFlag {
		fmt.Println(helpBlurb)
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}
}

func validateActions() {
	errors := []string{}
	if err := validateHookAction(*hookActionFlag); err != nil {
		errors = append(errors, err.Error())
	}
	if err := validateRequiredVars(); err != nil {
		errors = append(errors, err.Error())
	}
	if len(errors) != 0 {
		fmt.Println(strings.Join(errors, "\n"))
		os.Exit(1)
	}
}

func validateHookAction(hookAction string) error {
	if hookAction == "" {
		return fmt.Errorf("-a must to be specified")
	}
	hookActionValid := false
	switch hookAction {
	case actionAbandon:
		hookActionValid = true
	case actionHeartBeat:
		hookActionValid = true
	case actionContinue:
		hookActionValid = true
	}
	if !hookActionValid {
		return fmt.Errorf("Hook action %s is not valid", hookAction)
	}

	return nil
}

func validateRequiredVars() error {
	errors := []string{}
	if *asgNameFlag == "" {
		errors = append(errors, "-n autoscaling group name must be specified")
	}
	if *hookNameFlag == "" {
		errors = append(errors, "-l lifecycle hook name must be specified")
	}
	if *instanceIDFlag == "" {
		errors = append(errors, "-i instance_id must be specified")
	}

	if len(errors) != 0 {
		return fmt.Errorf("%s", strings.Join(errors, ","))
	}
	return nil
}

func run() {
	instanceID := ""
	if *instanceIDFlag == "-" {
		localInstanceID, err := ec2metadatareader.InstanceID()
		if err != nil {
			writeToStdErr(fmt.Sprintf("Could not determine instance id. Error: %s", err))
			os.Exit(1)
		}
		instanceID = localInstanceID
	} else {
		instanceID = *instanceIDFlag
	}

	switch *hookActionFlag {
	case actionAbandon:
		output, err := hookmanager.SetAbandon(*asgNameFlag, *hookNameFlag, instanceID)
		if err != nil {
			terminate(err.Error(), 1)
			return
		}
		terminate(output, 0)

	case actionContinue:
		output, err := hookmanager.SetContinue(*asgNameFlag, *hookNameFlag, instanceID)
		if err != nil {
			terminate(err.Error(), 1)
			return
		}
		terminate(output, 0)

	case actionHeartBeat:
		output, err := hookmanager.RecordHeartBeat(*asgNameFlag, *hookNameFlag, instanceID)
		if err != nil {
			terminate(err.Error(), 1)
			return
		}
		terminate(output, 0)
	}
}

func terminate(log string, exitCode int) {
	if exitCode == 0 {
		verboseLog(log)
	} else {
		writeToStdErr(log)
	}
	os.Exit(exitCode)
}

func writeToStdErr(s string) {
	fmt.Fprintln(os.Stderr, s)
}

func verboseLog(s string) {
	if *verboseFlag {
		fmt.Println(s)
	}
}

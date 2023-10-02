package commands

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"

	"github.com/denizgursoy/gotouch/internal/operator"
)

type (
	flag struct {
		name      string
		shortName string
	}
	goTouchCommand struct {
		commandName string
		flags       []flag
	}
)

func TestCreateRootCommand(t *testing.T) {
	t.Run("should have gotouch command", func(t *testing.T) {
		instance := operator.GetInstance()
		expectedCommand := goTouchCommand{
			commandName: "gotouch",
			flags: []flag{
				{
					shortName: "f",
					name:      "file",
				},
				{
					shortName: "i",
					name:      "inline",
				},
			},
		}

		actualCommand := CreateRootCommand(instance, BuildInfo{})
		checkCommandHasCorrectValues(t, expectedCommand, actualCommand)
	})

	t.Run("should have package command", func(t *testing.T) {
		instance := operator.GetInstance()
		expectedCommand := goTouchCommand{
			commandName: "package",
			flags: []flag{
				{
					shortName: "s",
					name:      "source",
				},
				{
					shortName: "t",
					name:      "target",
				},
			},
		}

		actualCommand := CreatePackageCommand(instance)
		checkCommandHasCorrectValues(t, expectedCommand, actualCommand)
	})
}

func checkCommandHasCorrectValues(t *testing.T, expectedCommand goTouchCommand, actualCommand *cobra.Command) {
	if expectedCommand.commandName == actualCommand.Name() {
		actualCommand.Flags().VisitAll(func(actualFlag *pflag.Flag) {
			flagFound := false
			for _, commandFlag := range expectedCommand.flags {
				if commandFlag.name == actualFlag.Name {
					flagFound = true
					require.EqualValues(t, commandFlag.shortName, actualFlag.Shorthand)
				}
			}
			require.Truef(t, flagFound, "could not find %s actualFlag on %s command ", actualFlag.Name, expectedCommand.commandName)
		})
	} else {
		require.Fail(t, "command name is not same")
	}
}

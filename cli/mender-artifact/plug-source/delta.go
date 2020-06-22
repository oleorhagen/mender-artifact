
package main


import (
	"fmt"

	"github.com/urfave/cli"
)

// This is the code we need to reproduce
// manifestos=""
// for manifesto in $MANIFESTOS; do
// if ! kubectl apply --dry-run -f $manifesto >/dev/null 2>&1; then
// echo "The ${manifesto} k8s manifesto is not valid. Aborting."
// show_help_and_exit_error
// fi
// manifestos="$manifestos -f ${manifesto}"
// done

func writeModuleImage(c *cli.Context) {
	// for _, manifesto := range c.StringSlice("manifestos") {
	// 	// Kubectl apply ...
	// }
	fmt.Println("Wrote the delta-artifact: " + c.String("artifact-name"))
	return
}


func CLI() cli.Command {

	//
	// Update modules: module-image
	//
	writeModuleCommand := cli.Command{
		Name:   "delta",
		Action: writeModuleImage,
		Usage:  "Writes Mender artifact for an delta update module",
		UsageText: "Writes a generic Mender artifact that will be used by an update module. " +
			"This command is not meant to be used directly, but should rather be wrapped by an " +
			"update module build command, which prepares all the necessary files and headers " +
			"for that update module.",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "from-artifact, from",
				Usage: "delta foo bar",
				Required: true,
			},
			cli.StringSliceFlag{
				Name:  "to-artifact, to",
				Usage: "delta foo bar",
				Required: true,
			},
		},
	}

	return writeModuleCommand

}

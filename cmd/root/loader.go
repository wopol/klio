package root

import (
	"plugin"

	"github.com/spf13/cobra"
	"github.com/g2a-com/klio/pkg/log"
)

func loadExternalCommand(rootCmd *cobra.Command, path string) {
	plug, err := plugin.Open(path)
	if err != nil {
		log.Warnf("cannot open plugin '%s': %v", path, err)
		return
	}

	symbol, err := plug.Lookup("NewCommand")
	if err != nil {
		log.Warnf("cannot find NewCommand symbol in the plugin '%s': %v", path, err)
		return
	}

	plugCommand, ok := symbol.(func() *cobra.Command)
	if !ok {
		log.Warnf("cannot find NewCommand symbol in the plugin '%s': %v", path, err)
		return
	}

	cmd := plugCommand()
	if cmd, _, _ := rootCmd.Find([]string{cmd.Name()}); cmd != rootCmd {
		log.Debugf("cannot register aleready registered command '%s' provided by '%s'", cmd.Name(), path)
		return
	}

	log.Debugf("registered command '%s' provided by '%s'", cmd.Name(), path)
	rootCmd.AddCommand(cmd)
}
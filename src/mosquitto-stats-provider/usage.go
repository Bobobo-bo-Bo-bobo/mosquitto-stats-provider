package main

import (
	"fmt"
	"runtime"
)

func showVersion() {

	fmt.Printf("mosquitto-stats-provider version %s\n"+
		"Copyright (C) by Andreas Maus <maus@ypbind.de>\n"+
		"This program comes with ABSOLUTELY NO WARRANTY.\n"+
		"\n"+
		"mosquitto-stats-provider is distributed under the Terms of the GNU General\n"+
		"Public License Version 3. (http://www.gnu.org/copyleft/gpl.html)\n"+
		"\n"+
		"Build with go version: %s\n"+
		"\n", version, runtime.Version())
}

func showUsage() {
	showVersion()
	fmt.Printf("Usage: mosquitto-stats-provider -config /path/to/config.file [-help] [-verbose] [-version]\n" +
		"\n" +
		"  -config /path/to/config.file Path to configuration file\n" +
		"  -help                        This text\n" +
		"  -verbose                     Verbose output\n" +
		"  -version                     Show version information\n" +
		"\n")
}

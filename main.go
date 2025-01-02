package main

import (
	"github.com/ryanparsa/gmail/cmd"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})

	logrus.SetLevel(logrus.DebugLevel)

}
func main() {
	cmd.Execute()
}

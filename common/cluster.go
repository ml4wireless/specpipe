package common

import "fmt"

var (
	ClusterSubjectPrefix string = "specpipe.cluster"
	DataSubjectPrefix    string = "specpipe.data"
)

type ClusterCmd string

var (
	RegisterCmd    ClusterCmd = "register"
	DeregisterCmd  ClusterCmd = "deregister"
	HealthCheckCmd ClusterCmd = "health"
)

type SDRType string

var (
	FM SDRType = "fm"
)

func ClusterSubject(sdrType SDRType, deviceName string, cmd ClusterCmd) string {
	return fmt.Sprintf("%s.%s.%s.%s", ClusterSubjectPrefix, sdrType, deviceName, cmd)
}

func DataSubject(sdrType SDRType, deviceName string) string {
	return fmt.Sprintf("%s.%s.%s", DataSubjectPrefix, sdrType, deviceName)
}

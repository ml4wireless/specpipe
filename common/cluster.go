package common

import (
	"fmt"
	"strings"
)

const (
	DataSubjectPrefix    string = "specpipe.data"    // stream subject
	IqDataSubjectPrefix  string = "specpipe-iq.data" // stream subject
	ClusterSubjectPrefix string = "specpipe-cluster" // simple subject

	KVStoreBucket string = "specpipe"
)

const OkMsg string = "ok"

type ClusterCmd string

const (
	HealthCheckCmd ClusterCmd = "health"
	WatchConfigCmd ClusterCmd = "watchcfg"
)

type SDRType string

const (
	FM SDRType = "fm"
	IQ SDRType = "iq"
)

func DataSubject(sdrType SDRType, deviceName string) string {
	if sdrType == IQ {
		return fmt.Sprintf("%s.%s.%s", IqDataSubjectPrefix, sdrType, deviceName)
	}
	return fmt.Sprintf("%s.%s.%s", DataSubjectPrefix, sdrType, deviceName)
}

func ClusterSubject(sdrType SDRType, deviceName string, cmd ClusterCmd) string {
	return fmt.Sprintf("%s.%s.%s.%s", ClusterSubjectPrefix, sdrType, deviceName, cmd)
}

func KVStoreKey(sdrType SDRType, deviceName string) string {
	return fmt.Sprintf("%s_%s", sdrType, deviceName)
}

func DeviceNameFromKey(key string) string {
	return strings.Split(key, "_")[1]
}

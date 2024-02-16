package common

type Device interface {
	Device()
}

type FMDevice struct {
	RegisterTs      int64   `json:"register_ts"`
	SpecpipeVersion string  `json:"specpipe_version"`
	Name            string  `json:"name"`
	Freq            string  `json:"freq"`
	SampleRate      string  `json:"sample_rate"`
	ResampleRate    string  `json:"resample_rate"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
}

func (*FMDevice) Device() {}

type IQDevice struct {
	RegisterTs      int64   `json:"register_ts"`
	SpecpipeVersion string  `json:"specpipe_version"`
	Name            string  `json:"name"`
	Freq            string  `json:"freq"`
	SampleRate      string  `json:"sample_rate"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	Forward         bool    `json:"forward"`
}

func (*IQDevice) Device() {}

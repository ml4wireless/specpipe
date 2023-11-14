package common

type Device interface {
	Device()
}

type FMDevice struct {
	Name       string  `json:"name"`
	Freq       string  `json:"freq"`
	SampleRate string  `json:"sample_rate"`
	Latitude   float32 `json:"latitude"`
	Longitude  float32 `json:"longitude"`
}

func (*FMDevice) Device() {}

type IQDevice struct {
	Name       string  `json:"name"`
	Freq       string  `json:"freq"`
	SampleRate string  `json:"sample_rate"`
	Latitude   float32 `json:"latitude"`
	Longitude  float32 `json:"longitude"`
}

func (*IQDevice) Device() {}

package common

type FMDevice struct {
	Name       string  `json:"name"`
	Freq       string  `json:"freq"`
	SampleRate string  `json:"sample_rate"`
	Latitude   float32 `json:"latitude"`
	Longitude  float32 `json:"longitude"`
}

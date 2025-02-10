package api

import (
	"servon/core/libs"
)

type Sample struct {
	sampleProvider libs.SampleProvider
}

func NewSample() Sample {
	return Sample{
		sampleProvider: libs.NewSampleProvider(),
	}
}

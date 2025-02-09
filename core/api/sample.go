package api

import (
	"servon/core/provider"
)

type Sample struct {
	sampleProvider provider.SampleProvider
}

func NewSample() Sample {
	return Sample{
		sampleProvider: provider.NewSampleProvider(),
	}
}

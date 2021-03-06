package main

import (
	"errors"
	"strconv"

	"github.com/JoshuaDoes/goprobe"
)

// VoiceServiceDirect exports the methods required to handle direct audio/video URLs
type VoiceServiceDirect struct {
}

// GetName returns the service's name
func (*VoiceServiceDirect) GetName() string {
	return "Direct URL"
}

// GetColor returns the service's color
func (*VoiceServiceDirect) GetColor() int {
	return 0x1C1C1C
}

// TestURL tests if the given URL has an audio stream
func (*VoiceServiceDirect) TestURL(url string) (bool, error) {
	probe, err := goprobe.ProbeMedia(url)
	if err != nil {
		return false, nil
	}

	if len(probe.Streams) == 0 {
		return false, nil
	}

	for _, stream := range probe.Streams {
		if stream.CodecType == "audio" {
			return true, nil
		}
	}

	return false, nil
}

// GetMetadata returns the metadata for a given direct audio/video URL
func (*VoiceServiceDirect) GetMetadata(url string) (*Metadata, error) {
	probe, err := goprobe.ProbeMedia(url)
	if err != nil {
		return nil, err
	}

	if len(probe.Streams) == 0 {
		return nil, errors.New("goprobe: no media streams found")
	}

	audioStream := &goprobe.Stream{}
	for _, stream := range probe.Streams {
		if stream.CodecType == "audio" {
			audioStream = stream
			break
		}
	}
	if audioStream == nil {
		return nil, errors.New("goprobe: no audio streams found")
	}

	duration, err := strconv.ParseFloat(audioStream.Duration, 64)
	if err != nil {
		return nil, err
	}

	metadata := &Metadata{
		Title:      url,
		DisplayURL: url,
		StreamURL:  url,
		Duration:   duration,
	}

	if audioStream.Tags != nil {
		if audioStream.Tags.Artist != "" && audioStream.Tags.Title != "" {
			trackArtist := &MetadataArtist{
				Name: audioStream.Tags.Artist,
				URL:  url,
			}
			metadata.Artists = append(metadata.Artists, *trackArtist)
			metadata.Title = audioStream.Tags.Title
		}
	} else {
		trackArtist := &MetadataArtist{
			Name: "Unknown",
			URL:  url,
		}
		metadata.Artists = append(metadata.Artists, *trackArtist)
	}

	return metadata, nil
}

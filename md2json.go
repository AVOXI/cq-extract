package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
)

// CallMetadata is metadata
type CallMetadata struct {
	MetadataType    string            `json:"metadataType"`
	ColumnQualifier string            `json:"columnQualifier"`
	Metadata        map[string]string `json:"Metadata"`
}

// Call is a call
type Call struct {
	CallId         string         `json:"callId"`
	UserId         string         `json:"UserId"`
	OrganizationId string         `json:"OrganizationId"`
	Metadata       []CallMetadata `json:"metadata"`
}

type MediaSource struct {
	TotalAudioEnergy     float64 `json:"totalAudioEnergy"`
	TotalSamplesDuration float64 `json:"totalSamplesDuration"`
}

type AudioTrack struct {
	JitterBufferDelay        float64 `json:"jitterBufferDelay"`
	JitterBufferEmittedCount float64 `json:"jitterBufferEmittedCount"`
	TotalAudioEnergy         float64 `json:"totalAudioEnergy"`
	TotalSamplesReceived     float64 `json:"totalSamplesReceived"`
	TotalSamplesDuration     float64 `json:"totalSamplesDuration"`
}

type InboundRTP struct {
	PacketsReceived int        `json:"packetsReceived"`
	PacketsLost     int        `json:"packetsLost"`
	Track           AudioTrack `json:"track"`
}

type SyntheticData struct {
	InboundRMS  float64 `json:"inbound_rms"`
	OutboundRMS float64 `json:"outbound_rms"`
}

const (
	CqTag   = "CALL_QUALITY_STATS"
	Version = "1.0.1"
)

func main() {
	file := os.Stdin
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()
	if *showVersion {
		fmt.Printf("Call Quality parser version: %s\n", Version)
		return
	}
	if len(flag.Args()) > 0 {
		path := flag.Args()[0]
		if f, err := os.Open(path); err != nil {
			log.Fatalf("Unable to open file [%v]: %v", path, err)
		} else {
			file = f
		}
	}
	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatalf("failed to read json: %v", err)
	}

	var c Call
	if err = json.Unmarshal(bytes, &c); err != nil {
		log.Fatalf("failed to parse json: %v", err)
	}

	fmt.Print("{")
	for _, m := range c.Metadata {
		if m.MetadataType == CqTag {
			count := len(m.Metadata)
			for k, cq := range m.Metadata {
				payload := strings.TrimSpace(cq)
				if strings.HasPrefix(payload, "[") || strings.HasPrefix(payload, "\"") {
					fmt.Printf("\"%v\": %v", k, cq)
				} else {
					// Add the missing double quotes
					fmt.Printf("\"%v\": \"%v\"", k, cq)
				}
				if count--; count > 0 {
					fmt.Print(",")
				}
			}
		}
	}

	sd := SyntheticData{}
	for _, m := range c.Metadata {
		if m.MetadataType == CqTag {
			if j, ok := m.Metadata["inbound-rtp"]; ok {
				var inboundRTP []InboundRTP
				err := json.Unmarshal([]byte(j), &inboundRTP)
				if err == nil && len(inboundRTP) > 0 {
					sd.InboundRMS = math.Sqrt(inboundRTP[0].Track.TotalAudioEnergy / inboundRTP[0].Track.TotalSamplesDuration)
				}
			}
			if j, ok := m.Metadata["media-source"]; ok {
				var mediaSrc []MediaSource
				err := json.Unmarshal([]byte(j), &mediaSrc)
				if err == nil && len(mediaSrc) > 0 {
					sd.OutboundRMS = math.Sqrt(mediaSrc[0].TotalAudioEnergy / mediaSrc[0].TotalSamplesDuration)
				}
			}
		}
	}
	if j, err := json.Marshal(&sd); err != nil {
		fmt.Print("}")
	} else {
		fmt.Printf(",\"synthetic-data\":%s}", string(j))
	}
}

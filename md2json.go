package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

const (
	CqTag = "CALL_QUALITY_STATS"
)

func main() {
	file := os.Stdin
	flag.Parse()
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
				fmt.Printf("\"%v\": %v", k, cq)
				if count--; count > 0 {
					fmt.Print(",")
				}
			}
		}
	}
	fmt.Print("}")

	//if _, err := fmt.Fprintf(os.Stderr, "%v\n", c); err != nil {
	//	log.Fatalf("unable to write")
	//}

}

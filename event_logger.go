package gohome

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome/log"
)

// EventLogger consumes events from the event bus and outputs them to
// the event log
type EventLogger struct {
	// Path the directory and file name where the log will be saved
	Path    string
	Verbose bool
}

func (c *EventLogger) ConsumerName() string {
	return "EventLogger"
}

func (c *EventLogger) StartConsuming(ch chan evtbus.Event) {
	log.V("EventLogger - start consuming events")

	go func() {
		f, err := os.Create(c.Path)
		if err != nil {
			log.E(fmt.Sprintf("failed to open event log for writing, log path: %s, err: %s", c.Path, err))
			return
		}
		defer f.Close()

		for e := range ch {
			var eventType string
			var data evtbus.Event

			switch evt := e.(type) {
			case *SensorAttrChangedEvt:
				eventType = "SensorAttrChangedEvt"
				data = evt
			case *ZoneLevelChangedEvt:
				eventType = "ZoneLevelChangedEvt"
				data = evt
			case *ButtonPressEvt:
				eventType = "ButtonPressEvt"
				data = evt
			case *ButtonReleaseEvt:
				eventType = "ButtonReleaseEvt"
				data = evt
			case *ClientConnectedEvt:
				eventType = "ClientConnectedEvt"
				data = evt
			case *ClientDisconnectedEvt:
				eventType = "ClientDisconnectedEvt"
				data = evt
			}

			// In verbose mode we log more information, useful for debugging
			if c.Verbose {
				switch evt := e.(type) {
				case *ZoneLevelReportingEvt:
					eventType = "ZoneLevelReportingEvt"
					data = evt
				case *SensorAttrReportingEvt:
					eventType = "SensorAttrReportingEvt"
					data = evt
				}
			}

			if data != nil {
				enc := json.NewEncoder(f)
				enc.Encode(struct {
					Type      string      `json:"type"`
					Timestamp string      `json:"timestamp"`
					Data      interface{} `json:"data"`
				}{
					Type:      eventType,
					Timestamp: time.Now().UTC().String(),
					Data:      data,
				})
			}
		}
		log.V("EventLogger - event channel has closed")
	}()
}

func (c *EventLogger) StopConsuming() {
	log.V("EventLogger - stop consuming events")
}
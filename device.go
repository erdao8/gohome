package gohome

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

type Device struct {
	Identifiable
	System         *System
	ConnectionInfo ConnectionInfo

	evpDone chan bool
	evpFire chan Event
	conn    Connection
}

func (d *Device) Connect() (Connection, error) {
	if d.conn != nil {
		return d.conn, nil
	}

	//TODO: When we support more than telnet, will use ConnectionInfo to
	//determine what type of connection we need to make
	conn := &TelnetConnection{}
	conn.SetConnectionInfo(d.ConnectionInfo)
	err := conn.Open()
	if err != nil {
		return nil, err
	}

	//TODO: Proper connection pooling
	d.conn = conn
	return conn, nil
}

func (d *Device) StartProducingEvents() (<-chan Event, <-chan bool) {
	//TODO: When to init these
	d.evpDone = make(chan bool)
	d.evpFire = make(chan Event)
	return d.evpFire, d.evpDone
}

//TODO: How to stop this?
func Stream(d *Device, r io.Reader) {
	scanner := bufio.NewScanner(r)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		//Match first instance of ~OUTPUT|~DEVICE.*\r\n
		str := string(data[0:])
		indices := regexp.MustCompile("[~|#][OUTPUT|DEVICE].+\r\n").FindStringIndex(str)
		//fmt.Printf("%s === %v\n", str, indices)

		//TODO: Don't let input grow forever - remove beginning chars after reaching max length

		if indices != nil {
			token = []byte(string([]rune(str)[indices[0]:indices[1]]))
			advance = indices[1]
			err = nil
		} else {
			advance = 0
			token = nil
			err = nil
		}
		return
	}

	scanner.Split(split)
	for scanner.Scan() {
		//fmt.Printf("scanner: %s\n", scanner.Text())

		if d.evpFire != nil {
			orig := scanner.Text()
			cmd := ParseCommandString(d, orig)
			d.evpFire <- NewEvent(d, cmd, orig)
		}
	}

	if d.evpDone != nil {
		close(d.evpDone)
	}

	//TODO: What about connect/disconnect event
	fmt.Println("Done scanning")
	if err := scanner.Err(); err != nil {
		fmt.Printf("something happened", err)
	}
}

package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

func CallNTPAsync(address string, quiet bool, channel chan *time.Time, group *sync.WaitGroup) {
	currentTime := CallNTP(address, quiet)

	if currentTime != nil {
		channel <- currentTime
	}

	group.Done()
}

// CallNTP Calls an NTP server according to [RFC 5905](https://datatracker.ietf.org/doc/html/rfc5905), and returns the current time.
// Uses the NTP v4
func CallNTP(address string, quiet bool) *time.Time {
	if !quiet {
		fmt.Printf("Calling NTP server at %s...\n", address)
	}

	packet := make([]byte, 48)

	// leap 0, version 4, mode 3 (client)
	// 0 4 3 -> 00 100 011 -> 0x23
	packet[0] = 0x23

	connection, err := net.Dial("udp", address)

	if err != nil {
		fmt.Println("An error occurred: ", err)
		return nil
	}

	defer connection.Close()

	// Set a deadline for the connection
	if deadlineExceededError := connection.SetDeadline(time.Now().Add(3 * time.Second)); deadlineExceededError != nil {
		fmt.Println("Deadline exceeded: ", deadlineExceededError)
		return nil
	}

	// Send the packet to the server
	if _, writeError := connection.Write(packet); writeError != nil {
		fmt.Println("Writing the packet failed: ", writeError)
		return nil
	}

	responsePacket := make([]byte, 48)

	// Read the response from the server
	if _, readError := connection.Read(responsePacket); readError != nil {
		fmt.Println("Reading the response failed: ", readError)
		return nil
	}

	ntpSeconds := binary.BigEndian.Uint32(responsePacket[40:44])

	// RFC 868: https://datatracker.ietf.org/doc/rfc868/
	unixSeconds := int64(ntpSeconds) - 2208988800

	unixTime := time.Unix(unixSeconds, 0)

	return &unixTime
}

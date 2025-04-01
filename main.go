package main

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"net"
	"os"
	"time"
)

// Calls an NTP server according to [RFC 5905](https://datatracker.ietf.org/doc/html/rfc5905), and returns the current time.
func callNTPAddress(address string, quiet bool) *time.Time {
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

func callNTP(hostname string, port int, quiet bool) *time.Time {
	return callNTPAddress(fmt.Sprintf("%s:%d", hostname, port), quiet)
}

func main() {
	var address string
	addressPresent := false

	var hostname string
	hostnamePresent := false

	var port int64
	portPresent := false

	var quiet bool

	command := &cli.Command{
		Name:  "ntp",
		Usage: "Get the current time from an NTP server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "address",
				Aliases:     []string{"a"},
				Usage:       "The address of the NTP server to call",
				Destination: &address,
				Value:       "pool.ntp.org:123",
				Action: func(_ context.Context, _ *cli.Command, _ string) error {
					addressPresent = true
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "hostname",
				Aliases:     []string{"h"},
				Usage:       "The hostname of the NTP server to call",
				Destination: &hostname,
				Value:       "pool.ntp.org",
				Action: func(_ context.Context, _ *cli.Command, _ string) error {
					hostnamePresent = true
					return nil
				},
			},
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "The port of the NTP server to call",
				Destination: &port,
				Value:       123,
				Action: func(_ context.Context, _ *cli.Command, _ int64) error {
					portPresent = true
					return nil
				},
			},
			&cli.BoolFlag{
				Name:        "quiet",
				Aliases:     []string{"q"},
				Usage:       "Suppress any output other than the current time or error messages",
				Destination: &quiet,
				Value:       false,
			},
		},

		Action: func(context context.Context, command *cli.Command) error {
			var currentTime *time.Time

			switch {
			case !addressPresent && hostnamePresent: // hostname is present, address is not (use default port)
				currentTime = callNTP(hostname, int(port), quiet)
			case addressPresent && !hostnamePresent && !portPresent: // address is present, hostname and port are not
				currentTime = callNTPAddress(address, quiet)
			case !addressPresent && !hostnamePresent && !portPresent: // neither address nor hostname nor port are present (use defaults)
				currentTime = callNTP(hostname, int(port), quiet)
			default:
				return errors.New("invalid arguments: you can either specify an address or a hostname, but not both")
			}

			if currentTime != nil {
				if !quiet {
					fmt.Printf("Current time: %s\n", currentTime.Format(time.UnixDate))
				} else {
					fmt.Printf("%s\n", currentTime.Format(time.UnixDate))
				}
			} else {
				return errors.New("failed to get the current time")
			}

			return nil
		},
	}

	if err := command.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	timeFormats := map[string]string{
		"Layout":      time.Layout,
		"ANSIC":       time.ANSIC,
		"UnixDate":    time.UnixDate,
		"RubyDate":    time.RubyDate,
		"RFC822":      time.RFC822,
		"RFC822Z":     time.RFC822Z,
		"RFC850":      time.RFC850,
		"RFC1123":     time.RFC1123,
		"RFC1123Z":    time.RFC1123Z,
		"RFC3339":     time.RFC3339,
		"RFC3339Nano": time.RFC3339Nano,
		"Kitchen":     time.Kitchen,
		"Stamp":       time.Stamp,
		"StampMilli":  time.StampMilli,
		"StampMicro":  time.StampMicro,
		"StampNano":   time.StampNano,
		"DateTime":    time.DateTime,
		"DateOnly":    time.DateOnly,
		"TimeOnly":    time.TimeOnly,
		"Seconds1900": "Seconds1900",
		"Seconds1970": "Seconds1970",
	}

	var address string
	addressPresent := false

	var hostname string
	hostnamePresent := false

	var port int64
	portPresent := false

	var quiet bool

	layout := time.UnixDate

	var parallels []string
	var fallbacks []string

	command := &cli.Command{
		Name:  "ntp",
		Usage: "Get the current time from an NTP server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "address",
				Aliases:     []string{"a"},
				Usage:       "The address of the NTP server",
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
				Usage:       "The hostname of the NTP server",
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
				Usage:       "The port of the NTP server",
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
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "The format to use for the output time",
				Value:   "UnixDate",
				Action: func(_ context.Context, _ *cli.Command, format string) error {
					if dateLayout := timeFormats[format]; dateLayout == "" {
						return fmt.Errorf("invalid format: %s", format)
					} else {
						layout = dateLayout
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:  "parallel",
				Usage: "Servers to call in parallel, alongside the main server. Separated by commas.",
				Action: func(_ context.Context, _ *cli.Command, serversString string) error {
					parallels = strings.Split(serversString, ",")

					return nil
				},
			},
			&cli.StringFlag{
				Name:  "fallback",
				Usage: "Fallback servers to call if the main server fails. Separated by commas.",
				Action: func(_ context.Context, _ *cli.Command, serversString string) error {
					fallbacks = strings.Split(serversString, ",")

					return nil
				},
			},
		},

		Action: func(context context.Context, command *cli.Command) error {
			var group sync.WaitGroup
			group.Add(1)

			channel := make(chan *time.Time, 1)

			switch {
			case !addressPresent && (hostnamePresent || (!hostnamePresent && !portPresent)):
				go CallNTPAsync(fmt.Sprintf("%s:%d", hostname, port), quiet, channel, &group)
			case addressPresent && !hostnamePresent && !portPresent: // address is present, hostname and port are not
				go CallNTPAsync(address, quiet, channel, &group)
			default:
				return errors.New("invalid arguments: you can either specify an address or a hostname, but not both")
			}

			if parallelsLen, fallbacksLen := len(parallels), len(fallbacks); parallelsLen > 0 && fallbacksLen > 0 {
				return errors.New("you can either specify parallel servers or fallback servers, but not both")
			} else if parallelsLen > 0 {
				for _, address := range parallels {
					go CallNTPAsync(address, quiet, channel, &group)
				}
			} else if fallbacksLen > 0 {
				group.Wait() // wait for the primary server to respond

				for _, address := range fallbacks {
					if len(channel) != 0 { // if the last server has already responded
						break
					}

					if result := CallNTP(address, quiet); result != nil {
						channel <- result
					}
				}
			}

			group.Wait()

			close(channel)

			currentTime := <-channel

			if currentTime != nil {
				timeString := FormatTime(currentTime, layout)

				if !quiet {
					fmt.Printf("Current time: ")
				}

				fmt.Printf("%s\n", timeString)

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

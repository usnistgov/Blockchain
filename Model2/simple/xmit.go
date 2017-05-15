package main

import (
  "compress/flate"
  "fmt"
  "io"
  "log"
  "strings"
  "sync"
)

// DEFLATE is suitable for transmitting compressed data across the network.
func main() {
	var wg sync.WaitGroup
	defer wg.Wait()

	// Use io.Pipe to simulate a network connection.
	// A real network application should take care to properly close the
	// underlying connection.
	rp, wp := io.Pipe()

	// Start a goroutine to act as the transmitter.
	wg.Add(1)
	go func() {
		defer wg.Done()

		zw, err := flate.NewWriter(wp, flate.BestSpeed)
		if err != nil {
			log.Fatal(err)
		}

		b := make([]byte, 256)
		for _, m := range strings.Fields("A long time ago in a galaxy far, far away...") {
			// We use a simple framing format where the first byte is the
			// message length, followed the message itself.
			b[0] = uint8(copy(b[1:], m))

			if _, err := zw.Write(b[:1+len(m)]); err != nil {
				log.Fatal(err)
			}

			// Flush ensures that the receiver can read all data sent so far.
			if err := zw.Flush(); err != nil {
				log.Fatal(err)
			}
		}

		if err := zw.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Start a goroutine to act as the receiver.
	wg.Add(1)
	go func() {
		defer wg.Done()

		zr := flate.NewReader(rp)

		b := make([]byte, 256)
		for {
			// Read the message length.
			// This is guaranteed to return for every corresponding
			// Flush and Close on the transmitter side.
			if _, err := io.ReadFull(zr, b[:1]); err != nil {
				if err == io.EOF {
					break // The transmitter closed the stream
				}
				log.Fatal(err)
			}

			// Read the message content.
			n := int(b[0])
			if _, err := io.ReadFull(zr, b[:n]); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Received %d bytes: %s\n", n, b[:n])
		}
		fmt.Println()

		if err := zr.Close(); err != nil {
			log.Fatal(err)
		}
	}()

} // end main.

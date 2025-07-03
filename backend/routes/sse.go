package routes

import (
	"bufio"
	"context"
	"io"
	"log"
	"time"

	"backend/broker"

	"github.com/gofiber/fiber/v2"
)

func SSEHandler(c *fiber.Ctx, brk *broker.Broker) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	log.Println("SSE client connected")
	fiberCtx := c.Context()

	ch := brk.RegisterClient()

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer func() {
			brk.UnregisterClient(ch)
			log.Println("SSE client disconnected")
		}()

		log.Println("Starting SSE stream")

		// Create a context that we can cancel
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Get the underlying connection
		conn := fiberCtx.Conn()
		if conn == nil {
			log.Println("No connection available")
			return
		}

		// Start a goroutine to monitor the connection
		go func() {
			buf := make([]byte, 1)
			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Try to read from the connection to detect closure
					_, err := conn.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Printf("Connection read error: %v", err)
						}
						cancel() // Cancel the main context
						return
					}
					time.Sleep(1 * time.Second)
				}
			}
		}()

		// Send initial message
		if _, err := w.WriteString(": connected\n\n"); err != nil {
			log.Printf("Error writing initial message: %v", err)
			return
		}
		if err := w.Flush(); err != nil {
			log.Printf("Error flushing initial message: %v", err)
			return
		}

		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Connection closed by monitor")
				return
			case <-fiberCtx.Done():
				log.Println("Fiber context done")
				return
			case <-ticker.C:
				if _, err := w.WriteString(": keep-alive\n\n"); err != nil {
					log.Printf("Error writing keep-alive: %v", err)
					return
				}
				if err := w.Flush(); err != nil {
					log.Printf("Error flushing keep-alive: %v", err)
					return
				}
			case msg, ok := <-ch:
				if !ok {
					log.Println("Client channel closed")
					return
				}

				log.Printf("Sending event: %s", msg)
				if _, err := w.WriteString("data: " + string(msg) + "\n\n"); err != nil {
					log.Printf("Error writing event: %v", err)
					return
				}
				if err := w.Flush(); err != nil {
					log.Printf("Error flushing event: %v", err)
					return
				}
			}
		}
	})

	return nil
}

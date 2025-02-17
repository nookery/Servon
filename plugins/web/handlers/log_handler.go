package handlers

import (
	"github.com/gin-gonic/gin"
)

// HandleStreamLogs streams logs from a specified channel using Server-Sent Events (SSE)
func (h *WebHandler) HandleStreamLogs(c *gin.Context) {
	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// Create a channel to notify if client disconnects
	clientGone := c.Writer.CloseNotify()

	// Get the log channel (you'll need to implement this based on your logging system)
	logChan := h.App.Printer.LogChan
	if logChan == nil {
		c.String(404, "Log channel not found")
		return
	}

	// Stream logs
	for {
		select {
		case <-clientGone:
			// Client disconnected
			return
		case msg, ok := <-logChan:
			if !ok {
				// Channel closed
				return
			}
			// Send log message as SSE
			c.SSEvent("log", msg)
			c.Writer.Flush()
		}
	}
}

package handler

import "fmt"

func createShortLink(host string, shorten string) string {
	return fmt.Sprintf("%s/%s", host, shorten)
}

func shortsChunk(shorts []string, batchSize int) [][]string {
	var chunks [][]string
	for {
		if len(shorts) == 0 {
			break
		}

		if len(shorts) < batchSize {
			batchSize = len(shorts)
		}

		chunks = append(chunks, shorts[0:batchSize])
		shorts = shorts[batchSize:]
	}

	return chunks
}

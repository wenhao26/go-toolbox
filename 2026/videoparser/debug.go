package main

import (
	"encoding/json"
	"fmt"
	"os"
	"toolbox/2026/videoparser/internal/ffmpeg"
)

func main() {
	version, err := ffmpeg.GetVersion()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get ffprobe version: %v\n", err)
		os.Exit(1)
	}
	
	output := map[string]string{
		"ffprobe_version": version,
		"status":          "available",
	}

	data, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(data))
}

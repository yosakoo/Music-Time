package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getDuration(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	duration := strings.TrimSpace(string(output))
	return duration, nil
}

func formatDuration(rawDuration string) (string, error) {
	seconds, err := strconv.ParseFloat(rawDuration, 64)
	if err != nil {
		return "", err
	}

	duration := time.Duration(seconds) * time.Second
	minutes := int(duration.Minutes())
	secondsRemainder := int(duration.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d", minutes, secondsRemainder), nil
}

func main() {
	file, err := os.Create("track_durations.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var totalSeconds float64

	for _, fileInfo := range files {
		if strings.HasSuffix(fileInfo.Name(), ".mp3") {
			duration, err := getDuration(fileInfo.Name())
			if err != nil {
				fmt.Printf("Error getting duration for %s: %v\n", fileInfo.Name(), err)
				continue
			}

			formattedDuration, err := formatDuration(strconv.FormatFloat(totalSeconds, 'f', -1, 64))
			if err != nil {
				fmt.Printf("Error formatting duration for %s: %v\n", fileInfo.Name(), err)
				continue
			}

			fileNameWithoutExtension := strings.TrimSuffix(fileInfo.Name(), ".mp3")

			line := fmt.Sprintf("%s - %s\n", formattedDuration, fileNameWithoutExtension)
			fmt.Print(line)
			writer.WriteString(line)
			seconds, _ := strconv.ParseFloat(duration, 64)
			totalSeconds += seconds
		}
	}

	writer.Flush()
	fmt.Println("Track durations written to track_durations.txt")
	os.Exit(0)
}

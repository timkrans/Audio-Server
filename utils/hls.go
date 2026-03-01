package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func GenerateHLS(inputFile string, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", err
	}

	m3u8File := filepath.Join(outputDir, "index.m3u8")
	//make this a an env variable eventially so path to ffmpeg
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-codec:a", "aac",
		"-b:a", "128k",
		"-hls_time", "10",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(outputDir, "segment_%03d.ts"),
		m3u8File,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %v, output: %s", err, string(output))
	}

	return m3u8File, nil
}
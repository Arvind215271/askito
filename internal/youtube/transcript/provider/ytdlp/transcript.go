package ytdlp

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	transcript "github.com/Arvind215271/askito/internal/youtube/transcript"
	parser "github.com/Arvind215271/askito/internal/youtube/transcript/parser"
)

func (c *Client) ValidateYTDLP() error {
	_, err := exec.LookPath("yt-dlp")
	return err
}

// GetTranscriptJSON fetches the pre-flattened structural JSON format from YouTube.
// This bypasses rolling text duplications and prevents post-processing conversion failures.
func (c *Client) GetTranscript(
	ctx context.Context,
	videoID string,
) (*transcript.Transcript, error) {

	videoURL := fmt.Sprintf(
		"https://www.youtube.com/watch?v=%s",
		videoID,
	)

	// Create a temporary folder to capture the download artifact
	tmpDir, err := os.MkdirTemp("", "askito-transcript-json-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// Execute the command targeting json3 specifically
	cmd := exec.CommandContext(
		ctx,
		"yt-dlp",
		"--write-auto-subs",
		"--write-subs",
		"--sub-langs", "en,en-orig",
		"--sub-format", "json3", // Target the raw layout tree directly
		"--skip-download",
		"--no-playlist",
		"--extractor-args", "youtube:player_client=android",
		"-o", filepath.Join(tmpDir, "%(id)s.%(ext)s"),
		videoURL,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf(
			"yt-dlp failed: %w: %s",
			err,
			stderr.String(),
		)
	}

	// Glob lookups search for the downloaded json3 asset
	files, err := filepath.Glob(
		filepath.Join(tmpDir, "*.json3"),
	)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("transcript JSON not available")
	}

	raw, err := os.ReadFile(files[0])
	if err != nil {
		return nil, err
	}

	// Parse and flatten the raw json3 tree down into a linear transcript string
	text, err := parser.ExtractTextFromJSON3(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json3 transcript: %w", err)
	}

	return &transcript.Transcript{
		Language: "en",
		Source:   transcript.TranscriptSourceYTDLP,
		Raw:      string(raw),
		Text:     text,
	}, nil
}



// // get the transcript using YTDLP as it is the only otpion to get those at cheap cost than youtube API quota
// func GetTranscript(
// 	ctx context.Context,
// 	videoID string,
// ) (*Transcript, error) {

// 	videoURL := fmt.Sprintf(
// 		"https://www.youtube.com/watch?v=%s",
// 		videoID,
// 	)

// 	// create a temporray folder
// 	// ytdlp gives us a .vtt file. Thus, we need to download it somewhere and then store it.
// 	// the * means any valid name that is possible ig? 
// 	tmpDir, err := os.MkdirTemp("", "askito-transcript-*")
// 	if err != nil {
// 		return nil, err
// 	}

// 	// remove this folder when the funciton exist 
// 	// as it is no longer needed to be used afterwards
// 	defer os.RemoveAll(tmpDir)

// 	// execute the command 	
// 	// auto-subs and subs both can be downlaoded.. 
// 	// here we have en.* means we only download english ones that exist here.
// 	// skip-download here skips the video from downloading and only download the subtitle here.
// 	// -o means output to be sotred in tmpDir and whatever id and extension we are using for it.
// 	// 1. The optimal yt-dlp command configuration
// 	cmd := exec.CommandContext(
// 		ctx,
// 		"yt-dlp",
// 		"--write-auto-subs",
// 		"--write-subs",
// 		"--sub-langs", "en,en-orig",
// 		"--sub-format", "vtt", // Pull the flat data structures first
// 		"--convert-subs", "vtt",           // Let yt-dlp convert it to clean VTT locally
// 		"--skip-download",
// 		"--no-playlist",
// 		"--extractor-args", "youtube:player_client=android",
// 		"-o", filepath.Join(tmpDir, "%(id)s.%(ext)s"),
// 		videoURL,
// 	)
// 	// here... this is for storing the error message
// 	// as error might occur. Why buffer? Because we are downloading it continuously... that is why
// 	var stderr bytes.Buffer
// 	cmd.Stderr = &stderr

	

// 	// this comamnd acutally runs the command. and chek for error here
// 	if err := cmd.Run(); err != nil {
// 		return nil, fmt.Errorf(
// 			"yt-dlp failed: %w: %s",
// 			err,
// 			stderr.String(),
// 		)
// 	}

// 	fmt.Println("TMP DIR:", tmpDir)

// 	filesX, _ := os.ReadDir(tmpDir)
// 	for _, f := range filesX {
// 		fmt.Println("FOUND:", f.Name())
// 	}

// 	// find the files in the path... if they exist. and if yes.. get the ones with .vtt at the  end.
// 	files, err := filepath.Glob(
// 		filepath.Join(tmpDir, "*.vtt"),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	entries, _ := os.ReadDir(tmpDir)

// 	fmt.Println("TMP DIR:", tmpDir)

// 	for _, e := range entries {
// 		fmt.Println("FILE:", e.Name())
// 	}

// 	if len(files) == 0 {
// 		return nil, fmt.Errorf(
// 			"transcript not available",
// 		)
// 	}

// 	raw, err := os.ReadFile(files[0])
// 	if err != nil {
// 		return nil, err
// 	}

	

// 	// create a text version for this file.
// 	text := ExtractCleanTranscriptTextAlreadyClean(string(raw))

// 	return &Transcript{
// 		Language: "en",
// 		Source:   TranscriptSourceYTDLP,
// 		Raw:      string(raw),
// 		Text:     text,
// 	}, nil
// }

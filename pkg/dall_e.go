package pkg

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/skratchdot/open-golang/open"
	"image/png"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

func DallE(ctx context.Context, client *openai.Client, prompt, outputFileName string) {
	// Example image as base64
	reqBase64 := openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}

	respBase64, err := client.CreateImage(ctx, reqBase64)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		fmt.Printf("Base64 decode error: %v\n", err)
		return
	}

	r := bytes.NewReader(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		fmt.Printf("PNG decode error: %v\n", err)
		return
	}

	if outputFileName == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Error getting home directory: %v", err)
		}

		re, err := regexp.Compile(`[\W]+`)
		if err != nil {
			log.Fatalf("Error cleaning prompt: %v", err)
		}

		filename := strings.TrimSuffix(re.ReplaceAllString(prompt, "-"), "-") + ".png"
		outputFileName = path.Join(homedir, "/Downloads", filename)
	}

	file, err := os.Create(outputFileName)
	if err != nil {
		fmt.Printf("File creation error: %v\n", err)
		return
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		fmt.Printf("PNG encode error: %v\n", err)
		return
	}

	fmt.Printf("The image was saved as %s\n", outputFileName)
	open.Run(outputFileName)
}

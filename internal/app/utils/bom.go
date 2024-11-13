package utils

import (
	"fmt"
	"io"
	"os"
)

func ReadFileWithUTF8BOM(filePath string) ([]byte, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Read the first 3 bytes to check for UTF-8 BOM
	bom := make([]byte, 3)
	n, err := file.Read(bom)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	// Check for UTF-8 BOM (0xEF, 0xBB, 0xBF)
	hasBOM := n >= 3 && bom[0] == 0xEF && bom[1] == 0xBB && bom[2] == 0xBF

	// Reset the file offset to after BOM if present, else to the beginning
	offset := int64(0)
	if hasBOM {
		offset = 3
	}
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("could not seek file: %w", err)
	}

	// Read the rest of the file
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file content: %w", err)
	}

	return content, nil
}

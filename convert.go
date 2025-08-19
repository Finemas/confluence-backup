package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func SavePageAsMarkdown(page *Page, dir string) error {
	md, err := ConvertToMarkdown(page.Body.Storage.Value)
	if err != nil {
		return fmt.Errorf("convert to markdown: %w", err)
	}

	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create dir: %w", err)
		}
	}

	filename := filepath.Join(dir, sanitizeFileName(page.Title)+".md")
	if err := os.WriteFile(filename, []byte(md), 0644); err != nil {
		return fmt.Errorf("save markdown: %w", err)
	}

	fmt.Printf("✅ Saved '%s' to %s\n", page.Title, filename)
	return nil
}

func ConvertToMarkdown(input string) (string, error) {
	cmd := exec.Command("pandoc", "-f", "html", "-t", "markdown")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("start pandoc: %w", err)
	}
	_, _ = io.WriteString(stdin, input)
	_ = stdin.Close()

	md, err := io.ReadAll(stdout)
	if err != nil {
		return "", fmt.Errorf("read markdown: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("pandoc wait: %w", err)
	}

	return string(md), nil
}

func sanitizeFileName(name string) string {
	invalid := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	return invalid.ReplaceAllString(name, "_")
}

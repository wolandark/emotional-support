package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Context struct {
	Program       string
	WindowTitle   string
	IsProgramming bool
	Language      string
	IsIDE         bool
	ProjectPath   string
}

type ContextDetector struct {
	programPatterns map[string]*regexp.Regexp
	languageExts    map[string][]string
}

func NewContextDetector() *ContextDetector {
	cd := &ContextDetector{
		programPatterns: make(map[string]*regexp.Regexp),
		languageExts:    make(map[string][]string),
	}

	// Define patterns for common editors/IDEs
	cd.programPatterns["vim"] = regexp.MustCompile(`(?i)vim|nvim|neovim`)
	cd.programPatterns["vscode"] = regexp.MustCompile(`(?i)code|visual studio code`)
	cd.programPatterns["emacs"] = regexp.MustCompile(`(?i)emacs`)
	cd.programPatterns["idea"] = regexp.MustCompile(`(?i)idea|intellij`)
	cd.programPatterns["sublime"] = regexp.MustCompile(`(?i)sublime`)
	cd.programPatterns["gedit"] = regexp.MustCompile(`(?i)gedit`)
	cd.programPatterns["kate"] = regexp.MustCompile(`(?i)kate`)
	cd.programPatterns["nano"] = regexp.MustCompile(`(?i)nano`)

	// Browsers
	cd.programPatterns["firefox"] = regexp.MustCompile(`(?i)firefox`)
	cd.programPatterns["chrome"] = regexp.MustCompile(`(?i)chrome`)
	cd.programPatterns["chromium"] = regexp.MustCompile(`(?i)chromium`)

	// Language file extensions
	cd.languageExts["go"] = []string{".go"}
	cd.languageExts["java"] = []string{".java"}
	cd.languageExts["python"] = []string{".py", ".pyw"}
	cd.languageExts["javascript"] = []string{".js", ".jsx", ".ts", ".tsx"}
	cd.languageExts["rust"] = []string{".rs"}
	cd.languageExts["cpp"] = []string{".cpp", ".cc", ".cxx", ".hpp", ".h", ".c"}
	cd.languageExts["ruby"] = []string{".rb"}
	cd.languageExts["php"] = []string{".php"}
	cd.languageExts["kotlin"] = []string{".kt", ".kts"}
	cd.languageExts["scala"] = []string{".scala"}
	cd.languageExts["swift"] = []string{".swift"}
	cd.languageExts["dart"] = []string{".dart"}

	return cd
}

func (cd *ContextDetector) DetectContext(windowInfo *WindowInfo) *Context {
	ctx := &Context{
		Program:       windowInfo.Process,
		WindowTitle:   windowInfo.Title,
		IsProgramming: false,
		Language:      "",
		IsIDE:         false,
	}

	// Detect program type
	titleLower := strings.ToLower(windowInfo.Title)
	processLower := strings.ToLower(windowInfo.Process)

	// Check for editors/IDEs and browsers
	programmingPrograms := map[string]bool{
		"vim": true, "vscode": true, "emacs": true, "idea": true,
		"sublime": true, "gedit": true, "kate": true, "nano": true,
	}

	for program, pattern := range cd.programPatterns {
		if pattern.MatchString(titleLower) || pattern.MatchString(processLower) {
			ctx.Program = program
			ctx.IsProgramming = programmingPrograms[program]
			ctx.IsIDE = (program == "vscode" || program == "idea" || program == "sublime")
			break
		}
	}

	// If no pattern matched but we have a process name, use it
	if ctx.Program == "" && windowInfo.Process != "" {
		ctx.Program = windowInfo.Process
	}

	// Try to extract file path from window title
	ctx.ProjectPath = cd.extractPathFromTitle(windowInfo.Title)

	// Detect language
	ctx.Language = cd.detectLanguage(windowInfo.Title, ctx.ProjectPath, ctx.Program)

	return ctx
}

func (cd *ContextDetector) extractPathFromTitle(title string) string {
	// Try to extract file path from common title formats
	// Examples: "file.py - Editor", "/path/to/file.py", "file.py (Project Name)"

	// Look for absolute paths
	if strings.HasPrefix(title, "/") {
		parts := strings.Fields(title)
		if len(parts) > 0 {
			path := parts[0]
			if strings.Contains(path, "/") {
				return filepath.Dir(path)
			}
		}
	}

	// Look for relative paths or filenames
	parts := strings.Fields(title)
	for _, part := range parts {
		if strings.Contains(part, ".") && !strings.Contains(part, " ") {
			// Might be a filename
			return filepath.Dir(part)
		}
	}

	return ""
}

func (cd *ContextDetector) detectLanguage(title, projectPath, program string) string {
	// First, try to detect from file extension in title
	titleLower := strings.ToLower(title)
	for lang, exts := range cd.languageExts {
		for _, ext := range exts {
			if strings.Contains(titleLower, ext) {
				return lang
			}
		}
	}

	// If we have a project path, try to read IDE workspace files
	if projectPath != "" {
		lang := cd.detectFromWorkspace(projectPath)
		if lang != "" {
			return lang
		}
	}

	// Try to detect from program name (e.g., "golang" in process name)
	processLower := strings.ToLower(program)
	if strings.Contains(processLower, "go") {
		return "go"
	}

	return ""
}

func (cd *ContextDetector) detectFromWorkspace(projectPath string) string {
	// Try to read VSCode workspace settings
	vsCodePath := filepath.Join(projectPath, ".vscode", "settings.json")
	if data, err := os.ReadFile(vsCodePath); err == nil {
		var settings map[string]interface{}
		if json.Unmarshal(data, &settings) == nil {
			// Check for language-specific settings
			if files, ok := settings["files.associations"].(map[string]interface{}); ok {
				for pattern, lang := range files {
					langStr := fmt.Sprintf("%v", lang)
					// Map common language identifiers
					if strings.Contains(langStr, "go") {
						return "go"
					} else if strings.Contains(langStr, "java") {
						return "java"
					} else if strings.Contains(langStr, "python") {
						return "python"
					} else if strings.Contains(langStr, "javascript") || strings.Contains(langStr, "typescript") {
						return "javascript"
					}
					_ = pattern // Use pattern if needed
				}
			}
		}
	}

	// Try to detect from common project files
	if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); err == nil {
		return "go"
	}
	if _, err := os.Stat(filepath.Join(projectPath, "package.json")); err == nil {
		return "javascript"
	}
	if _, err := os.Stat(filepath.Join(projectPath, "requirements.txt")); err == nil {
		return "python"
	}
	if _, err := os.Stat(filepath.Join(projectPath, "setup.py")); err == nil {
		return "python"
	}
	if _, err := os.Stat(filepath.Join(projectPath, "Pipfile")); err == nil {
		return "python"
	}
	if _, err := os.Stat(filepath.Join(projectPath, "pom.xml")); err == nil {
		return "java"
	}
	if _, err := os.Stat(filepath.Join(projectPath, "build.gradle")); err == nil {
		return "java"
	}
	if _, err := os.Stat(filepath.Join(projectPath, "Cargo.toml")); err == nil {
		return "rust"
	}

	return ""
}

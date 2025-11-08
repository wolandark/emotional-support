package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type MessageGenerator struct {
	rng *rand.Rand
}

func NewMessageGenerator() *MessageGenerator {
	return &MessageGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (mg *MessageGenerator) GetTimeBasedMessage(ctx *Context, duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	// Format time string
	timeStr := mg.formatDuration(hours, minutes)
	programName := mg.formatProgramName(ctx.Program)

	// Extract meaningful info from window title
	fileInfo := mg.extractFileInfo(ctx.WindowTitle)
	projectInfo := mg.extractProjectInfo(ctx.WindowTitle, ctx.ProjectPath)

	messages := []string{}

	// Program-specific messages with window title context
	if ctx.Program == "vim" || ctx.Program == "nvim" {
		if fileInfo != "" {
			messages = []string{
				fmt.Sprintf("Wow, you've been editing %s in %s for %s! I'm so proud of you! 🎉", fileInfo, programName, timeStr),
				fmt.Sprintf("%s in %s working on %s? You're a true wizard! ✨", timeStr, programName, fileInfo),
				fmt.Sprintf("Your %s skills are amazing! %s of focus on %s! 💪", programName, timeStr, fileInfo),
			}
		} else {
			messages = []string{
				fmt.Sprintf("Wow, you've been in %s for %s! I'm so proud of you! 🎉", programName, timeStr),
				fmt.Sprintf("%s in %s? You're a true wizard! ✨", timeStr, programName),
				fmt.Sprintf("Your %s skills are amazing! %s of focus! 💪", programName, timeStr),
			}
		}
	} else if ctx.Program == "vscode" {
		if projectInfo != "" {
			messages = []string{
				fmt.Sprintf("You've been coding in %s on %s for %s! Keep up the amazing work! 🚀", programName, projectInfo, timeStr),
				fmt.Sprintf("%s of dedication in %s working on %s! You're doing great! 💚", timeStr, programName, projectInfo),
				fmt.Sprintf("Look at you go! %s of focused coding in %s on %s! 🌟", timeStr, programName, projectInfo),
			}
		} else if fileInfo != "" {
			messages = []string{
				fmt.Sprintf("You've been coding in %s on %s for %s! Keep up the amazing work! 🚀", programName, fileInfo, timeStr),
				fmt.Sprintf("%s of dedication in %s! You're doing great! 💚", timeStr, programName),
				fmt.Sprintf("Look at you go! %s of focused coding in %s! 🌟", timeStr, programName),
			}
		} else {
			messages = []string{
				fmt.Sprintf("You've been coding in %s for %s! Keep up the amazing work! 🚀", programName, timeStr),
				fmt.Sprintf("%s of dedication in %s! You're doing great! 💚", timeStr, programName),
				fmt.Sprintf("Look at you go! %s of focused coding in %s! 🌟", timeStr, programName),
			}
		}
	} else if ctx.IsProgramming {
		if fileInfo != "" {
			messages = []string{
				fmt.Sprintf("You've been coding in %s on %s for %s! Keep up the amazing work! 🚀", programName, fileInfo, timeStr),
				fmt.Sprintf("%s of dedication in %s working on %s! You're doing great! 💚", timeStr, programName, fileInfo),
				fmt.Sprintf("Look at you go! %s of focused coding in %s on %s! 🌟", timeStr, programName, fileInfo),
			}
		} else {
			messages = []string{
				fmt.Sprintf("You've been coding in %s for %s! Keep up the amazing work! 🚀", programName, timeStr),
				fmt.Sprintf("%s of dedication in %s! You're doing great! 💚", timeStr, programName),
				fmt.Sprintf("Look at you go! %s of focused coding in %s! 🌟", timeStr, programName),
			}
		}
	} else {
		// Non-programming apps - use window title if available
		if ctx.Program == "firefox" || ctx.Program == "chrome" || ctx.Program == "chromium" {
			if ctx.WindowTitle != "" && len(ctx.WindowTitle) < 50 {
				messages = []string{
					fmt.Sprintf("%s is truly the best! You've been on '%s' for %s! 🌐", programName, mg.truncateTitle(ctx.WindowTitle, 40), timeStr),
					fmt.Sprintf("You've been browsing '%s' in %s for %s! Hope you're having fun! 💚", mg.truncateTitle(ctx.WindowTitle, 40), programName, timeStr),
					fmt.Sprintf("%s for %s browsing '%s'? That's some serious browsing! 🚀", programName, timeStr, mg.truncateTitle(ctx.WindowTitle, 40)),
				}
			} else {
				messages = []string{
					fmt.Sprintf("%s is truly the best! It's been %s! 🌐", programName, timeStr),
					fmt.Sprintf("You've been browsing in %s for %s! Hope you're having fun! 💚", programName, timeStr),
					fmt.Sprintf("%s for %s? That's some serious browsing! 🚀", programName, timeStr),
				}
			}
		} else if ctx.Program != "" {
			if ctx.WindowTitle != "" && len(ctx.WindowTitle) < 50 {
				messages = []string{
					fmt.Sprintf("You've been using %s working on '%s' for %s! Keep it up! 💪", programName, mg.truncateTitle(ctx.WindowTitle, 40), timeStr),
					fmt.Sprintf("%s in %s on '%s'? You're focused! 🌟", timeStr, programName, mg.truncateTitle(ctx.WindowTitle, 40)),
					fmt.Sprintf("Wow, %s in %s working on '%s'! You're doing great! 💚", timeStr, programName, mg.truncateTitle(ctx.WindowTitle, 40)),
				}
			} else {
				messages = []string{
					fmt.Sprintf("You've been using %s for %s! Keep it up! 💪", programName, timeStr),
					fmt.Sprintf("%s for %s? You're focused! 🌟", programName, timeStr),
					fmt.Sprintf("Wow, %s in %s! You're doing great! 💚", timeStr, programName),
				}
			}
		}
	}

	if len(messages) > 0 {
		return messages[mg.rng.Intn(len(messages))]
	}
	return ""
}

func (mg *MessageGenerator) extractFileInfo(windowTitle string) string {
	if windowTitle == "" {
		return ""
	}

	// Try to extract filename from common patterns
	// Examples: "main.go - Editor", "/path/to/file.py", "file.py (Project Name)"

	// Look for file extensions
	extensions := []string{".go", ".py", ".js", ".ts", ".java", ".rs", ".cpp", ".c", ".h", ".rb", ".php", ".kt", ".swift", ".dart", ".scala"}
	for _, ext := range extensions {
		if idx := strings.Index(windowTitle, ext); idx > 0 {
			// Extract filename
			start := strings.LastIndex(windowTitle[:idx], "/")
			if start == -1 {
				start = strings.LastIndex(windowTitle[:idx], " ")
			}
			if start >= 0 {
				filename := strings.TrimSpace(windowTitle[start+1 : idx+len(ext)])
				if len(filename) > 0 && len(filename) < 50 {
					return filename
				}
			} else {
				filename := strings.TrimSpace(windowTitle[:idx+len(ext)])
				if len(filename) > 0 && len(filename) < 50 {
					return filename
				}
			}
		}
	}

	return ""
}

func (mg *MessageGenerator) extractProjectInfo(windowTitle, projectPath string) string {
	// First try project path
	if projectPath != "" {
		parts := strings.Split(projectPath, "/")
		if len(parts) > 0 {
			projectName := parts[len(parts)-1]
			if projectName != "" && len(projectName) < 40 {
				return projectName
			}
		}
	}

	// Try to extract from window title (e.g., "file.py (Project Name)")
	if windowTitle != "" {
		if idx := strings.Index(windowTitle, "("); idx > 0 {
			if endIdx := strings.Index(windowTitle[idx:], ")"); endIdx > 0 {
				projectName := strings.TrimSpace(windowTitle[idx+1 : idx+endIdx])
				if len(projectName) > 0 && len(projectName) < 40 {
					return projectName
				}
			}
		}
	}

	return ""
}

func (mg *MessageGenerator) truncateTitle(title string, maxLen int) string {
	if len(title) <= maxLen {
		return title
	}
	return title[:maxLen-3] + "..."
}

func (mg *MessageGenerator) formatDuration(hours, minutes int) string {
	if hours > 0 && minutes > 0 {
		return fmt.Sprintf("%d hour%s and %d minute%s", hours, plural(hours), minutes, plural(minutes))
	} else if hours > 0 {
		return fmt.Sprintf("%d hour%s", hours, plural(hours))
	} else {
		return fmt.Sprintf("%d minute%s", minutes, plural(minutes))
	}
}

func (mg *MessageGenerator) formatProgramName(program string) string {
	// Capitalize and format program names nicely
	if program == "" {
		return "this app"
	}

	// Handle common program names
	names := map[string]string{
		"vim":      "Vim",
		"nvim":     "Neovim",
		"vscode":   "VS Code",
		"emacs":    "Emacs",
		"idea":     "IntelliJ IDEA",
		"sublime":  "Sublime Text",
		"firefox":  "Firefox",
		"chrome":   "Chrome",
		"chromium": "Chromium",
		"gedit":    "gedit",
		"kate":     "Kate",
		"nano":     "Nano",
	}

	if niceName, ok := names[program]; ok {
		return niceName
	}

	// Capitalize first letter
	if len(program) > 0 {
		return strings.ToUpper(program[:1]) + program[1:]
	}
	return program
}

func (mg *MessageGenerator) GetLanguageMessage(language string) string {
	messages := map[string][]string{
		"java": {
			"I know Java is hard, but you got it! 💪",
			"Java can be tricky, but you're handling it like a pro! 🌟",
			"Keep pushing through those Java challenges! You're doing great! 💚",
		},
		"cpp": {
			"C++ is complex, but you're tackling it! Keep going! 🚀",
			"Memory management is tough, but you've got this! 💪",
			"You're doing amazing work with C++! 🌟",
		},
		"rust": {
			"Rust's borrow checker can be challenging, but you're learning! 💚",
			"Keep fighting the good fight with Rust! You're awesome! 🦀",
			"Rust is hard, but you're making progress! Keep it up! ✨",
		},
		"go": {
			"Go is a great choice! You're doing fantastic! 🐹",
			"Keep up the great work with Go! 💪",
			"Your Go code is going to be amazing! 🌟",
		},
		"python": {
			"Python is fun! Keep enjoying the journey! 🐍",
			"You're doing great with Python! 💚",
			"Keep up the awesome Python work! ✨",
		},
		"javascript": {
			"JavaScript can be wild, but you're taming it! 🚀",
			"Keep up the great work with JavaScript/TypeScript! 💪",
			"You're doing amazing with JS/TS! 🌟",
		},
	}

	if langMsgs, ok := messages[language]; ok {
		return langMsgs[mg.rng.Intn(len(langMsgs))]
	}

	// Generic programming message
	generic := []string{
		fmt.Sprintf("You're doing great with %s! Keep it up! 💚", language),
		fmt.Sprintf("Keep pushing forward with %s! You've got this! 💪", language),
	}
	return generic[mg.rng.Intn(len(generic))]
}

func (mg *MessageGenerator) GetHealthReminder() string {
	messages := []string{
		"💧 Remember to stay hydrated! Take a sip of water!",
		"👀 Blink your eyes! Give them a break from the screen!",
		"💚 Take a deep breath! You're doing great!",
		"🪑 Stretch a bit! Your body will thank you!",
		"☕ Time for a quick break? Maybe some water or tea?",
		"👁️ Look away from the screen for 20 seconds! Your eyes need it!",
		"🧘 Take a moment to relax your shoulders!",
		"💧 Hydration check! Have you had water recently?",
	}
	return messages[mg.rng.Intn(len(messages))]
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

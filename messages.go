package main

import (
	"fmt"
	"math/rand"
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

	messages := []string{}

	if ctx.Program == "vim" {
		messages = []string{
			fmt.Sprintf("Wow, you've been in vim for %d hour%s! I'm so proud of you! 🎉", hours, plural(hours)),
			fmt.Sprintf("%d hour%s in vim? You're a true wizard! ✨", hours, plural(hours)),
			fmt.Sprintf("Your vim skills are amazing! %d hour%s of focus! 💪", hours, plural(hours)),
		}
	} else if ctx.IsProgramming {
		messages = []string{
			fmt.Sprintf("You've been coding for %d hour%s! Keep up the amazing work! 🚀", hours, plural(hours)),
			fmt.Sprintf("%d hour%s of dedication! You're doing great! 💚", hours, plural(hours)),
			fmt.Sprintf("Look at you go! %d hour%s of focused coding! 🌟", hours, plural(hours)),
		}
	}

	if len(messages) > 0 {
		return messages[mg.rng.Intn(len(messages))]
	}
	return ""
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

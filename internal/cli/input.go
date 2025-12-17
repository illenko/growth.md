package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func PromptString(prompt string, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultValue
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}

	return input
}

func PromptStringRequired(prompt string) string {
	for {
		input := PromptString(prompt, "")
		if input != "" {
			return input
		}
		PrintError(fmt.Errorf("this field is required"))
	}
}

func PromptConfirm(prompt string) bool {
	input := PromptString(prompt+" (y/n)", "")
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes"
}

func PromptConfirmDefault(prompt string, defaultYes bool) bool {
	defaultStr := "n"
	if defaultYes {
		defaultStr = "y"
	}

	fmt.Printf("%s (y/n) [%s]: ", prompt, defaultStr)
	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultYes
	}

	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" {
		return defaultYes
	}

	return input == "y" || input == "yes"
}

func PromptSelect(prompt string, options []string) string {
	if len(options) == 0 {
		return ""
	}

	fmt.Println(prompt)
	for i, option := range options {
		fmt.Printf("  %d) %s\n", i+1, option)
	}
	fmt.Print("Select (1-", len(options), "): ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return options[0]
	}

	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(options) {
		PrintWarning("Invalid selection, using first option")
		return options[0]
	}

	return options[choice-1]
}

func PromptSelectWithDefault(prompt string, options []string, defaultOption string) string {
	if len(options) == 0 {
		return defaultOption
	}

	defaultIdx := -1
	for i, opt := range options {
		if opt == defaultOption {
			defaultIdx = i
			break
		}
	}

	fmt.Println(prompt)
	for i, option := range options {
		marker := " "
		if i == defaultIdx {
			marker = "*"
		}
		fmt.Printf("  %s %d) %s\n", marker, i+1, option)
	}

	defaultDisplay := ""
	if defaultIdx >= 0 {
		defaultDisplay = fmt.Sprintf(" [%d]", defaultIdx+1)
	}
	fmt.Printf("Select (1-%d)%s: ", len(options), defaultDisplay)

	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultOption
	}

	input = strings.TrimSpace(input)
	if input == "" && defaultIdx >= 0 {
		return defaultOption
	}

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(options) {
		if defaultIdx >= 0 {
			PrintWarning("Invalid selection, using default option")
			return defaultOption
		}
		PrintWarning("Invalid selection, using first option")
		return options[0]
	}

	return options[choice-1]
}

func PromptInt(prompt string, defaultValue int) int {
	defaultStr := ""
	if defaultValue != 0 {
		defaultStr = fmt.Sprintf("%d", defaultValue)
	}

	input := PromptString(prompt, defaultStr)
	value, err := strconv.Atoi(input)
	if err != nil {
		return defaultValue
	}

	return value
}

func PromptMultiline(prompt string) string {
	fmt.Println(prompt)
	fmt.Println("(Press Ctrl+D or enter a line with just '.' to finish)")

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "." {
			break
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func ValidateNotEmpty(value string, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func ValidateOneOf(value string, options []string) error {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return nil
	}

	for _, option := range options {
		if strings.ToLower(option) == value {
			return nil
		}
	}

	return fmt.Errorf("must be one of: %s", strings.Join(options, ", "))
}

func ValidatePositive(value int, fieldName string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive", fieldName)
	}
	return nil
}

func ValidateRange(value int, min, max int, fieldName string) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be between %d and %d", fieldName, min, max)
	}
	return nil
}

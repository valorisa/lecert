package wizard

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Answers struct {
	Domain    string
	Email     string
	Challenge string
}

func Run() (*Answers, error) {
	reader := bufio.NewReader(os.Stdin)
	answers := &Answers{}

	fmt.Println("=== Let's Encrypt Certificate Wizard ===")
	fmt.Println()

	// Question 1: Domain
	fmt.Print("1/3 Domain name (e.g. example.com): ")
	domain, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("read domain: %w", err)
	}
	answers.Domain = strings.TrimSpace(domain)
	if answers.Domain == "" {
		return nil, fmt.Errorf("domain cannot be empty")
	}

	// Question 2: Email
	fmt.Print("2/3 Your email (for certificate expiry alerts): ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("read email: %w", err)
	}
	answers.Email = strings.TrimSpace(email)
	if answers.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	// Question 3: Challenge type (simplified)
	fmt.Println("3/3 How should we verify domain ownership?")
	fmt.Println("    [1] HTTP (needs port 80 open) — recommended for most setups")
	fmt.Println("    [2] DNS  (needs DNS provider API access)")
	fmt.Print("    Choice [1]: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	switch choice {
	case "2":
		answers.Challenge = "dns-01"
	default:
		answers.Challenge = "http-01"
	}

	fmt.Println()
	fmt.Printf("Got it! Requesting certificate for %s via %s challenge...\n", answers.Domain, answers.Challenge)
	return answers, nil
}

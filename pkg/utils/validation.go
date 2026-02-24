package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var (
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
	ErrPasswordNoUppercase = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLowercase = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoDigit     = errors.New("password must contain at least one digit")
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrInvalidRepoURL      = errors.New("invalid repository URL")
)

// ValidatePassword checks if a password meets security requirements
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUppercase
	}
	if !hasLower {
		return ErrPasswordNoLowercase
	}
	if !hasDigit {
		return ErrPasswordNoDigit
	}

	return nil
}

// ValidateEmail checks if an email address is valid
func ValidateEmail(email string) error {
	// Basic email regex pattern
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidateRepoURL validates a Git repository URL (GitHub, GitLab, Bitbucket, or generic)
func ValidateRepoURL(url string) error {
	// Accept HTTPS URLs from major Git providers and generic Git servers
	httpsPatterns := []*regexp.Regexp{
		// GitHub: https://github.com/user/repo.git or https://github.com/user/repo
		regexp.MustCompile(`^https://github\.com/[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+(\.git)?$`),
		// GitLab: https://gitlab.com/user/repo.git or https://gitlab.com/user/repo
		regexp.MustCompile(`^https://gitlab\.com/[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+(\.git)?$`),
		// Bitbucket: https://bitbucket.org/user/repo.git or https://bitbucket.org/user/repo
		regexp.MustCompile(`^https://bitbucket\.org/[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+(\.git)?$`),
		// Generic HTTPS Git URL
		regexp.MustCompile(`^https?://[a-zA-Z0-9_.-]+\.[a-zA-Z]{2,}/[a-zA-Z0-9_.-/]+(\.git)?$`),
	}

	// Accept SSH URLs from major Git providers and generic Git servers
	sshPatterns := []*regexp.Regexp{
		// GitHub SSH: git@github.com:user/repo.git
		regexp.MustCompile(`^git@github\.com:[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+(\.git)?$`),
		// GitLab SSH: git@gitlab.com:user/repo.git
		regexp.MustCompile(`^git@gitlab\.com:[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+(\.git)?$`),
		// Bitbucket SSH: git@bitbucket.org:user/repo.git
		regexp.MustCompile(`^git@bitbucket\.org:[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+(\.git)?$`),
		// Generic SSH Git URL
		regexp.MustCompile(`^git@[a-zA-Z0-9_.-]+\.[a-zA-Z]{2,}:[a-zA-Z0-9_.-/]+(\.git)?$`),
	}

	// Check HTTPS patterns
	for _, pattern := range httpsPatterns {
		if pattern.MatchString(url) {
			return nil
		}
	}

	// Check SSH patterns
	for _, pattern := range sshPatterns {
		if pattern.MatchString(url) {
			return nil
		}
	}

	return ErrInvalidRepoURL
}

// SanitizeRepoName extracts and sanitizes repository name from URL
func SanitizeRepoName(repoURL string) string {
	// Remove .git suffix
	url := strings.TrimSuffix(repoURL, ".git")

	// Get the last part of the URL
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return "project"
	}

	name := parts[len(parts)-1]

	// Only allow alphanumeric, dash, underscore
	var sanitized strings.Builder
	for _, char := range name {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '-' || char == '_' {
			sanitized.WriteRune(char)
		}
	}

	result := sanitized.String()
	if result == "" {
		return "project"
	}

	return result
}

// SanitizeDeployScript validates and sanitizes deploy script
// Returns error if dangerous patterns are detected
func ValidateDeployScript(script string) error {
	// Check for dangerous patterns
	dangerousPatterns := []string{
		"rm -rf /",
		"rm -rf /*",
		"; rm ",
		"&& rm -rf",
		"| rm ",
		"`rm ",
		"$(rm ",
		"curl | bash",
		"wget | bash",
		"curl | sh",
		"wget | sh",
	}

	lowerScript := strings.ToLower(script)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerScript, strings.ToLower(pattern)) {
			return errors.New("deploy script contains potentially dangerous commands")
		}
	}

	return nil
}

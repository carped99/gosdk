package outbox

import (
	"fmt"
	"regexp"
	"strings"
)

// Security utilities for preventing SQL injection and other attacks

var (
	// SQL identifier validation regex
	// Allows: letters, numbers, underscores, dots (for schema.table format)
	// Must start with a letter
	sqlIdentifierRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*(\.[a-zA-Z][a-zA-Z0-9_]*)*$`)

	// Blacklisted SQL keywords that could be used in injection attacks
	sqlKeywords = map[string]bool{
		"select": true, "insert": true, "update": true, "delete": true,
		"drop": true, "create": true, "alter": true, "truncate": true,
		"union": true, "exec": true, "execute": true, "script": true,
		"javascript": true, "vbscript": true, "expression": true,
		"declare": true, "begin": true, "end": true, "if": true,
		"while": true, "for": true, "loop": true, "goto": true,
		"waitfor": true, "delay": true, "sleep": true, "benchmark": true,
		"load_file": true, "into": true, "outfile": true, "dumpfile": true,
		"information_schema": true, "sys": true, "master": true,
		"xp_": true, "sp_": true, "fn_": true,
	}

	// Dangerous characters that could be used in SQL injection
	dangerousChars = []string{
		"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_", "fn_",
		"0x", "0b", "\\", "`", "=", "<", ">", "!", "&", "|",
	}
)

// ValidateSQLIdentifier validates a SQL identifier (table name, column name, etc.)
// to prevent SQL injection attacks
func ValidateSQLIdentifier(identifier string) error {
	if identifier == "" {
		return fmt.Errorf("identifier cannot be empty")
	}

	// Check length limits
	if len(identifier) > 128 {
		return fmt.Errorf("identifier too long (max 128 characters): %s", identifier)
	}

	// Check for dangerous characters
	for _, char := range dangerousChars {
		if strings.Contains(identifier, char) {
			return fmt.Errorf("identifier contains dangerous character '%s': %s", char, identifier)
		}
	}

	// Check for SQL keywords
	lowerIdentifier := strings.ToLower(identifier)
	if sqlKeywords[lowerIdentifier] {
		return fmt.Errorf("identifier is a reserved SQL keyword: %s", identifier)
	}

	// Check for keyword patterns
	for keyword := range sqlKeywords {
		if strings.Contains(lowerIdentifier, keyword) {
			return fmt.Errorf("identifier contains SQL keyword '%s': %s", keyword, identifier)
		}
	}

	// Validate format using regex
	if !sqlIdentifierRegex.MatchString(identifier) {
		return fmt.Errorf("invalid identifier format: %s (must start with a letter and contain only alphanumeric characters, underscores, and dots)", identifier)
	}

	// Additional checks for common injection patterns
	if strings.Contains(strings.ToLower(identifier), "union") {
		return fmt.Errorf("identifier contains UNION keyword: %s", identifier)
	}

	if strings.Contains(strings.ToLower(identifier), "select") {
		return fmt.Errorf("identifier contains SELECT keyword: %s", identifier)
	}

	return nil
}

// sanitizeTableName sanitizes and validates a table name
func sanitizeTableName(tableName string) (string, error) {
	// Trim whitespace
	tableName = strings.TrimSpace(tableName)

	// Validate the table name
	if err := ValidateSQLIdentifier(tableName); err != nil {
		return "", fmt.Errorf("invalid table name: %w", err)
	}

	// Normalize the table name (convert to lowercase for consistency)
	// Note: Some databases are case-sensitive, so we keep the original case
	// but normalize for comparison
	normalized := strings.ToLower(tableName)

	// Additional checks for table-specific patterns
	if strings.Contains(normalized, "information_schema") {
		return "", fmt.Errorf("table name cannot reference information_schema: %s", tableName)
	}

	if strings.Contains(normalized, "pg_") || strings.Contains(normalized, "mysql.") {
		return "", fmt.Errorf("table name cannot reference system tables: %s", tableName)
	}

	return tableName, nil
}

// IsSQLInjectionAttempt checks if a string contains potential SQL injection patterns
func IsSQLInjectionAttempt(input string) bool {
	if input == "" {
		return false
	}

	lowerInput := strings.ToLower(input)

	// Check for common SQL injection patterns
	injectionPatterns := []string{
		"union select", "union all select", "union distinct select",
		"'; drop table", "'; delete from", "'; truncate table",
		"'; insert into", "'; update ", "'; alter table",
		"exec(", "execute(", "sp_", "xp_", "fn_",
		"waitfor delay", "benchmark(", "sleep(",
		"load_file(", "into outfile", "into dumpfile",
		"information_schema", "sys.tables", "sys.columns",
		"0x", "0b", "\\x", "\\u", "\\n", "\\r", "\\t",
		"<!--", "-->", "<script", "</script>",
		"javascript:", "vbscript:", "expression(",
	}

	for _, pattern := range injectionPatterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}

	// Check for balanced quotes (potential for breaking out of strings)
	singleQuotes := strings.Count(input, "'")
	doubleQuotes := strings.Count(input, "\"")

	if singleQuotes%2 != 0 || doubleQuotes%2 != 0 {
		return true
	}

	// Check for comment patterns
	if strings.Contains(lowerInput, "--") || strings.Contains(lowerInput, "/*") {
		return true
	}

	return false
}

package main

import (
	"net"
	"strings"
	"unicode/utf8"
)

// Based on http://blog.onyxbits.de/validating-email-addresses-with-a-regex-do-yourself-a-favor-and-dont-391/

func IsValid(input string) bool {
	// RFC 5321 limits the total length of an address
	if len(input) == 0 || len(input) > 254 {
		return false
	}
	var ch rune
	var local []rune
	var domain []rune
	state := 0
	index := 0
	size := 0
	runeCount := utf8.RuneCountInString(input)
	for len(input) > 0 && state != -1 {
		ch, size = utf8.DecodeRuneInString(input)

		switch state {
		case 0: // local part start
			if ch == '(' {
				consumeComment(&input, &index, &size, runeCount, &ch)
				break
			}
			if ch == '"' {
				local = append(local, ch)
				state = 5
				break
			}
			if isAllowedLocal(ch) {
				local = append(local, ch)
				state = 1
				break
			}
			state = -1
		case 1: // local part
			if ch == '(' {
				consumeComment(&input, &index, &size, runeCount, &ch)
				break
			}
			if isAllowedLocal(ch) {
				local = append(local, ch)
				break
			}
			if ch == '.' {
				local = append(local, ch)
				state = 0
				break
			}
			if ch == '@' { // Endof local part
				state = 2
				break
			}
			state = -1
		case 2: // domain part start
			if ch == '(' {
				consumeComment(&input, &index, &size, runeCount, &ch)
				break
			}
			if ch == '[' {
				ip := readIPAddress(&input, &index, &size, runeCount, &ch)
				if strings.Contains(ip, ":") {
					// we have IPv6 address
					if !strings.HasPrefix(ip, "IPv6:") {
						// does not include IPv6: prefix
						state = -1
						break
					}
					ip = ip[5:]
				}
				if net.ParseIP(ip) != nil {
					state = 8
					break
				} else {
					// invalid ip address
					state = -1
					break
				}
			}
			if isAllowedDomain(ch) {
				domain = append(domain, ch)
				state = 3
				break
			}
			state = -1
		case 3: //domain part
			if ch == '(' {
				consumeComment(&input, &index, &size, runeCount, &ch)
				break
			}
			if isAllowedDomain(ch) {
				domain = append(domain, ch)
				break
			}
			if ch == '-' || ch == '.' {
				domain = append(domain, ch)
				state = 4
				break
			}
			state = -1
		case 4: // domain part
			if isAllowedDomain(ch) {
				domain = append(domain, ch)
				state = 3
				break
			}
			if ch == '-' {
				domain = append(domain, ch)
				break
			}
			state = -1
		case 5: // in quoted string
			if ch == '\\' {
				local = append(local, ch)
				state = 7
				break
			}
			if ch == '"' {
				local = append(local, ch)
				state = 6
				break
			}
			local = append(local, ch)
		case 6: // after double quote in quoted part that is not perceded by backslash
			if ch == '.' {
				local = append(local, ch)
				state = 0
				break
			}
			if ch == '@' { // End of local part
				state = 2
				break
			}
			state = -1
		case 7: // in quoted string preceded by backslash
			if ch == '\\' || ch == '"' {
				local = append(local, ch)
				state = 5
				break
			}
			state = -1
		case 8: // domain part end
			// We're here because ip was good. Should not have any characters anymore exept comments
			if ch == '(' {
				consumeComment(&input, &index, &size, runeCount, &ch)
				break
			}
			state = -1
		}
		input = input[size:]
		index++
	}
	if state != 3 && state != 8 {
		return false
	}
	if len(local) > 64 { // local part limited to 64 code points
		return false
	}
	return true
}

func isAllowedLocal(r rune) bool {
	return ((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' || r == '+' || r == '!' || r == '#' || r == '$' || r == '%' || r == '&' || r == '\'' || r == '*' || r == '/' || r == '=' || r == '?' || r == '^' || r == '_' || r == '`' || r == '{' || r == '|' || r == '}' || r == '~' || r > 127)
}

func isAllowedDomain(r rune) bool {
	return ((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'))
}

func isAllowedInQuoted(r rune) bool {
	return (r == ' ' || r == '"' || r == '(' || r == ')' || r == ',' || r == ':' || r == ';' || r == '<' || r == '>' || r == '@' || r == '[' || r == '\\' || r == ']')
}

func consumeComment(input *string, index, size *int, runeCount int, ch *rune) {
	for *ch != ')' && *index < runeCount {
		*input = (*input)[*size:]
		*index++
		*ch, *size = utf8.DecodeRuneInString(*input)
	}
}

func readIPAddress(input *string, index, size *int, runeCount int, ch *rune) string {
	var ip []rune
	for *ch != ']' && *index < runeCount {
		*input = (*input)[*size:]
		*index++
		*ch, *size = utf8.DecodeRuneInString(*input)
		ip = append(ip, *ch)
	}
	return string(ip[:len(ip)-1])
}

package regexp

import (
	"regexp"
)

// Verify mobile phone, most 11 byte length.
func MobilePhone(str ...string) bool {
	var yes bool
	for _, s := range str {
		yes, _ = regexp.MatchString("^1[0-9]{10}$", s)
		if !yes {
			return yes
		}
	}
	return yes
}

// Verify telephone, most 8 byte length.
func TellPhone(str ...string) bool {
	var yes bool
	for _, s := range str {
		yes, _ = regexp.MatchString("^[0-9]{8}$", s)
		if !yes {
			return yes
		}
	}
	return yes
}

// Verify email
func Email(str ...string) bool {
	var yes bool
	for _, s := range str {
		yes, _ = regexp.MatchString("^([a-z0-9_\\.-]+)@([\\da-z\\.-]+)\\.([a-z\\.]{2,6})$", s)
		if !yes {
			return yes
		}
	}
	return yes
}

// Find bulk chinese word.
func Chinese(str ... string) bool {
	var yes bool
	for _, v := range str {
		yes, _ = regexp.MatchString(`[\p{Han}]+`, v)
		if !yes {
			return yes
		}
	}
	return yes
}

// Find special character.
func SpecialChar(str ... string) bool {
	var yes bool
	for _, s := range str {
		yes, _ = regexp.MatchString(`[\f\t\n\r\v\123\x7F\x{10FFFF}\\\^\$\.\*\+\?\{\}\(\)\[\]|]`, s)
		if !yes {
			return yes
		}
	}
	return yes
}

// Find html tag.
func HtmlTag(str ... string) bool {
	var yes bool
	for _, s := range str {
		yes, _ = regexp.MatchString(`<[^>]+>`, s)
		if !yes {
			return yes
		}
	}
	return yes
}

// Find html Script.
func HtmlScript(str ... string) bool {
	var yes bool
	for _, s := range str {
		yes, _ = regexp.MatchString(`<script[^>]*?>[\\s\\S]*?<\\/script>`, s)
		if !yes {
			return yes
		}
	}
	return yes
}

// Find css style.
func HtmlStyle(str ... string) bool {
	var yes bool
	for _, s := range str {
		yes, _ = regexp.MatchString(`<style[^>]*?>[\\s\\S]*?<\\/style>`, s)
		if !yes {
			return yes
		}
	}
	return yes
}

// Find space, tap, newline
func ESChar(str ... string) bool {
	var yes bool
	for _, s := range str {
		yes, _ = regexp.MatchString(`\\s*|\t|\r|\n`, s)
		if !yes {
			return yes
		}
	}
	return yes
}

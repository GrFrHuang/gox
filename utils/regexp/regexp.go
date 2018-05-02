package regexp

import (
	"regexp"
	"strings"
	"strconv"
	"math"
)

// Verify mobile phone, most 11 byte length.
func MobilePhone(str ...string) bool {
	var yes bool
	for _, s := range str {
		yes, _ = regexp.MatchString(`^(13[0-9]|15[0-9]|18[0-9]|14[0-9]|17[0-9])\d{8}$`, s)
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

// Find space, tap, newline.
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

// Verify password special.
func Password(str ...string) bool {
	var yes bool
	for _, s := range str {
		yes, _ = regexp.MatchString(`([A-Z]|[a-z]|[0-9]|[\~!@#$%^&*()+=|{}':;',\\\\[\\\\].<>/?~！@#￥%……&*（）——+|{}【】‘；：”“'。，、？]){6,20}$`, s)
		if !yes {
			return yes
		}
	}
	return yes
}

// Verify person id card number.
func VerifyIdCard(idStr string) (bool) {
	var value [18]float64
	var i int = 0
	//权值
	var verify_num = [17]float64{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	//计算出的检验位
	var verify_code string
	//加权求和的值
	var s float64
	//是否检验通过
	var end_value string

	IdSlice := strings.SplitN(idStr, "", 18)
	//转换，用户身份证的最后一位
	for k, v := range IdSlice {
		value[k], _ = strconv.ParseFloat(v, 1)
		end_value = v

	}
	//计算权值和
	for i < 17 {
		s = s + (value[i] * verify_num[i])
		i++
	}
	sum := math.Mod(s, 11)
	//value := strconv.FormatFloat(y,1,1,1)
	//no break?
	switch sum {
	case 0:
		verify_code = "1"
	case 1:
		verify_code = "0"
	case 2:
		verify_code = "X"
	case 3:
		verify_code = "9"
	case 4:
		verify_code = "8"
	case 5:
		verify_code = "7"
	case 6:
		verify_code = "6"
	case 7:
		verify_code = "5"
	case 8:
		verify_code = "4"
	case 9:
		verify_code = "3"
	case 10:
		verify_code = "2"
	default:
		return false
	}
	if verify_code != end_value {
		return false
	} else {
		return true
	}
}

// Step.1 -- Prevent sql fight.
// Step.2 -- Variable type conversion.
// Step.3 -- Don't return the sql native error.
// Step.4 -- Give the DB user minimal executive authority.
func SqlFilter() {

}

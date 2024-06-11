package util

import (
	"regexp"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type CompanyDetails struct {
	Name      string `json:"name"`
	AllRights string `json:"all_rights"`
}

// HttpWebResponseSuccess :
func HttpWebResponseSuccess(code int, message string) gin.H {
	response := gin.H{
		"all_rights": "© 2024 AI Resume Automation. All Rights Reserved.",
		"status":     "success",
		"code":       code,
		"error":      nil,
		"message":    message,
	}

	return response
}

// HttpWebResponseWithDataSuccess :
func HttpWebResponseWithDataSuccess(code int, message string, data interface{}) gin.H {
	response := gin.H{
		"all_rights": "© 2024 AI Resume Automation. All Rights Reserved.",
		"status":     "success",
		"code":       code,
		"error":      nil,
		"message":    message,
		"data":       data,
	}

	return response
}

// HttpWebResponseError :
func HttpWebResponseError(code int, errors, message string) gin.H {
	response := gin.H{
		"all_rights": "© 2024 AI Resume Automation. All Rights Reserved.",
		"status":     "error",
		"code":       code,
		"error":      errors,
		"message":    message,
	}

	return response
}

// CheckPathExists :
func CheckPathExists(dirDump string, createDirDump int) (bool, string) {
	var err error
	if _, err = os.Stat(dirDump); err != nil {
		if os.IsNotExist(err) {
			if createDirDump == 1 {
				err = os.MkdirAll(dirDump, 0700)
				if err != nil {
					errmsg := fmt.Sprintf("Err: Path to backup \"%s\". %s", dirDump, err.Error())
					return false, errmsg
				} else {
					return true, ""
				}
			} else {
				errmsg := fmt.Sprintf("Err: Path to backup \"%s\" doesn't exists", dirDump)
				return false, errmsg
			}
		} else {
			errmsg := fmt.Sprintf("Err: Path to backup \"%s\" doesn't exists", dirDump)
			return false, errmsg
		}
	} else {
		return true, ""
	}
}

// ExtractPercentage :
func ExtractPercentage(str string) (string, error) {
    // Regular expression to match a percentage value
    re := regexp.MustCompile(`(\d+(\.\d+)?)%`)

    // Find the first match in the string
    match := re.FindStringSubmatch(str)
    if len(match) < 2 {
        return "0", fmt.Errorf("percentage not found in the string")
    }

    // Extract the percentage value from the matched string
    percentageStr := match[1]

    // Parse the extracted percentage value to float64
    // percentage, err := strconv.ParseFloat(percentageStr, 64)
    // if err != nil {
    //     return 0, fmt.Errorf("failed to parse percentage value: %v", err)
    // }

    return percentageStr, nil
}
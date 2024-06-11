package parser

import (
	"bytes"
	"fmt"
	"path/filepath"

	"users/logger"

	"github.com/dslipak/pdf"
)

// ReadFiles :
func ReadFiles(filelocation string) (string, error) {

		fileExtension := filepath.Ext(filelocation)

		fmt.Println("file name ===", filelocation, fileExtension)
		// var err error

		switch fileExtension {
		case ".pdf":

			pdfPath := filelocation

			fmt.Println("this pdf in plain text :")
			// pdf.DebugOn = true
			content, err := readPdfPlainText(pdfPath)
			if err != nil {
				fmt.Println("error in read pdf:", err)
				return "", err
			} 
			return content, err

		case ".doc", ".docx":

			// _, err := readDocxPlainText(filelocation)
			// if err != nil {
			// 	fmt.Println("error in "+fileExtension + ":", err)
			// 	return "", err
			// }

		default:
			logger.Log.Error().Err(nil).Msg("parser: invalid file format, ProcessUploadedUserData()")

		}


		return "", nil
}

// readPdfPlainText :
func readPdfPlainText(path string) (string, error) {
	r, err := pdf.Open(path)
	// remember close file
	// defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

/*
// readDocxPlainText :
func readDocxPlainText(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	
	var r io.Reader
	r = f
	
	tmpl, _, err := docconv.ConvertDocx(r)
	if err != nil {
		return "", err
	}

	fmt.Println(tmpl)

	return tmpl, nil
}
*/
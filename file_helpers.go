package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

func getMapContent(rawContent string) map[string]string {
	var item = make(map[string]string)
	contentLines := strings.Split(rawContent, "\n")

	fields := [15]string{
		"Message-ID",
		"Date",
		"From",
		"To",
		"Subject",
		"Cc",
		"Mime-Version",
		"Content-Type",
		"Content-Transfer-Encoding",
		"X-From",
		"X-To",
		"X-cc",
		"X-Folder",
		"X-Origin",
		"X-FileName",
	}
	for _, v := range fields {
		lineInfo := Filter(contentLines, func(i string) bool {
			return strings.HasPrefix(i, v)
		})

		if len(lineInfo) > 0 {
			splitedline := strings.Split(lineInfo[0], ": ")

			if len(splitedline) > 2 {
				splitedline[1] = strings.Join(splitedline[1:], ": ")
			}

			if len(splitedline) > 1 {
				switch v {
				case "To": // obtains all "To" emails when are in more than one line
					to := splitedline[1]
					index := 3
					for {
						nextline := contentLines[index+1]
						if !strings.HasPrefix(nextline, "\t") {
							break
						} else {
							to = to + strings.Replace(nextline, "\t", "", -1)
						}
						index++
					}
					item[v] = to
				case "Date":
					currentFormat := "Mon, _2 Jan 2006 15:04:05 -0700 (MST)"
					date, err := time.Parse(currentFormat, splitedline[1])
					if err != nil {
						fmt.Println(err)
						date = time.Now()
					}
					item[v] = date.Format(time.RFC3339)
				default:
					item[v] = splitedline[1]
				}
			}
		}
	}
	spaceSplitedText := strings.Split(rawContent, "\n\n")
	if len(spaceSplitedText) >= 2 {
		item["content"] = strings.Join(spaceSplitedText[1:], "\n")
	} else {
		// TODO: validate empty files
		item["content"] = strings.Join(contentLines[15:], "\n")

	}
	item["_id"] = item["Message-ID"]

	return item
}

func getMapFile(header *tar.Header, tarReader *tar.Reader) (map[string]string, error) {
	name := header.Name

	var contentBuffer bytes.Buffer

	var item map[string]string

	switch header.Typeflag {
	case tar.TypeDir:
		return item, errors.New("not file")
	case tar.TypeReg:
		if _, err := io.Copy(&contentBuffer, tarReader); err != nil {
			fmt.Println("Error...")
			return item, err
		}
		rawContent := contentBuffer.String()
		rawContent = strings.Replace(rawContent, "\r\n", "\n", -1)

		item = getMapContent(rawContent)
		item["fileName"] = header.Name
		return item, nil
	default:
		fmt.Printf("%s : %c %s %s\n",
			"Yikes! Unable to figure out type",
			header.Typeflag,
			"in file",
			name,
		)
		return item, errors.New("not file")
	}
}

func GetFileAsMapList(tarFile *gzip.Reader) ([]map[string]string, error) {
	tarReader := tar.NewReader(tarFile)
	var fileContents []map[string]string

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
		}

		item, err := getMapFile(header, tarReader)

		if err != nil {
			if err.Error() == "not file" {
				continue
			} else {
				fmt.Println(err)

			}
		} else {
			fileContents = append(fileContents, item)
		}
	}
	return fileContents, nil
}

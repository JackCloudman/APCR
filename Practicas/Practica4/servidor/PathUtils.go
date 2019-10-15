package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
)

type FileInfo struct {
	Name string
	Size string
	Type bool
	Date string
	Mote string
}

func getPathInfo(Path string) []byte {
	var fileinfos []*FileInfo
	files, err := ioutil.ReadDir(Path)
	if err != nil {
		fmt.Println(err)
	}
	FilePath := splitregex.Split(Path, 2)[1][1:]
	for _, f := range files {
		newFile := new(FileInfo)
		newFile.Name = f.Name()
		newFile.Type = f.IsDir()
		newFile.Size = ByteCountDecimal(f.Size())
		newFile.Date = f.ModTime().String()
		fileinfos = append(fileinfos, newFile)
		if len(newFile.Name) > 20 {
			newFile.Mote = newFile.Name[0:17] + "..."
		} else {
			newFile.Mote = newFile.Name
		}
		newFile.Name = FilePath + "/" + newFile.Name
	}
	parent, _ := path.Split(FilePath)
	sort.SliceStable(fileinfos, func(i, j int) bool {
		return fileinfos[i].Type
	})
	return GenerateIndexOf(fileinfos, FilePath, parent)
}
func GenerateIndexOf(files []*FileInfo, FilePath, parent string) []byte {
	htmltemplate, _ := ioutil.ReadFile("indexof.html")
	if len(parent) > 0 {
		parent = parent[:len(parent)-1]
	}
	data := struct {
		Title  string
		Parent string
		Items  []*FileInfo
	}{
		Title:  FilePath,
		Items:  files,
		Parent: parent,
	}
	t, err := template.New("webpage").Parse(string(htmltemplate))
	fmt.Println("ERROR:", err)
	var tpl bytes.Buffer
	t.Execute(&tpl, data)
	return tpl.Bytes()

}
func GetFileContentType(filePath string) (string, error) {
	f, _ := os.Open(filePath)
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := f.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

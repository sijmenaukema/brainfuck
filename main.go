package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

type program struct {
	data         []byte
	stack        []byte
	stackPointer int
	dataPointer  int
	output       string
}

func main() {
	filePaths, err := getFilePath()
	if err != nil {
		log.WithField("file path", filePaths).Error(err)
	}

	for _, filePath := range filePaths {
		start := time.Now()
		reader, err := newReader(filePath)
		if err != nil {
			log.WithField("reader", reader).Error(err)
		}
		defer reader.Close()

		p := &program{
			data:         make([]byte, 30000),
			stack:        make([]byte, 30000),
			stackPointer: 0,
			dataPointer:  0,
			output:       "",
		}

		//Now read data into byte slice
		n, err := reader.Read(p.data)
		if err != nil {
			log.Fatalf("error reading from reader file: %v", err)
		}
		log.WithField("bytes read ", n).Info()
		p.interpret()
		if err != nil {
			log.WithField("interpreter", p.output).Error(err)
		}
		duration := time.Since(start)
		fmt.Printf("%s \nexecuted in: %v\n", p.output, duration)
	}
}

func newReader(path string) (*os.File, error) {
	input, err := os.Open(path)
	if err != nil {
		return input, err
	}
	return input, nil
}

func getFilePath() ([]string, error) {
	absPath, _ := os.Getwd()
	var filePath []string
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/local/", absPath))
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filePath = append(filePath, fmt.Sprintf(`%s/local/%s`, absPath, file.Name()))
		fmt.Println(filePath)
	}
	return filePath, err
}

func (p *program) interpret() {
	for ; p.dataPointer < len(p.data); p.dataPointer++ {
		switch p.data[p.dataPointer] {
		case '>':
			p.stackPointer++
		case '<':
			p.stackPointer--
		case '+':
			p.stack[p.stackPointer]++
		case '-':
			p.stack[p.stackPointer]--
		case '.':
			p.output += string(p.stack[p.stackPointer])
		case ',':
			p.stack[p.stackPointer] = p.data[p.dataPointer]
		case '[':
			if p.stack[p.stackPointer] == 0 {
				loop := 1
				for loop > 0 {
					p.dataPointer++
					if p.data[p.dataPointer] == '[' {
						loop++
					}
					if p.data[p.dataPointer] == ']' {
						loop--
					}
				}
			}
		case ']':
			if p.stack[p.stackPointer] != 0 {
				loop := 1
				for loop > 0 {
					p.dataPointer--
					if p.data[p.dataPointer] == '[' {
						loop--
					}
					if p.data[p.dataPointer] == ']' {
						loop++
					}
				}
			}
		default:
			break
		}
	}
}

package ini

import (
	"fmt"
	"io"
	"os"
)

func (d *IniDoc) Save(filename string) error {
	if filename == "" {
		return fmt.Errorf("invalid filepath: '%s'", filename)
	}

	result := d.ToString()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(result)

	return err
}

func Load(filename string) (*IniDoc, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	doc := Parse(string(contents))
	return doc, nil
}

package emulator

import (
	"encoding/csv"
	"os"
)

type CSV struct {
	fh  *os.File
	csv *csv.Writer
}

func (c *CSV) Write(d []string) error {
	return c.csv.Write(d)
}

func (c *CSV) Close() error {
	c.csv.Flush()
	err := c.csv.Error()
	if err != nil {
		return err
	}

	return c.fh.Close()
}

func MustNewCSV(f string) *CSV {
	c, err := NewCSV(f)
	if err != nil {
		panic(err)
	}

	return c
}

func NewCSV(f string) (*CSV, error) {
	fh, err := os.Create(f)
	if err != nil {
		return nil, err
	}

	return &CSV{fh, csv.NewWriter(fh)}, nil
}

package discover

import (
	"os"

	"github.com/andygrunwald/vdf"
)

func parseVDFMap(path string) (map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	p := vdf.NewParser(f)
	return p.Parse()
}

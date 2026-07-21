package server

import (
	"io/fs"
	"path/filepath"
	"strings"

	solidManager "github.com/wrhz/solid/manager"
)

func LoadHTML() error {
	render := solidManager.GetTemplateConfig().GetTemplateRender()
	dir := filepath.Join("resource", "views")

	err := filepath.WalkDir("./resource/views", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".html") {
			relPath, err := filepath.Rel(dir, path)

			if err != nil {
				return err
			}

			return render.New(relPath)
		}
		
		return nil
	})

	return err
}
package helper

import (
	"fmt"
	"github.com/denormal/go-gitignore"
	"github.com/ngyewch/asdf-helper/asdf"
	"github.com/ngyewch/asdf-helper/util"
	"os"
	"path/filepath"
	"strings"
)

func walk(handler func(asdfHelper *asdf.Helper, name string, version string) error) error {
	ignore, err := gitignore.NewRepository(".")
	if err != nil {
		return err
	}

	helper, err := asdf.NewHelper()
	if err != nil {
		return err
	}
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		match := ignore.Relative(path, info.IsDir())
		if match != nil {
			if match.Ignore() {
				if info.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
			}
		}
		if info.Name() == ".tool-versions" {
			fmt.Println()
			fmt.Println(path)

			pluginMap := make(map[string]string, 0)
			dir := filepath.Dir(path)
			err := util.ScanFile(filepath.Join(dir, ".plugin-versions"), func(line string) error {
				if line == "" {
					return nil
				}
				parts := strings.Split(line, " ")
				if len(parts) == 2 {
					name := parts[0]
					gitUrl := parts[1]
					pluginMap[name] = gitUrl
				}
				return nil
			})
			if err != nil && !os.IsNotExist(err) {
				return err
			}

			err = util.ScanFile(path, func(line string) error {
				if line == "" {
					return nil
				}
				parts := strings.Split(line, " ")
				if len(parts) == 2 {
					name := parts[0]
					version := parts[1]
					hasPlugin, err := helper.CheckPlugin(name)
					if err != nil {
						return err
					}
					if !hasPlugin {
						gitUrl, ok := pluginMap[name]
						if ok {
							err = helper.AddCustomPlugin(name, gitUrl)
							if err != nil {
								return err
							}
						} else {
							err = helper.AddPlugin(name)
							if err != nil {
								return err
							}
						}
					}

					err = handler(helper, name, version)
					if err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

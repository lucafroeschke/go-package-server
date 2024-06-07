package cmd

import (
	"fmt"
	"github.com/lucafroeschke/go-package-server/args"
	"github.com/lucafroeschke/go-package-server/config"
	"github.com/lucafroeschke/go-package-server/server"
	"github.com/lucafroeschke/go-package-server/templates"
	"html/template"
	"os"
	"path/filepath"
)

func RunBuildCommand() {
	cfg := config.GetConfig()

	output := args.GetString("output")
	isCloudflarePages := args.GetBool("cloudflare-pages")

	if output == "" {
		dir, _ := os.Getwd()
		output = filepath.Join(dir, "build")
	}

	if _, err := os.Stat(output); !os.IsNotExist(err) {
		fmt.Println("Output directory already exists, please remove it or specify a different output directory")
		return
	}

	fmt.Println("Building to", output)
	os.MkdirAll(output, 0755)

	tmpl, _ := template.ParseFS(templates.Templates, "package.html")

	for _, pkg := range cfg.Packages {
		fmt.Println("Building package", pkg.Name)

		if isCloudflarePages {
			pkgFile, err := os.Create(filepath.Join(output, fmt.Sprintf("%s.html", pkg.Name)))
			if err != nil {
				fmt.Println("Error creating package html file:", err)
				continue
			}

			err = tmpl.Execute(pkgFile, server.PackageResponse{Config: *cfg, Package: pkg})
		} else {
			pkgDir := filepath.Join(output, pkg.Name)
			os.MkdirAll(pkgDir, 0755)

			indexFile, err := os.Create(filepath.Join(pkgDir, "index.html"))
			if err != nil {
				fmt.Println("Error creating index.html:", err)
				continue
			}

			err = tmpl.Execute(indexFile, server.PackageResponse{Config: *cfg, Package: pkg})
			if err != nil {
				fmt.Println("Error executing template:", err)
			}

			err = indexFile.Close()
			if err != nil {
				return
			}
		}
	}

	tmpl, _ = template.ParseFS(templates.Templates, "index.html")

	indexFile, err := os.Create(filepath.Join(output, "index.html"))
	if err != nil {
		fmt.Println("Error creating index.html:", err)
		return
	}

	err = tmpl.Execute(indexFile, cfg)
	if err != nil {
		fmt.Println("Error executing template:", err)
	}
}

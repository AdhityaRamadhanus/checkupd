package commands

import (
	"html/template"
	"log"
	"os"

	"io/ioutil"

	"github.com/AdhityaRamadhanus/checkup"
	"github.com/AdhityaRamadhanus/checkupd/config"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// Main Function
func setupS3StatusPage(cliContext *cli.Context) {
	s3config := config.S3Config{
		AccessKeyID:     cliContext.String("accesskeyid"),
		SecretAccessKey: cliContext.String("secretaccesskey"),
		Region:          cliContext.String("region"),
		Bucket:          cliContext.String("bucket"),
	}
	var s3Error error
	if len(s3config.AccessKeyID) == 0 {
		log.Println("Please Provide AccessKeyID")
		s3Error = errors.New("S3 Config Error")
	}

	if len(s3config.SecretAccessKey) == 0 {
		log.Println("Please Provide Secret AccessKey")
		s3Error = errors.New("S3 Config Error")
	}

	if len(s3config.Region) == 0 {
		log.Println("Please Provide S3 Region")
		s3Error = errors.New("S3 Config Error")
	}

	if len(s3config.Bucket) == 0 {
		log.Println("Please Provide S3 Bucket")
		s3Error = errors.New("S3 Config Error")
	}

	if s3Error != nil {
		log.Println(s3Error)
		return
	}

	if err := setupS3(cliContext.String("url"), s3config); err != nil {
		log.Println(err)
	}
}

func setupFSStatusPage(cliContext *cli.Context) {
	if err := setupFS(cliContext.String("url")); err != nil {
		log.Println(err)
	}
}

// Helper
func setupFS(url string) error {
	// Create directory for logs
	// ignore the error
	log.Println("Setting up directory")
	os.Mkdir("./logs", 0777)
	os.Mkdir("./caddy-logs", 0777)
	os.Mkdir("./caddy-errors", 0777)
	// setup checkup.json
	log.Println("Creating checkup.json")
	checkup := checkup.Checkup{
		Checkers: []checkup.Checker{},
		Storage: checkup.FS{
			Dir: "/checkup/logs",
		},
	}
	jsonBytes, err := checkup.MarshalJSON()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(config.DefaultCheckupJSON, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	defer file.Close()
	if err != nil {
		return errors.Wrap(err, "Failed opening checkup.json file")
	}
	if _, err = file.Write(jsonBytes); err != nil {
		return errors.Wrap(err, "Failed to write checkup.json")
	}
	// setup Caddyfile
	// Read the template
	log.Println("Creating Caddyfile")
	caddyTemplate, err := template.ParseFiles(config.DefaultTplCaddyfile)
	if err != nil {
		return errors.Wrap(err, "Error parsing caddy template config")
	}
	caddyFile, err := os.OpenFile(config.DefaultCaddyFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
	defer caddyFile.Close()
	if err != nil {
		return errors.Wrap(err, "Error opening Caddyfile")
	}
	// Execute the template
	if err := caddyTemplate.Execute(caddyFile, struct{ URL string }{URL: url}); err != nil {
		return errors.Wrap(err, "Failed writing to Caddyfile")
	}

	// setup config status page
	log.Println("Creating config.js for status page")
	srcConfigBytes, err := ioutil.ReadFile(config.DefaultTplFSJS)
	if err != nil {
		return errors.Wrap(err, "Failed opening config.js")
	}

	dstConfigFile, err := os.OpenFile(config.DefaultConfigJS, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	defer dstConfigFile.Close()
	if _, err = dstConfigFile.Write(srcConfigBytes); err != nil {
		return errors.Wrap(err, "Failed to write config.js")
	}

	// setup config status page
	log.Println("Creating index.html for status page")

	srcIndexHTML, err := ioutil.ReadFile(config.DefaultTplIndexFS)
	if err != nil {
		return errors.Wrap(err, "Failed opening index_fs.html")
	}

	dstIndexHTML, err := os.OpenFile(config.DefaultIndexHtml, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	defer dstConfigFile.Close()
	if _, err = dstIndexHTML.Write(srcIndexHTML); err != nil {
		return errors.Wrap(err, "Failed to write index.html")
	}
	log.Println("Success!")
	return nil
}

// Helper
func setupS3(url string, s3 config.S3Config) error {
	// Create directory for logs
	// ignore the error
	log.Println("Setting up directory")
	os.Mkdir("./logs", 0777)
	os.Mkdir("./caddy-logs", 0777)
	os.Mkdir("./caddy-errors", 0777)
	// setup checkup.json
	log.Println("Creating checkup.json")
	checkup := checkup.Checkup{
		Checkers: []checkup.Checker{},
		Storage: checkup.S3{
			AccessKeyID:     s3.AccessKeyID,
			SecretAccessKey: s3.SecretAccessKey,
			Region:          s3.Region,
			Bucket:          s3.Bucket,
		},
	}
	jsonBytes, err := checkup.MarshalJSON()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(config.DefaultCheckupJSON, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	defer file.Close()
	if err != nil {
		return errors.Wrap(err, "Failed opening checkup.json file")
	}
	if _, err = file.Write(jsonBytes); err != nil {
		return errors.Wrap(err, "Failed to write checkup.json")
	}
	// setup Caddyfile
	// Read the template
	log.Println("Creating Caddyfile")
	caddyTemplate, err := template.ParseFiles(config.DefaultTplCaddyfile)
	if err != nil {
		return errors.Wrap(err, "Error parsing caddy template config")
	}
	caddyFile, err := os.OpenFile(config.DefaultCaddyFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
	defer caddyFile.Close()
	if err != nil {
		return errors.Wrap(err, "Error opening Caddyfile")
	}
	// Execute the template
	if err := caddyTemplate.Execute(caddyFile, struct{ URL string }{URL: url}); err != nil {
		return errors.Wrap(err, "Failed writing to Caddyfile")
	}

	// setup config status page
	log.Println("Creating config.js for status page")
	s3ConfigTemplate, err := template.ParseFiles(config.DefaultTplS3JS)
	if err != nil {
		return errors.Wrap(err, "Error parsing s3 template config")
	}

	dstConfigFile, err := os.OpenFile(config.DefaultConfigJS, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	defer dstConfigFile.Close()
	if err := s3ConfigTemplate.Execute(dstConfigFile, s3); err != nil {
		return errors.Wrap(err, "Failed to write config.js")
	}

	// setup config status page
	log.Println("Creating index.html for status page")

	srcIndexHTML, err := ioutil.ReadFile(config.DefaultTplIndexS3)
	if err != nil {
		return errors.Wrap(err, "Failed opening index_fs.html")
	}

	dstIndexHTML, err := os.OpenFile(config.DefaultIndexHtml, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	defer dstConfigFile.Close()
	if _, err = dstIndexHTML.Write(srcIndexHTML); err != nil {
		return errors.Wrap(err, "Failed to write index.html")
	}
	log.Println("Success!")
	return nil
}

package cmd

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/google/go-github/v34/github"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type GithubConnection struct {
	client      *github.Client
	ctx         context.Context
	url         string
	owner       string
	repo        string
	ansibleRole AnsibleRole
}

func NewConnectionBuilder(urlPath string) *GithubConnection {
	return &GithubConnection{
		client:      github.NewClient(nil),
		ctx:         context.Background(),
		url:         urlPath,
		owner:       strings.Split(urlPath, "/")[1],
		repo:        strings.Split(urlPath, "/")[2],
		ansibleRole: AnsibleRole{},
	}
}

// func getContents get Ansible role info from GitHub repository
// look only into role directories such as meta, defaults, tasks etc.
func (cb *GithubConnection) getContents(path string, parentDirName string, tempDir string) (err error) {
	_, directoryContent, _, err := cb.client.Repositories.GetContents(cb.ctx, cb.owner, cb.repo, path, nil)
	if err != nil {
		return err
	}

	for _, c := range directoryContent {
		switch *c.Type {
		case "dir":
			switch *c.Path {
			case "defaults", "handlers", "meta", "tasks", "vars":
				cb.parseFile(path, *c.Path)
			}
		case "templates", "files":
			var dirName = filepath.Join(path, *c.Path)
			var name = ""
			switch dirName {
			case "templates", "files":
				name, err := ioutil.TempDir("", dirName)
				if err != nil {
					log.Print(err)
				}
				tempDir = dirName
				fmt.Println("Temp dir name:", name)
			}

			err := cb.getContents(filepath.Join(path, *c.Path), *c.Name, name)
			if err != nil {
				log.Print(err)
			}
		}
	}
	return nil
}
func (cb *GithubConnection) parseFile(path string, parentDirName string) (err error) {
	_, directoryContent, _, err := cb.client.Repositories.GetContents(cb.ctx, cb.owner, cb.repo, path, nil)
	if err != nil {
		return err
	}

	for _, c := range directoryContent {
		switch *c.Type {
		case "file":
			if *c.Name == "main.yaml" || *c.Name == "main.yml" {
				b, err := downloadContents(cb, c)
				if err != nil {
					return err
				} else {
					ra := reflect.ValueOf(&cb.ansibleRole).Elem()
					ra.FieldByName(strings.Title(parentDirName) + "Main").SetBytes(b)
				}
			} else {
				b, err := downloadContents(cb, c)
				if err != nil {
					return err
				} else {
					ra := reflect.ValueOf(&cb.ansibleRole).Elem()
					ra.FieldByName(strings.Title(parentDirName)).SetBytes(b)
				}
			}

		}
	}
	return
}
func downloadContents(cb *GithubConnection, content *github.RepositoryContent) ([]byte, error) {
	rc, _, _, errDownload := cb.client.Repositories.DownloadContentsWithMeta(cb.ctx, cb.owner, cb.repo, *content.Path, nil)
	if errDownload != nil {
		return nil, errDownload
	}
	defer rc.Close()
	var b, errRead = ioutil.ReadAll(rc)
	if errRead != nil {
		return nil, errRead
	}
	sha := calculateGitSHA1(b)
	if *content.SHA == hex.EncodeToString(sha) {
		log.Print("SHA verified for a file: " + *content.Path)
		return b, nil
	} else {
		return nil, &cmdError{
			s:         "Invalid SHA, retry download for a file: " + *content.Path,
			userError: false,
		}
	}
}

func calculateGitSHA1(contents []byte) []byte {
	contentLen := len(contents)
	blobSlice := []byte("blob " + strconv.Itoa(contentLen))
	blobSlice = append(blobSlice, '\x00')
	blobSlice = append(blobSlice, contents...)
	h := sha1.New()
	h.Write(blobSlice)
	bs := h.Sum(nil)
	return bs
}

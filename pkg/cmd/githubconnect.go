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

type SHAError struct {
	s string
}

func (S SHAError) Error() string {
	return S.s
}

type GithubConnection struct {
	client      *github.Client
	ctx         context.Context
	url         string
	owner       string
	repo        string
	ansibleRole ansibleRole
}

func NewConnectionBuilder(urlPath string) *GithubConnection {
	return &GithubConnection{
		client:      github.NewClient(nil),
		ctx:         context.Background(),
		url:         urlPath,
		owner:       strings.Split(urlPath, "/")[1],
		repo:        strings.Split(urlPath, "/")[2],
		ansibleRole: ansibleRole{},
	}
}

func (cb *GithubConnection) getContents(path string, parentDirName string) (err error) {
	_, directoryContent, _, err := cb.client.Repositories.GetContents(cb.ctx, cb.owner, cb.repo, path, nil)
	if err != nil {
		return err
	}

	for _, c := range directoryContent {
		switch *c.Type {
		case "file":
			switch parentDirName {
			case "defaults", "templates", "handlers", "meta", "tasks", "vars", "files", "library":
				if *c.Name == "main.yaml" || *c.Name == "main.yml" {
					b, err := downloadContents(cb, c)
					if err != nil {
						return err
					} else {
						//TODO try capitalize struct fields
						reflect.ValueOf(&cb.ansibleRole).Elem().FieldByName(parentDirName).SetBytes(b)
						cb.ansibleRole.defaults = b
					}
				} else {
					fmt.Print("TODO")
				}
			}

		case "dir":
			err := cb.getContents(filepath.Join(path, *c.Path), *c.Name)
			if err != nil {
				log.Print(err)
			}
		}
	}
	return nil
}

func downloadContents(cb *GithubConnection, content *github.RepositoryContent) ([]byte, error) {
	rc, _, _, err_download := cb.client.Repositories.DownloadContentsWithMeta(cb.ctx, cb.owner, cb.repo, *content.Path, nil)
	if err_download != nil {
		return nil, err_download
	}
	defer rc.Close()
	var b, err_read = ioutil.ReadAll(rc)
	if err_read != nil {
		return nil, err_read
	}
	sha := calculateGitSHA1(b)
	if *content.SHA == hex.EncodeToString(sha) {
		log.Print("SHA verified for" + *content.Path)
		fmt.Println(string(b))
		return b, nil
	} else {
		return nil, &SHAError{"invalid SHA, retry for file"}
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
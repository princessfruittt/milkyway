package cmd

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/google/go-github/v34/github"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

type generateCmd struct {
	*baseBuilderCmd
}

func (b *cmdsBuilder) newGenerateCmd() *generateCmd {
	cc := &generateCmd{}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate tosca types",
		Run:   cc.generateTypes,
	}
	cc.baseBuilderCmd = b.newBuilderCmd(cmd)
	return cc
}
func (c *generateCmd) printOutput(cmd *cobra.Command, args []string) {
	fmt.Println("I am in generated func")
}
func (c *generateCmd) generateTypes(cmd *cobra.Command, args []string) {
	client := github.NewClient(nil)
	ctx := context.Background()
	u, err := url.Parse(c.roleURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u.Path)
	info := strings.Split(u.Path, "/")
	getContents(ctx, client, info[1], info[2], "", "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(opt.GetSource().Source)
}

func getContents(ctx context.Context, client *github.Client, owner string, repo string, path string, parentDirName string) {
	fileContent, directoryContent, response, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", fileContent)
	fmt.Printf("%#v\n", directoryContent)
	fmt.Printf("%#v\n", response)

	for _, c := range directoryContent {
		fmt.Println(*c.Type, *c.Size, *c.SHA)

		switch *c.Type {
		case "file":
			switch parentDirName {
			case "templates":
				fmt.Println(parentDirName)
				downloadContents(ctx, owner, repo, client, c)
			case "defaults":
				fmt.Println(parentDirName)
				downloadContents(ctx, owner, repo, client, c)
			case "handlers":
				fmt.Println(parentDirName)
				downloadContents(ctx, owner, repo, client, c)
			case "meta":
				fmt.Println(parentDirName)
				downloadContents(ctx, owner, repo, client, c)
			case "tasks":
				fmt.Println(parentDirName)
				downloadContents(ctx, owner, repo, client, c)
			case "vars":
				fmt.Println(parentDirName)
				downloadContents(ctx, owner, repo, client, c)
			}

		case "dir":
			getContents(ctx, client, owner, repo, filepath.Join(path, *c.Path), *c.Name)
		}
	}
}

func downloadContents(ctx context.Context, owner string, repo string, client *github.Client, content *github.RepositoryContent) {
	rc, _, _, err_download := client.Repositories.DownloadContentsWithMeta(ctx, owner, repo, *content.Path, nil)
	if err_download != nil {
		fmt.Println(err_download)
		return
	}
	defer rc.Close()
	var b, err_read = ioutil.ReadAll(rc)
	if err_read != nil {
		fmt.Println(err_read)
		return
	}
	sha := calculateGitSHA1(b)
	if *content.SHA == hex.EncodeToString(sha) {
		fmt.Println("no need to update this file, the SHA is the same")
		fmt.Print(string(b))
	} else {
		fmt.Println("SHA err")
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

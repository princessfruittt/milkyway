package cmd

import (
	"context"
	"github.com/google/go-github/v34/github"
	"reflect"
	"testing"
)

func TestGithubConnection_getContents(t *testing.T) {
	type fields struct {
		client      *github.Client
		ctx         context.Context
		url         string
		owner       string
		repo        string
		ansibleRole AnsibleRole
	}
	type args struct {
		path          string
		parentDirName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := &GithubConnection{
				client:      tt.fields.client,
				ctx:         tt.fields.ctx,
				url:         tt.fields.url,
				owner:       tt.fields.owner,
				repo:        tt.fields.repo,
				ansibleRole: tt.fields.ansibleRole,
			}
			if err := cb.getContents(tt.args.path, tt.args.parentDirName); (err != nil) != tt.wantErr {
				t.Errorf("getContents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewConnectionBuilder(t *testing.T) {
	type args struct {
		urlPath string
	}
	tests := []struct {
		name string
		args args
		want *GithubConnection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConnectionBuilder(tt.args.urlPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConnectionBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateGitSHA1(t *testing.T) {
	type args struct {
		contents []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateGitSHA1(tt.args.contents); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateGitSHA1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_downloadContents(t *testing.T) {
	type args struct {
		cb      *GithubConnection
		content *github.RepositoryContent
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := downloadContents(tt.args.cb, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("downloadContents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("downloadContents() got = %v, want %v", got, tt.want)
			}
		})
	}
}

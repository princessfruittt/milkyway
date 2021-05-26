package cmd

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// AnsibleRole struct for Ansible role v2.9 and less
type AnsibleRole struct {
	TasksMain    []byte
	VarsMain     []byte
	DefaultsMain []byte
	MetaMain     []byte
	Templates    fs.FileInfo
	Tasks        [][]byte
	Vars         [][]byte
	Defaults     [][]byte
	Handlers     fs.FileInfo
	Files        fs.FileInfo
	//Library       []byte
	Version string
}

func NewAnsibleRole() *AnsibleRole {
	return &AnsibleRole{
		TasksMain:    nil,
		VarsMain:     nil,
		DefaultsMain: nil,
		MetaMain:     nil,
		Templates:    nil,
		Tasks:        nil,
		Vars:         nil,
		Defaults:     nil,
		Handlers:     nil,
		Files:        nil,
		Version:      "",
	}
}

type ansibleRoleMeta struct {
	Meta         GalaxyMeta `yaml:"galaxy_info,omitempty"`
	Dependencies []string   `yaml:"dependencies,omitempty"`
}

type GalaxyMeta struct {
	RoleName          string     `yaml:"role_name,omitempty"`
	Author            string     `yaml:"author,omitempty"`
	Description       string     `yaml:"description,omitempty"`
	Platforms         []Platform `yaml:"platforms,omitempty"`
	MinAnsibleVersion string     `yaml:"min_ansible_version,omitempty"`
}
type Platform struct {
	Name     string   `yaml:"name,omitempty"`
	Versions []string `yaml:"versions,omitempty"`
}
type Data interface {
}
type Datayaml struct {
	Data []Data `yaml:",inline"`
}

type role struct {
	r AnsibleRole
}

func newRole() *role {
	return &role{
		r: AnsibleRole{},
	}
}

// GetRoleFromPath return struct with filled AnsibleRole fields from absolute path input.
func GetRoleFromPath(rolePath string) AnsibleRole {
	r := newRole()
	e := filepath.WalkDir(rolePath, func(path string, de fs.DirEntry, err error) error {
		switch de.IsDir() {
		case true:
			switch de.Name() {
			case "defaults", "meta", "tasks", "vars":
				err := getRoleContent(filepath.Join(rolePath, de.Name()), de.Name(), r)
				if err != nil {
					return err
				}
			case "templates", "files", "handlers":
				i, err := de.Info()
				if err != nil {
					return err
				}
				f := reflect.ValueOf(&r.r).Elem()
				f.FieldByName(strings.Title(de.Name())).Set(reflect.ValueOf(i))
			}
			return err
		}
		return err
	})

	if e != nil {
		log.Fatal(e)
	}
	return r.r
}

func getRoleContent(rolePath string, parentDirName string, r *role) error {
	e := filepath.WalkDir(rolePath, func(path string, de fs.DirEntry, err error) error {
		if de.IsDir() == false {
			//By default Ansible will look in each directory within a role for a main.yml file for relevant content (also main.yaml and main)
			n := strings.ToLower(de.Name())
			content, err := os.ReadFile(filepath.Join(rolePath, de.Name()))
			if err != nil {
				return err
			}
			if n == "main.yaml" || n == "main.yml" || n == "main" {
				ra := reflect.ValueOf(&r.r).Elem()
				ra.FieldByName(strings.Title(parentDirName) + "Main").SetBytes(content)
			} else {
				ra := reflect.ValueOf(&r.r).Elem()
				if ra.FieldByName(strings.Title(parentDirName)).Kind() == reflect.Array {
					srcArr := [2]byte{}
					copy(srcArr[:], content)
					reflect.Copy(ra, reflect.ValueOf(srcArr))
				}
			}
			return err
		}
		return err
	})
	return e
}

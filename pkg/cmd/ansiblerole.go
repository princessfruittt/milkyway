package cmd

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
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
	Artifacts map[string]string
	Version   string
	Name      string
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
		Artifacts:    make(map[string]string),
		Version:      "",
		Name:         "Default",
	}
}

type ansibleRoleMeta struct {
	Meta         GalaxyMeta    `yaml:"galaxy_info,omitempty"`
	Dependencies []interface{} `yaml:"dependencies,omitempty"`
}

type GalaxyMeta struct {
	RoleName          string     `yaml:"role_name,omitempty"`
	Author            string     `yaml:"author,omitempty"`
	Description       string     `yaml:"description,omitempty"`
	Platforms         []Platform `yaml:"platforms,omitempty"`
	MinAnsibleVersion string     `yaml:"min_ansible_version,omitempty"`
}
type Platform struct {
	Name     string        `yaml:"name,omitempty"`
	Versions []interface{} `yaml:"versions,omitempty"`
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
		r: *NewAnsibleRole(),
	}
}

var dirWithMain = []string{"defaults", "meta", "tasks", "vars"}

func NilFields(x AnsibleRole) bool {
	rv := reflect.ValueOf(&x).Elem()

	for i := 0; i < rv.NumField(); i++ {
		if !rv.Field(i).IsNil() {
			return false
		}
	}
	return true
}

// GetRoleFromPath return struct with filled AnsibleRole fields from absolute path input.
func GetRoleFromPath(rolePath string) AnsibleRole {
	r := newRole()
	e := filepath.WalkDir(rolePath, func(path string, de fs.DirEntry, err error) error {
		if filepath.IsAbs(rolePath) == false {
			newPath, err := filepath.Abs(rolePath)
			if err != nil {
				return err
			}
			rolePath = newPath
		}
		if err != nil {
			return err
		}
		switch de.IsDir() {
		case true:
			dn := strings.ToLower(de.Name())
			switch dn {
			case "defaults", "meta", "tasks", "vars":
				err := getRoleContent(filepath.Join(rolePath, de.Name()), de.Name(), r)
				if err != nil {
					return err
				}
			case "templates", "files", "handlers":
				r.r.Artifacts[dn] = dn
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
		i := sort.Search(len(dirWithMain), func(i int) bool { return dirWithMain[i] >= de.Name() })
		if de.IsDir() == true && i < len(dirWithMain) && dirWithMain[i] == de.Name() {
			//By default Ansible will look in each directory within a role for a main.yml file for relevant content (also main.yaml and main)

		} else {
			dn := strings.ToLower(de.Name())
			content, err := os.ReadFile(filepath.Join(rolePath, de.Name()))
			if err != nil {
				return err
			}
			if dn == "main.yaml" || dn == "main.yml" || dn == "main" {
				ar := reflect.ValueOf(&r.r).Elem()
				ar.FieldByName(strings.Title(parentDirName) + "Main").SetBytes(content)
			} else {
				pdn := strings.ToLower(parentDirName)
				r.r.Artifacts[pdn] = pdn
				ar := reflect.ValueOf(&r.r).Elem()
				if ar.FieldByName(strings.Title(parentDirName)).Kind() == reflect.Slice {
					ar.FieldByName(strings.Title(parentDirName)).Set(reflect.Append(ar.FieldByName(strings.Title(parentDirName)), reflect.ValueOf(content)))
				}
			}
			return err
		}
		return err
	})
	return e
}

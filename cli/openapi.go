package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const ApplicationJSON = "application/json"
const MultipartFormData = "multipart/form-data"

type OpenAPI struct {
	Status     *string
	url        string
	Paths      map[string]Path `json:"paths"`
	Components Components      `json:"components"`
}

type Components struct {
	Schemas map[string]json.RawMessage `json:"schemas"`
}

type Path struct {
	Post PathPost `json:"post"`
}

type PathPost struct {
	ID              string              `json:"operationId"`
	Summary         string              `json:"summary"`
	Description     string              `json:"description"`
	Tags            []string            `json:"tags"`
	RequestBody     RequestBody         `json:"requestBody"`
	Responses       map[string]Response `json:"responses"`
	XPangeaUISchema *XPangeaUISchema    `json:"x-pangea-ui-schema"`
}

type XPangeaUISchema struct {
	IsConfiguration *bool `json:"isConfiguration"`
}

type RequestBody struct {
	Content Content `json:"content"`
}

type Response struct {
	Description string  `json:"description"`
	Content     Content `json:"content"`
}

type Content map[string]struct {
	Schema *Schema `json:"schema"`
}

type Schema struct {
	Ref           string         `json:"$ref"`
	OneOf         []Schema       `json:"oneOf"`
	AnyOf         []Schema       `json:"anyOf"`
	Discriminator *Discriminator `json:"discriminator"`

	Title       string     `json:"title"`
	Description string     `json:"description"`
	Required    []string   `json:"required"`
	Properties  Properties `json:"properties"`
}

type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]Schema `json:"mapping"`
}

type Properties map[string]Property

type Property struct {
	Ref string `json:"$ref"`

	Type        json.RawMessage `json:"type"` // TODO: this can be a string or an array of strings
	Default     any             `json:"default"`
	Description string          `json:"description"`
	Format      string          `json:"format"`
	Enum        []any           `json:"enum"`
	Const       any             `json:"const"`
}

func (oa *OpenAPI) resolveReferences() error {
	for name, path := range oa.Paths {
		// request body
		for mime, content := range path.Post.RequestBody.Content {
			err := oa.resolveSchema(&content.Schema)
			if err != nil {
				return nil
			}

			err = oa.resolveDiscriminator(content.Schema.Discriminator)
			if err != nil {
				return nil
			}

			err = oa.resolvePropertiesReferences(content.Schema.Properties)
			if err != nil {
				return err
			}
			oa.Paths[name].Post.RequestBody.Content[mime] = content
		}
		// responses
		for mime, response := range path.Post.Responses {
			for mime, content := range response.Content {
				err := oa.resolveSchema(&content.Schema)
				if err != nil {
					return nil
				}
				response.Content[mime] = content
				err = oa.resolvePropertiesReferences(content.Schema.Properties)
				if err != nil {
					return err
				}
			}
			path.Post.Responses[mime] = response
		}
	}
	return nil
}

func (oa *OpenAPI) resolveSchema(source **Schema) error {
	if (*source).Ref != "" {
		var target Schema
		err := oa.schemaRef((*source).Ref, &target)
		if err != nil {
			return err
		}
		*source = &target
		err = oa.resolveSchema(source)
		if err != nil {
			return err
		}
	}

	for i, sch := range (*source).OneOf {
		psch := &sch
		err := oa.resolveSchema(&psch)
		if err != nil {
			return err
		}
		(*source).OneOf[i] = *psch
	}

	for i, sch := range (*source).AnyOf {
		psch := &sch
		err := oa.resolveSchema(&psch)
		if err != nil {
			return err
		}
		(*source).AnyOf[i] = *psch
	}

	err := oa.resolvePropertiesReferences((*source).Properties)
	if err != nil {
		return err
	}

	return nil
}

func (oa *OpenAPI) resolveDiscriminator(d *Discriminator) error {
	if d == nil {
		return nil
	}

	for k, v := range d.Mapping {
		sch := &v
		err := oa.resolveSchema(&sch)
		if err != nil {
			continue
		}
		d.Mapping[k] = v
	}
	return nil
}

func (oa *OpenAPI) resolvePropertiesReferences(props Properties) error {
	for k, v := range props {
		if v.Ref != "" {
			var prop Property
			err := oa.schemaRef(v.Ref, &prop)
			if err != nil {
				return err
			}
			props[k] = prop
		}
	}
	return nil
}

func (oa *OpenAPI) schemaRef(r string, val any) error {
	p := strings.Split(r, "#")
	if p[0] != oa.url {
		return errors.New("ref not allowed")
	}
	if !strings.HasPrefix(p[1], "/components/schemas/") {
		return errors.New("ref not allowed")
	}

	pp := strings.Split(p[1], "/")
	targetName := pp[len(pp)-1]
	d, ok := oa.Components.Schemas[targetName]
	if !ok {
		return errors.New("schema not found")
	}

	err := json.Unmarshal(d, val)
	if err != nil {
		return err
	}

	return nil
}

func LoadFile(fname string, url string) (*OpenAPI, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadReader(f, url)
}

func LoadReader(r io.Reader, url string) (*OpenAPI, error) {
	s, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	oapi, err := Unmarshal(s)
	if err != nil {
		return nil, err
	}

	oapi.url = url
	err = oapi.resolveReferences()
	if err != nil {
		return nil, err
	}

	return oapi, nil
}

func LoadURL(u string) (*OpenAPI, error) {
	data := LoadFromCache(u)
	if data == nil {
		fmt.Fprintf(os.Stderr, "Downloading from %s...\n", u)
		resp, err := http.Get(u)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = SaveToCache(u, data)
		if err == nil {
			fmt.Fprintf(os.Stderr, "Downloaded. Saved to cache.\n")
		} else {
			fmt.Fprintf(os.Stderr, "Downloaded but failed to save to cache.\n")
		}
	}

	reader := bytes.NewReader(data)
	return LoadReader(reader, u)
}

func Unmarshal(s []byte) (*OpenAPI, error) {
	var oapi OpenAPI
	err := json.Unmarshal(s, &oapi)
	if err != nil {
		return nil, err
	}
	return &oapi, nil
}

func GetCacheFilename(u string) (string, error) {
	pu, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	t := GetDaySinceEpoch()

	cacheDir, err := GetCacheFolder()
	if err != nil {
		return "", err
	}

	name := strings.ReplaceAll(fmt.Sprintf("%s_%s", pu.Host, pu.Path), "/", "_")
	path := filepath.Join(cacheDir, t, name)
	return path, nil
}

func LoadFromCache(u string) []byte {
	filename, err := GetCacheFilename(u)
	if err != nil {
		return nil
	}

	// Read the entire file into a byte slice
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	return fileContent
}

func SaveToCache(u string, data []byte) error {
	filename, err := GetCacheFilename(u)
	if err != nil {
		return err
	}

	dir := filepath.Dir(filename)

	// Check if the directory exists
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		// Directory does not exist, create it along with any necessary parent directories
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		// Other error occurred while checking directory existence
		return err
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func RemoveCachedFileFromURL(url string) error {
	filename, err := GetCacheFilename(url)
	if err != nil {
		return err
	}
	return os.Remove(filename)
}

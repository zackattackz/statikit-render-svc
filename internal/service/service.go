package service

import (
	"bytes"
	"io"
	"text/template"
)

type RenderService interface {
	Render(string, Schema, io.Writer) error
}

// Chainable behavior modifier for Renderer.
type Middleware func(RenderService) RenderService

// Implementation of RenderService
type renderService struct {
	cache Cache
}

// Used to inject dependencies into a render service
func New(cache Cache) renderService {
	return renderService{
		cache,
	}
}

func (rs renderService) Render(contents string, schema Schema, w io.Writer) error {
	// First check the cache for a result,
	// return it if it exists
	k := CacheKey{schema, contents}
	if res, err := rs.cache.Get(k); err == nil {
		_, err = w.Write([]byte(res))
		return err
	}

	// Render the template with the given contents/schema

	// Create the template with given contents
	t, err := template.New("input").Parse(contents)
	if err != nil {
		return err
	}

	// Execute template and store the result in a buffer
	buff := bytes.Buffer{}
	err = t.Execute(&buff, schema)
	if err != nil {
		return err
	}

	// Write the result to the output writer
	res := buff.Bytes()
	_, err = w.Write(res)
	if err != nil {
		return err
	}

	// Update the cache (ignore errors)
	rs.cache.Set(k, string(res))
	return nil
}

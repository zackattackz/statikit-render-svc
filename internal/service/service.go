package service

import (
	"bytes"
	"html/template"
	"io"

	"github.com/zackattackz/statikit-render-svc/internal/models"
	"github.com/zackattackz/statikit-render-svc/internal/ports"
)

type RenderService interface {
	// Uses html/template to render the input context with the given schema, writes result to input writer
	//
	// Returns error != nil if failed
	Render(string, models.Schema, io.Writer) error
}

// Chainable behavior modifier for RenderService.
type Middleware func(RenderService) RenderService

// Default implementation of RenderService
type renderService struct {
	cache ports.Cache
}

// Returns default implementation of a RenderService, using a cache.
func New(cache ports.Cache) RenderService {
	return renderService{
		cache,
	}
}

// Uses html/template to render the input context with the given schema, writes result to input writer
//
// Returns error != nil if failed
func (rs renderService) Render(contents string, schema models.Schema, w io.Writer) error {
	// First check the cache for a result,
	// return it if it exists
	k := ports.CacheKey{Schema: schema, Contents: contents}
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

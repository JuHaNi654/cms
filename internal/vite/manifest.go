package vite

import (
	"encoding/json"
	"io"
	"strings"
)

type Manifest map[string]*Chunk

type Chunk struct {
	File           string   `json:"file"`
	Name           string   `json:"name"`
	Src            string   `json:"src"`
	CSS            []string `json:"css"`
	IsDynamicEntry bool     `json:"IsDynamicEntry"`
	IsEntry        bool     `json:"isEntry"`
	Imports        []string `json:"imports"`
	DynamicImports []string `json:"dynamicImports"`
}

func parseManifest(r io.Reader) (*Manifest, error) {
	var m Manifest
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (m Manifest) GetEntryPoint() *Chunk {
	for _, chunk := range m {
		if chunk.IsEntry {
			return chunk
		}
	}

	return nil
}

func (m Manifest) GetEntryPoints() []*Chunk {
	var entryPoints []*Chunk
	for _, chunk := range m {
		if chunk.IsEntry {
			entryPoints = append(entryPoints, chunk)
		}
	}

	return entryPoints
}

func (m Manifest) GetChunk(name string) (*Chunk, bool) {
	chunk, ok := m[name]
	return chunk, ok
}

func (m Manifest) GenerateCSS(name string) string {
	var sb strings.Builder
	seen := make(map[string]bool)
	var addCSS func(string)

	addCSS = func(name string) {
		if seen[name] {
			return
		}
		seen[name] = true
		chunk, ok := m[name]

		if !ok {
			return
		}

		for _, css := range chunk.CSS {
			sb.WriteString(`<link rel="stylesheet" href="/`)
			sb.WriteString(css)
			sb.WriteString(`">`)
		}

		for _, imp := range chunk.Imports {
			addCSS(imp)
		}
	}

	addCSS(name)

	return sb.String()
}

func (m Manifest) GenerateModules(name string) string {
  chunk, ok := m[name]
  if !ok {
    return ""
  }
  
  var sb strings.Builder
  if chunk.File != "" {
    sb.WriteString(`<script type="module" src="/`)
    sb.WriteString(chunk.File)
    sb.WriteString(`"></script>`)
  }

  return sb.String()
}

func (m Manifest) GeneratePreloadModules(name string) string {
  var sb strings.Builder
  seen := make(map[string]bool)

  var addModulePreload func(string)
  addModulePreload = func(name string) {
    if seen[name] {
      return
    }
    seen[name] = true 
    chunk, ok := m[name]
    if !ok {
      return
    }

    if chunk.File != "" {
      sb.WriteString(`<script type="preloadmodule" src="/`)
      sb.WriteString(chunk.File)
      sb.WriteString(`"></script>`)
    }

    for _, imp := range chunk.Imports {
      addModulePreload(imp)
    }
  }

  addModulePreload(name)
  return sb.String()
}

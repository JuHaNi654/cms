package vite

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

type FsHandler struct {
	Fs      fs.FS
	HttpFs  http.FileSystem
	Handler http.Handler
}

func NewFsHandler(fs fs.FS) *FsHandler {
	return &FsHandler{
		Fs:      fs,
		HttpFs:  http.FS(fs),
		Handler: http.FileServerFS(fs),
	}
}

func (handler *FsHandler) setFs(fs fs.FS) {
	handler.Fs = fs
	handler.HttpFs = http.FS(fs)
	handler.Handler = http.FileServerFS(fs)
}

type Handler struct {
	fs              *FsHandler
	pub             *FsHandler // Assets
	manifest        *Manifest
	isProd          bool
	viteEntry       string
	viteURL         string
	defaultMetadata *Metadata
}

type Config struct {
	FS        fs.FS
	PublicFs  fs.FS // Assets
	IsProd    bool
	ViteEntry string
	ViteURL   string
}

func NewHandler(config Config) (*Handler, error) {
	if config.FS == nil {
		return nil, fmt.Errorf("vite: fs is nill")
	}

	handler := &Handler{
		fs:        NewFsHandler(config.FS),
		isProd:    config.IsProd,
		viteEntry: config.ViteEntry,
		viteURL:   config.ViteURL,
	}

	if handler.isProd { // Production mode
		manifest, err := handler.fs.Fs.Open(".vite/manifest.json")
		if err != nil {
			return nil, fmt.Errorf("vite: cannot open manifest.json: %w", err)
		}
		defer manifest.Close()

		handler.manifest, err = parseManifest(manifest)
		if err != nil {
			return nil, fmt.Errorf("vite: error while parsing manifest: %w", err)
		}
	} else { // Development mode
		if handler.viteURL == "" {
			handler.viteURL = "http://localhost:5173"
		}
	}

	return handler, nil
}

func (handler *Handler) SetDefaultMetadata(data *Metadata) {
	handler.defaultMetadata = data
}

func (h *Handler) GetPageData(r *http.Request) (*PageData, error) {
	page := &PageData{
		IsDev:   h.isProd == false,
		ViteURL: h.viteURL,
	}
	var chunk *Chunk

	ctx := r.Context()
	md := MetadataFromContext(ctx)
	if md == nil {
		md = h.defaultMetadata
	} else {
		page.Metadata = template.HTML(md.String())
	}

	if h.isProd {
		if h.viteEntry == "" {
			chunk = h.manifest.GetEntryPoint()
		} else {
			entries := h.manifest.GetEntryPoints()
			for _, entry := range entries {
				if h.viteEntry == entry.Src {
					chunk = entry
					break
				}
			}
		}

		if chunk == nil {
			return nil, fmt.Errorf("internal server error")
		}

		page.StyleSheets = template.HTML(h.manifest.GenerateCSS(chunk.Src))
		page.Modules = template.HTML(h.manifest.GenerateModules(chunk.Src))
		page.PreloadModules = template.HTML(h.manifest.GeneratePreloadModules(chunk.Src))
	}

	return page, nil
}

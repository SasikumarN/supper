package application

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/set"
	"github.com/tympanix/supper/api"
	"github.com/tympanix/supper/list"
	"github.com/tympanix/supper/media"
	"github.com/tympanix/supper/provider"
	"github.com/tympanix/supper/types"
	"github.com/urfave/cli"
	"golang.org/x/text/language"
)

var filetypes = []string{
	".avi", ".mkv", ".mp4", ".m4v", ".flv", ".mov", ".wmv", ".webm", ".mpg", ".mpeg",
}

// Application is an configuration instance of the application
type Application struct {
	types.Provider
	*http.ServeMux
	context *cli.Context
}

func New(context *cli.Context) types.App {
	app := &Application{
		Provider: provider.Subscene(),
		ServeMux: http.NewServeMux(),
		context:  context,
	}

	static := context.String("static")

	fs := http.FileServer(http.Dir(filepath.Join(static, "/static")))
	app.ServeMux.Handle("/static/", http.StripPrefix("/static", fs))

	api := api.New(app)
	app.ServeMux.Handle("/api/", http.StripPrefix("/api", api))

	index := IndexHandler(filepath.Join(static, "index.html"))
	app.ServeMux.Handle("/", index)

	return app
}

// Context returns the CLI context of the application (e.g. flags, args ect.)
func (a *Application) Context() *cli.Context {
	return a.context
}

func (a *Application) Languages() set.Interface {
	// Parse all language flags into slice of tags
	lang := set.New()
	for _, tag := range a.Context().GlobalStringSlice("lang") {
		_lang, err := language.Parse(tag)
		if err != nil {
			os.Exit(5)
		}
		lang.Add(_lang)
	}
	return lang
}

// IndexHandler always serves the same file (e.g. index.html)
func IndexHandler(entrypoint string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	})
}

func fileIsMedia(f os.FileInfo) bool {
	for _, ext := range filetypes {
		if ext == filepath.Ext(f.Name()) {
			return true
		}
	}
	return false
}

// FindMedia searches for media files
func (a *Application) FindMedia(roots ...string) (types.LocalMediaList, error) {
	medialist := make([]types.LocalMedia, 0)

	for _, root := range roots {
		if _, err := os.Stat(root); os.IsNotExist(err) {
			return nil, err
		}

		err := filepath.Walk(root, func(filepath string, f os.FileInfo, err error) error {
			if f.IsDir() {
				return nil
			}
			if !fileIsMedia(f) {
				return nil
			}
			_media, err := media.New(filepath)
			if err != nil {
				return nil
			}
			if media.IsSample(_media) {
				return nil
			}
			medialist = append(medialist, _media)
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return list.NewLocalMedia(medialist...), nil
}
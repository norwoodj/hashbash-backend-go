package frontend

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
)

type PageModel struct {
	Error          string
	RainbowTableId string
}

func AddFrontendFlags(flags *pflag.FlagSet) {
	flags.StringP("frontend-template-path", "f", "", "Path to directory containing frontend templates")
}

func createTemplateRegex(templateFileNames []string) *regexp.Regexp {
	templateNames := make([]string, 0)

	for _, t := range templateFileNames {
		if strings.HasPrefix(t, "_") {
			continue
		}

		templateBaseName := path.Base(t)
		templateNames = append(templateNames, strings.TrimSuffix(templateBaseName, ".html.gotmpl"))
	}

	regexString := fmt.Sprintf("^(%s)$", strings.Join(templateNames, "|"))
	return regexp.MustCompile(regexString)
}

func renderTemplate(
	w http.ResponseWriter,
	templateName string,
	templates *template.Template,
	pageModel PageModel,
) {
	err := templates.ExecuteTemplate(w, templateName+".html.gotmpl", pageModel)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RegisterTemplateHandler(router *mux.Router, frontendDirPath string) error {
	templateFiles, err := ioutil.ReadDir(frontendDirPath)
	if err != nil {
		return err
	}

	templateFilenames := make([]string, 0)

	for _, f := range templateFiles {
		templateFilenames = append(templateFilenames, path.Join(frontendDirPath, f.Name()))
	}

	templates := template.Must(template.ParseFiles(templateFilenames...))
	templateRegex := createTemplateRegex(templateFilenames)

	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/home", http.StatusPermanentRedirect)
	})

	router.HandleFunc("/{template}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		templateName := vars["template"]
		queryParams := request.URL.Query()

		if !templateRegex.MatchString(templateName) {
			http.NotFound(writer, request)
			return
		}

		rainbowTableId := queryParams.Get("rainbowTableId")

		if templateName == "search-rainbow-table" && rainbowTableId == "" {
			rainbowTableId := queryParams.Get("rainbowTableId")

			if rainbowTableId == "" {
				http.Redirect(writer, request, fmt.Sprintf("/rainbow-tables?error=%s", url.PathEscape(rainbowTableIdQueryParamRequired)), 303)
				return
			}
		}

		pageModel := PageModel{
			Error:          queryParams.Get("error"),
			RainbowTableId: rainbowTableId,
		}

		renderTemplate(writer, templateName, templates, pageModel)
	})

	return nil
}

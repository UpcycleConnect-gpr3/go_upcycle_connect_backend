package {{.PackageName}}

import (
	"{{.ModuleName}}/utils/log"
	"{{.ModuleName}}/utils/response"
	"net/http"
)

func Index{{.ResourceName}}Handler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	// TODO: Implement
	response.NewSuccessData(w, []string{})
}

func Store{{.ResourceName}}Handler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	// TODO: Implement
	response.NewSuccessData(w, map[string]string{"id": "new-id"})
}

func Show{{.ResourceName}}Handler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	// TODO: Implement
	response.NewSuccessData(w, map[string]string{"id": "123"})
}

func Update{{.ResourceName}}Handler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	// TODO: Implement
	response.NewSuccessMessage(w, "{{.ResourceName}} updated successfully")
}

func Delete{{.ResourceName}}Handler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	// TODO: Implement
	response.NewSuccessMessage(w, "{{.ResourceName}} deleted successfully")
}

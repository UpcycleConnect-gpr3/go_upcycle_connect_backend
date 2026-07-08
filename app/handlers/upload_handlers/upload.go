package upload_handlers

import (
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/response"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const uploadDir = "uploads"

var allowedExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

func StoreUploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.NewErrorMessage(w, "Fichier trop volumineux ou invalide", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.NewErrorMessage(w, "Champ 'file' manquant", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExt[ext] {
		response.NewErrorMessage(w, "Format d'image non supporté", http.StatusBadRequest)
		return
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		response.NewErrorMessage(w, "Stockage indisponible", http.StatusInternalServerError)
		return
	}

	filename := uuid.New().String() + ext
	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		response.NewErrorMessage(w, "Impossible d'enregistrer le fichier", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		response.NewErrorMessage(w, "Impossible d'enregistrer le fichier", http.StatusInternalServerError)
		return
	}

	response.NewSuccessData(w, map[string]string{"path": "/" + uploadDir + "/" + filename})
}

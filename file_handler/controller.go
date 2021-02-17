package file_handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/storage-blobs-go-quickstart/errors"
	"github.com/Azure/storage-blobs-go-quickstart/usecase"
	"github.com/labstack/echo"
)

func Upload(c echo.Context) error {
	err := c.Request().ParseMultipartForm(10 << 20)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	clienteID := strings.ReplaceAll(strings.ToLower(c.Request().FormValue("clienteID")), " ", "_")
	bcID := strings.ReplaceAll(strings.ToLower(c.Request().FormValue("bancoClienteID")), " ", "_")
	pacienteID := strings.ReplaceAll(strings.ToLower(c.Request().FormValue("pacienteID")), " ", "_")
	fileType := strings.ReplaceAll(strings.ToLower(c.Request().FormValue("fileType")), " ", "_")
	if clienteID == "" || bcID == "" || pacienteID == "" || fileType == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": errors.ErrInvalidPayload.Error(),
		})
	}

	if fileType != "anexo" && fileType != "foto" && fileType != "video" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": errors.ErrInvalidPayload.Error(),
		})
	}

	// falta verificar se IDs são válidos e comparar com os IDs do token

	fileHeader, err := c.FormFile("myFile")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	if err = usecase.Upload(fileHeader, clienteID, bcID, pacienteID, fileType); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func GetURL(c echo.Context) error {
	clienteID := c.Param("cliente_id")
	bcID := c.Param("bc_id")
	pacienteID := c.Param("paciente_id")
	fileType := c.Param("filetype")
	filename := c.Param("filename")

	log.Printf("clienteID: %s, bcID: %s, pacienteID: %s, fileType: %s, filename: %s", clienteID, bcID, pacienteID, fileType, filename)

	if clienteID == "" || bcID == "" || pacienteID == "" || fileType == "" || filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": errors.ErrInvalidPayload.Error(),
		})
	}

	if fileType != "anexo" && fileType != "foto" && fileType != "video" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": errors.ErrInvalidPayload.Error(),
		})
	}

	url, err := usecase.GetURL(clienteID, bcID, pacienteID, fileType, filename)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	getURLResponse := &getURLResponse{Url: url}
	return c.JSON(http.StatusOK, getURLResponse)
}

func Download(c echo.Context) error {

	clienteID := c.Param("cliente_id")
	bcID := c.Param("bc_id")
	pacienteID := c.Param("paciente_id")
	fileType := c.Param("file_type")
	filename := c.Param("filename")
	mode := c.QueryParam("mode")

	log.Printf("clienteID: %s, bcID: %s, pacienteID: %s, fileType: %s, filename: %s, mode: %s", clienteID, bcID, pacienteID, fileType, filename, mode)

	if clienteID == "" || bcID == "" || pacienteID == "" || fileType == "" || filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": errors.ErrInvalidPayload.Error(),
		})
	}

	if fileType != "anexo" && fileType != "foto" && fileType != "video" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": errors.ErrInvalidPayload.Error(),
		})
	}

	if mode != "attachment" && mode != "inline" && mode != "download" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": errors.ErrInvalidPayload.Error(),
		})
	}

	file, err := usecase.Download(clienteID, bcID, pacienteID, fileType, filename)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	defer os.Remove(file)
	name := filepath.Base(file)

	log.Println(file, name)
	switch mode {
	case "attachment": // força download no navegador
		return c.Attachment(file, name)
	case "inline": // abre no navegador
		return c.Inline(file, name)
	default: // abre no navegador ou força o download pra vc colocar o nome
		return c.File(file)
	}
}

type getURLResponse struct {
	Url string `json:"url"`
}

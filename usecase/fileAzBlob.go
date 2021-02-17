package usecase

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/storage-blobs-go-quickstart/errors"
)

type azureConnection struct {
	key          string
	account      string
	containerURL azblob.ContainerURL
	pipeline     pipeline.Pipeline
	ctx          context.Context
}

func newAzureConnection() (*azureConnection, error) {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	if len(accountName) == 0 || len(accountKey) == 0 {
		return &azureConnection{}, errors.ErrInvalidCredentials
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return &azureConnection{}, errors.ErrInvalidCredentials
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	ctx := context.Background()

	return &azureConnection{
		key:      accountKey,
		account:  accountName,
		pipeline: p,
		ctx:      ctx,
	}, nil
}

func GetURL(clienteID, bcID, pacienteID, fileType, filename string) (string, error) {
	conn, err := newAzureConnection()
	if err != nil {
		return "", err
	}

	URL, err := conn.get(clienteID, bcID, pacienteID, fileType, filename)
	if err != nil {
		return "", err
	}

	return URL, nil
}

func Download(clienteID, bcID, pacienteID, fileType, filename string) (string, error) {
	conn, err := newAzureConnection()
	if err != nil {
		return "", err
	}

	file, err := conn.download(clienteID, bcID, pacienteID, fileType, filename)
	if err != nil {
		return "", err
	}

	return file, nil
}

func (az *azureConnection) download(clienteID, bcID, pacienteID, fileType, filename string) (string, error) {
	basePath := fmt.Sprintf("https://%s.blob.core.windows.net/%s", az.account, clienteID)
	URL, _ := url.Parse(basePath)
	containerURL := azblob.NewContainerURL(*URL, az.pipeline)
	newBlob := fmt.Sprintf("%s/%s/%s/%s", bcID, pacienteID, fileType, filename)

	blobURL := containerURL.NewBlobURL(newBlob)

	downloadResponse, err := blobURL.Download(az.ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return "", err
	}
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	if err != nil {
		return "", err
	}

	tempFile, err := os.Create("./temp/" + filename)
	if err != nil {
		return "", err
	}
	defer func() {
		log.Println("fechando arquivo...")
		tempFile.Close()
	}()

	_, err = io.Copy(tempFile, &downloadedData)
	if err != nil {
		return "", err
	}

	return "./temp/" + filename, nil
}

func (az *azureConnection) get(clienteID, bcID, pacienteID, fileType, filename string) (string, error) {
	basePath := fmt.Sprintf("https://%s.blob.core.windows.net/%s", az.account, clienteID)
	URL, _ := url.Parse(basePath)
	containerURL := azblob.NewContainerURL(*URL, az.pipeline)
	newBlob := fmt.Sprintf("%s/%s/%s/%s", bcID, pacienteID, fileType, filename)

	blobURL := containerURL.NewBlobURL(newBlob)

	_, err := blobURL.Download(az.ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false) // só testa se blob está lá, não faz o download
	if err != nil {
		return "", err
	}
	// bodyStream := download.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20}) // se quiser o conteúdo do blob
	// downloadedData := bytes.Buffer{}
	// _, err = downloadedData.ReadFrom(bodyStream)
	// if err != nil {
	// 	return "", err
	// }

	validURL := blobURL.URL()

	return validURL.String(), nil
}

func Upload(fileHeader *multipart.FileHeader, clienteID, bcID, pacienteID, fileType string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	filename := html.EscapeString(strings.TrimSpace(strings.ReplaceAll(fileHeader.Filename, " ", "_")))

	conn, err := newAzureConnection()
	if err != nil {
		return err
	}

	if err = conn.save(file, clienteID, bcID, pacienteID, fileType, filename); err != nil {
		return err
	}

	return nil
}

func (az *azureConnection) save(file multipart.File, clienteID, bcID, pacienteID, fileType, filename string) error {
	log.Println("Azure working...")

	if err := az.ensureContainer(clienteID); err != nil {
		return err
	}

	newBlob := fmt.Sprintf("%s/%s/%s/%s", bcID, pacienteID, fileType, filename)
	blobURL := az.containerURL.NewBlockBlobURL(newBlob)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if _, err := azblob.UploadBufferToBlockBlob(az.ctx, fileBytes, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16}); err != nil {
		return err
	}

	return nil
}

func (az *azureConnection) ensureContainer(containerName string) error {

	basePath := fmt.Sprintf("https://%s.blob.core.windows.net/%s", az.account, containerName)
	URL, _ := url.Parse(basePath)
	containerURL := azblob.NewContainerURL(*URL, az.pipeline)

	_, err := containerURL.Create(az.ctx, azblob.Metadata{}, azblob.PublicAccessContainer)
	if err != nil {
		if serr, ok := err.(azblob.StorageError); ok { // This error is a Service-specific
			switch serr.ServiceCode() { // Compare serviceCode to ServiceCodeXxx constants
			case azblob.ServiceCodeContainerAlreadyExists:
				log.Println("Container already exists")
			default:
				return err
			}
		} else {
			return err
		}

	}

	az.containerURL = containerURL
	return nil
}

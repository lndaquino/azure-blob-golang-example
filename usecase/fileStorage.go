package usecase

// ESTUDANDO COMO ENVIAR EM BLOCKS PRO AZURE, POR ENQUANTO ENVIA ATÉ 256MB PELO AZBLOB EM THREADS PARALELAS
// https://stackoverflow.com/questions/43187362/golang-processing-images-via-multipart-and-streaming-to-azure
import (
	"html"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/storage-blobs-go-quickstart/errors"
)

type azureStorage struct {
	client storage.Client
}

func newAzureStorage() (*azureStorage, error) {
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	if len(accountName) == 0 || len(accountKey) == 0 {
		return &azureStorage{}, errors.ErrInvalidCredentials
	}

	azClient, err := storage.NewBasicClient(accountName, accountKey)
	if err != nil {
		return &azureStorage{}, errors.ErrInvalidCredentials
	}

	return &azureStorage{
		client: azClient,
	}, nil
}

func (az *azureStorage) saveInBlocks(file multipart.File, clienteID, bcID, pacienteID, fileType, filename string) error {
	const BLOB_LENGTH_LIMITS uint64 = 64 * 1024 * 1024

	if err := az.ensureContainerBlock(clienteID); err != nil {
		return err
	}

	// newBlob := fmt.Sprintf("%s/%s/%s/%s", bcID, pacienteID, fileType, filename)
	// blobClient := storage.Blob{}

	// blobClient.PutBlock()
	// r := bufio.NewReader(file)
	// b := make([]byte, BLOB_LENGTH_LIMITS)
	// i := 0
	// var blocks []azblob.Block
	// for {
	// 	_, err := r.Read(b)
	// 	if err != nil {
	// 		if err != io.EOF {
	// 			return err
	// 		}
	// 		break
	// 	}
	// 	blockID := base64.StdEncoding.EncodeToString(b)
	// 	block := azblob.Block{
	// 		Name: blockID,
	// 		Size: int32(len(b)),
	// 	}
	// 	blocks = append(blocks, block)
	// }
	return nil
}

func (az *azureStorage) ensureContainerBlock(containerName string) error {
	// az.client
	var container storage.Container

	created, err := container.CreateIfNotExists(&storage.CreateContainerOptions{})
	if err != nil {
		return err
	}
	if created {
		log.Println("container criado...")
	} else {
		log.Println("container já existente...")
	}

	return nil
}

func UploadBlocks(fileHeader *multipart.FileHeader, clienteID, bcID, pacienteID, fileType string) error {
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

	conn, err := newAzureStorage()
	if err != nil {
		return err
	}

	if err = conn.saveInBlocks(file, clienteID, bcID, pacienteID, fileType, filename); err != nil {
		return err
	}

	return nil
}

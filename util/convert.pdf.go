package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	tokenURL   = "https://api.aspose.cloud/connect/token"
	baseAPIURL = "https://api.aspose.cloud/v4.0/words"
)

type AsposeClient struct {
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client
}

func NewAsposeClient(clientID, clientSecret string) *AsposeClient {
	return &AsposeClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (a *AsposeClient) GetToken(ctx context.Context) (string, error) {

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", a.ClientID)
	form.Set("client_secret", a.ClientSecret)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		tokenURL,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error obteniendo token: %s", string(body))
	}

	var tr tokenResponse

	if err := json.Unmarshal(body, &tr); err != nil {
		return "", err
	}

	if tr.AccessToken == "" {
		return "", fmt.Errorf("token vacío")
	}

	return tr.AccessToken, nil
}

func (a *AsposeClient) UploadFile(ctx context.Context, token, localFile, remoteName string) error {

	file, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	u := fmt.Sprintf("%s/storage/file/%s", baseAPIURL, url.PathEscape(remoteName))

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u, file)
	if err != nil {
		return err
	}

	req.ContentLength = info.Size()
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error subiendo archivo: %s", string(body))
	}

	return nil
}

type saveAsRequest struct {
	SaveFormat string `json:"SaveFormat"`
	FileName   string `json:"FileName"`
}

func (a *AsposeClient) ConvertToPDF(
	ctx context.Context,
	token string,
	remoteDocx string,
	remotePDF string,
	localPDF string,
) error {

	reqBody := saveAsRequest{
		SaveFormat: "pdf",
		FileName:   remotePDF,
	}

	data, _ := json.Marshal(reqBody)

	u := fmt.Sprintf("%s/%s/saveAs", baseAPIURL, url.PathEscape(remoteDocx))

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		u,
		bytes.NewBuffer(data),
	)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", mime.TypeByExtension(".json"))

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error convirtiendo archivo: %s", string(body))
	}

	return a.DownloadFile(ctx, token, remotePDF, localPDF)
}

func (a *AsposeClient) DownloadFile(
	ctx context.Context,
	token string,
	remoteFile string,
	localFile string,
) error {

	u := fmt.Sprintf("%s/storage/file/%s", baseAPIURL, url.PathEscape(remoteFile))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error descargando PDF: %s", string(body))
	}

	if err := os.MkdirAll(filepath.Dir(localFile), 0755); err != nil {
		return err
	}

	out, err := os.Create(localFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}

type DrawingObjectInsert struct {
	RelativeHorizontalPosition string  `json:"RelativeHorizontalPosition,omitempty"`
	RelativeVerticalPosition   string  `json:"RelativeVerticalPosition,omitempty"`
	Width                      float64 `json:"Width,omitempty"`
	Height                     float64 `json:"Height,omitempty"`
	WrapType                   string  `json:"WrapType,omitempty"`
}

type ReplaceTextRequest struct {
	OldValue    string `json:"OldValue"`
	NewValue    string `json:"NewValue"`
	IsMatchCase bool   `json:"IsMatchCase"`
}

type SearchResponse struct {
	SearchResults struct {
		List []struct {
			RangeStart struct {
				Node struct {
					NodeId string `json:"NodeId"`
				} `json:"Node"`
			} `json:"RangeStart"`
		} `json:"List"`
	} `json:"SearchResults"`
}

func (a *AsposeClient) SearchText(ctx context.Context, token, remoteDocx, pattern string) (string, error) {
	u := fmt.Sprintf("%s/%s/search?pattern=%s", baseAPIURL, url.PathEscape(remoteDocx), url.QueryEscape(pattern))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error buscando texto '%s': %s", pattern, string(body))
	}

	var sr SearchResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &sr); err != nil {
		return "", err
	}

	if len(sr.SearchResults.List) > 0 {
		nodeId := sr.SearchResults.List[0].RangeStart.Node.NodeId
		return nodeId, nil
	}

	return "", nil
}

func (a *AsposeClient) ReplaceText(ctx context.Context, token, remoteDocx, oldValue, newValue string) error {
	u := fmt.Sprintf("%s/%s/replaceText", baseAPIURL, url.PathEscape(remoteDocx))

	payload := ReplaceTextRequest{
		OldValue:    oldValue,
		NewValue:    newValue,
		IsMatchCase: false,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error reemplazando texto '%s': %s", oldValue, string(body))
	}

	return nil
}

func (a *AsposeClient) InsertImageAtNode(
	ctx context.Context,
	token string,
	remoteDocx string,
	nodePath string,
	localImagePath string,
) error {
	paragraphPath := nodePath
	if parts := strings.Split(nodePath, "."); len(parts) > 4 {
		paragraphPath = strings.Join(parts[:4], ".")
	}

	var u string
	if strings.HasPrefix(paragraphPath, "sections/") || strings.HasPrefix(paragraphPath, "paragraphs/") {
		u = fmt.Sprintf("%s/%s/%s/drawingObjects", baseAPIURL, url.PathEscape(remoteDocx), paragraphPath)
	} else if paragraphPath != "" {
		u = fmt.Sprintf("%s/%s/paragraphs/%s/drawingObjects", baseAPIURL, url.PathEscape(remoteDocx), url.PathEscape(paragraphPath))
	} else {
		u = fmt.Sprintf("%s/%s/drawingObjects", baseAPIURL, url.PathEscape(remoteDocx))
	}

	file, err := os.Open(localImagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	drawingObj := DrawingObjectInsert{
		WrapType: "Inline",
		Width:    200,
		Height:   100,
	}
	drawingData, _ := json.Marshal(drawingObj)

	part, _ := writer.CreateFormField("DrawingObject")
	part.Write(drawingData)

	filePart, err := writer.CreateFormFile("File", filepath.Base(localImagePath))
	if err != nil {
		return err
	}
	io.Copy(filePart, file)
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error insertando imagen en nodo %s: %s", nodePath, string(respBody))
	}

	return nil
}

func (a *AsposeClient) InsertImageAtBookmark(
	ctx context.Context,
	token string,
	remoteDocx string,
	bookmarkName string,
	localImagePath string,
) error {
	u := fmt.Sprintf("%s/%s/bookmarks/%s/drawingObjects", baseAPIURL, url.PathEscape(remoteDocx), url.PathEscape(bookmarkName))

	file, err := os.Open(localImagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add JSON part for DrawingObject
	drawingObj := DrawingObjectInsert{
		WrapType: "Inline",
		Width:    200,
		Height:   100,
	}
	drawingData, _ := json.Marshal(drawingObj)

	part, _ := writer.CreateFormField("DrawingObject")
	part.Write(drawingData)

	// Add file part
	filePart, err := writer.CreateFormFile("File", filepath.Base(localImagePath))
	if err != nil {
		return err
	}
	io.Copy(filePart, file)

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error insertando imagen en marcador: %s", string(respBody))
	}

	return nil
}

func (a *AsposeClient) InsertImageAtTextOrBookmark(
	ctx context.Context,
	token string,
	remoteDocx string,
	localImagePath string,
) error {
	// 1. Probar patrones de texto plano
	rawPatterns := []string{"{{firma}}", "{{Firma}}", "{{FIRMA}}", "{firma}", "{Firma}", "{FIRMA}", "[FIRMA]", "[firma]"}

	for _, raw := range rawPatterns {
		escapedPattern := regexp.QuoteMeta(raw)
		for _, pattern := range []string{escapedPattern, raw} {
			nodeId, err := a.SearchText(ctx, token, remoteDocx, pattern)
			if err == nil && nodeId != "" {
				// Borrar el texto de la etiqueta
				_ = a.ReplaceText(ctx, token, remoteDocx, raw, "")

				// Insertar la imagen en la ubicación/nodo encontrado
				if err := a.InsertImageAtNode(ctx, token, remoteDocx, nodeId, localImagePath); err == nil {
					return nil
				}
			}
		}
	}

	// 2. Si no se encontró texto plano, intentar marcadores "Firma" o "firma"
	bookmarks := []string{"Firma", "firma"}
	for _, bm := range bookmarks {
		if err := a.InsertImageAtBookmark(ctx, token, remoteDocx, bm, localImagePath); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no se encontró la etiqueta de firma '{{firma}}' ni el marcador 'Firma' en la plantilla DOCX")
}


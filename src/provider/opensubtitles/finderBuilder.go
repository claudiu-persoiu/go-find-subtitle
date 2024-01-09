package opensubtitles

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/claudiu-persoiu/go-find-subtitle/src/provider/opensubtitles/hash"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Finder struct {
	config Config
	token  string
	client *Client
}

type Config struct {
	User      string `json:"user"`
	Pass      string `json:"pass"`
	Key       string `json:"key"`
	Languages string `json:"languages"`
}

const URL = "https://api.opensubtitles.com"
const timeout = 10

func NewFinder(credentials Config, client *Client) *Finder {
	if credentials.Key == "" {
		log.Fatalln("No Opensutitles credentials provided, please see documentation at: https://github.com/claudiu-persoiu/go-find-subtitle")
	}

	f := &Finder{
		config: credentials,
		client: client,
	}

	return f
}

func (f *Finder) Find(path string) (bool, error) {

	if f.token == "" {
		err := f.login()
		if err != nil {
			return false, err
		}
	}

	fileId, err := f.search(path)
	if err != nil {
		return false, err
	}

	if fileId == "" {
		return false, err
	}

	url, err := f.download(fileId)
	if err != nil {
		return false, err
	}

	err = f.downloadFile(url, path)
	if err != nil {
		return false, err
	}

	return true, nil
}

type loginResponse struct {
	Token   string `json:"token"`
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

func (f *Finder) login() error {
	body := []byte(`{
		"username": "` + f.config.User + `",
		"password": "` + f.config.Pass + `"
	}`)

	resp, err := f.client.Request(http.MethodPost, "/api/v1/login", body, "")

	response := &loginResponse{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(response)
	if err != nil {
		return fmt.Errorf("could not parse response: %s", err)
	}

	if response.Token == "" {
		return errors.New("Could not create token: " + response.Message)
	}

	f.token = response.Token

	return nil
}

type searchResponse struct {
	Data []searchResultItem `json:"data,omitempty"`
}

type searchResultItem struct {
	Attributes searchResultAttributes `json:"attributes,omitempty"`
}

type searchResultAttributes struct {
	Language string             `json:"language"`
	Files    []searchResultFile `json:"files"`
}

type searchResultFile struct {
	Id   int    `json:"file_id"`
	Name string `json:"name"`
}

func (f *Finder) search(path string) (fileId string, err error) {
	var hashSearch string
	hash, err := hash.BuildHash(path)
	if err != nil {
		fmt.Println(fmt.Errorf("could not create hash: %s for file %s", err, path))
	} else {
		hashSearch = "&moviehash=" + hash
	}

	resp, err := f.client.Request(http.MethodGet, "/api/v1/subtitles?&languages="+f.config.Languages+"&"+hashSearch+"&query="+filepath.Base(path), nil, f.token)

	response := &searchResponse{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(response)
	if err != nil {
		return "", fmt.Errorf("could not parse response: %s", err)
	}

	if response != nil && len(response.Data) > 0 && len(response.Data[0].Attributes.Files) > 0 {
		fileId = strconv.Itoa(response.Data[0].Attributes.Files[0].Id)
	}

	return
}

type downloadResponse struct {
	Link string `json:"link"`
}

func (f *Finder) download(fileId string) (url string, err error) {

	body := []byte(`{
		"file_id": "` + fileId + `"
	}`)

	resp, err := f.client.Request(http.MethodPost, "/api/v1/download", body, f.token)

	response := &downloadResponse{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(response)
	if err != nil {
		return url, fmt.Errorf("could not parse response: %s", err)
	}

	return response.Link, nil
}

func (f *Finder) downloadFile(url, path string) error {
	out, err := os.Create(path[:len(path)-4] + ".srt")
	if err != nil {
		return fmt.Errorf("could not create subtitle file: %s", err)
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not download file: %s", err)
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}

func (f *Finder) logout() error {
	if f.token == "" {
		return nil
	}

	_, err := f.client.Request(http.MethodDelete, "/api/v1/logout", nil, f.token)

	if err != nil {
		return err
	}

	return nil
}

func (f *Finder) Close() error {
	return f.logout()
}

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"unbewohnte.xyz/Unbewohnte/gochat/api"
	"unbewohnte.xyz/Unbewohnte/gochat/log"
)

// Handle incoming attachments
func (s *Server) handleAttachments(w http.ResponseWriter, req *http.Request) {
	userAuthHeader := api.GetUserAuthHeaderData(req)
	if !s.db.DoesUserExist(userAuthHeader.Name) {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	userDB, _ := s.db.GetUser(userAuthHeader.Name)
	if userDB.SecretHash != userAuthHeader.SecretHash {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	// accept incoming file
	file, header, err := req.FormFile(api.AttachmentFormPostKey)
	if err != nil {
		log.Error("could not get attached file: %s", err)
		http.Error(w, "Error getting attached file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if uint64(header.Size) > api.MaxAttachmentSize {
		http.Error(w, "Too big file", http.StatusBadRequest)
		return
	}

	localFilename := fmt.Sprintf("%s_%d_%s", userDB.Name, header.Size, header.Filename)
	localFile, err := os.Create(filepath.Join(s.workingDir, attachmentsDirName, localFilename))
	if err != nil {
		log.Error("could not create local attachment file: %s", err)
		http.Error(w, "Could not create file", http.StatusInternalServerError)
		return
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, file)
	if err != nil {
		log.Error("could not copy attachment file contents: %s", err)
		http.Error(w, "Could not copy file contents", http.StatusInternalServerError)
		return
	}

	// send partial URL pointing to the file
	url := api.PartialAttachmentURL{
		URL: fmt.Sprintf("%s/%s", attachmentsDirName, localFilename),
	}
	urlJsonBytes, err := json.Marshal(&url)
	if err != nil {
		log.Error("could not marshal partial attachment URL for \"%s\": %s", localFilename, err)
		http.Error(w, "Error constructing a request", http.StatusInternalServerError)
		return
	}

	w.Write(urlJsonBytes)
}

// Remove the oldest attachments when the memory limit is exceeded
func manageAttachmentsStorage(attachmentsDirPath string, sizeLimit uint64, checkDelay time.Duration) {
	for {
		dirEntries, err := os.ReadDir(attachmentsDirPath)
		if err != nil {
			log.Error("error reading attachments directory: %s", err)
		}

		var dirSize uint64 = 0
		var oldestAttachmentModTime time.Time = time.Now()
		var oldestAttachmentPath string = ""
		var oldestAttachmentSize uint64 = 0
		if dirEntries != nil {
			for _, entry := range dirEntries {
				entryInfo, err := entry.Info()
				if err != nil {
					continue
				}

				entrySize := entryInfo.Size()
				dirSize += uint64(entrySize)

				entryModTime := entryInfo.ModTime()
				if entryModTime.Before(oldestAttachmentModTime) {
					oldestAttachmentModTime = entryModTime
					oldestAttachmentPath = filepath.Join(attachmentsDirPath, entry.Name())
					oldestAttachmentSize = uint64(entrySize)
				}
			}

			if dirSize > sizeLimit {
				// A cleanup !
				os.Remove(oldestAttachmentPath)
				log.Info(
					"removed %s during attachments storage management. Cleared %d bytes",
					oldestAttachmentPath, oldestAttachmentSize)
			}
		}

		time.Sleep(checkDelay)
	}
}

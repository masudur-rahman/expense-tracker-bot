package google

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	dbFileName   = "expense-tracker.db"
	dbFolderName = ".expense-tracker"
)

func getDriveService() (*drive.Service, error) {
	creds, err := google.FindDefaultCredentials(context.Background(), drive.DriveScope)
	if err != nil {
		return nil, err
	}

	return drive.NewService(context.Background(), option.WithCredentials(creds))
}

func readFileFromDrive(svc *drive.Service, upstreamFilePath string) ([]byte, error) {
	ff := strings.Split(upstreamFilePath, "/")

	folderID, err := getFolderID(svc, ff[0])
	if err != nil {
		return nil, err
	}

	fileQuery := fmt.Sprintf("name='%s' and '%s' in parents", ff[1], folderID)
	fileList, err := svc.Files.List().Q(fileQuery).Do()
	if err != nil {
		return nil, err
	}

	if len(fileList.Files) == 0 {
		return nil, fmt.Errorf("file '%s' not found in folder '%s'", ff[1], ff[0])
	}

	fileID := fileList.Files[0].Id

	resp, err := svc.Files.Get(fileID).Download()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func uploadFileToDrive(svc *drive.Service, upstreamFilePath string, localFilePath string) error {
	ff := strings.Split(upstreamFilePath, "/")
	folderID, err := getFolderID(svc, ff[0])
	if err != nil {
		return err
	}

	fileMetadata := &drive.File{
		Name:     filepath.Base(ff[1]),
		Parents:  []string{folderID},
		MimeType: "application/octet-stream",
	}

	file, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = svc.Files.Create(fileMetadata).Media(file).Do()
	return err
}

func getFolderID(svc *drive.Service, folderName string) (string, error) {
	folderQuery := fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder'", folderName)
	folderList, err := svc.Files.List().Q(folderQuery).Do()
	if err != nil {
		return "", err
	}

	if len(folderList.Files) == 0 {
		parentFolder := &drive.File{
			Name:     folderName,
			MimeType: "application/vnd.google-apps.folder",
		}
		parentFolder, err = svc.Files.Create(parentFolder).Do()
		if err != nil {
			return "", err
		}
		return parentFolder.Id, nil
	}

	return folderList.Files[0].Id, nil
}

func SyncDatabaseFromDrive() error {
	svc, err := getDriveService()
	if err != nil {
		return err
	}

	upstreamDBPath := fmt.Sprintf("%s/%s", dbFolderName, dbFileName)
	data, err := readFileFromDrive(svc, upstreamDBPath)
	if err != nil {
		return err
	}

	return os.WriteFile(dbFileName, data, 0666)
}

func SyncDatabaseToDrive() error {
	svc, err := getDriveService()
	if err != nil {
		return err
	}

	upstreamDBPath := fmt.Sprintf("%s/%s", dbFolderName, dbFileName)
	return uploadFileToDrive(svc, upstreamDBPath, dbFileName)
}

func SyncDatabaseToDrivePeriodically(interval time.Duration) {
	if interval == time.Duration(0) {
		interval = time.Hour
	}

	if err := SyncDatabaseToDrive(); err != nil {
		logr.DefaultLogger.Errorw("Sync database to drive failed", "error", err.Error())
		return
	}
	logr.DefaultLogger.Infof("Sqlite database synced to google drive")

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			if err := SyncDatabaseToDrive(); err != nil {
				logr.DefaultLogger.Errorw("Sync database to drive failed", "error", err.Error())
				return
			}
			logr.DefaultLogger.Infof("Sqlite database synced to google drive")
		}
	}
}
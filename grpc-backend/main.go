package main

import (
	v1 "backend/gen"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"regexp"
	"strings"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/juli0n21/go-osu-parser/parser"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

	// Run gRPC + grpc-gateway servers
	if err := runGrpcAndGateway(s, port); err != nil {
		log.Fatalf("Failed to run servers: %v", err)
	}
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
	url := GetEnv("PROXY_URL", "https://proxy.illegalesachen.download/settings")

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
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	req.AddCookie(&http.Cookie{
		Name:  "session_cookie",
		Value: cookie,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error in request: %s", resp.Status)
	}
	return nil
}

func runGrpcAndGateway(s *Server, port string) error {
	grpcPort := ":9090" // gRPC server port
	httpPort := port    // REST gateway port (e.g. ":8080")

	grpcLis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", grpcPort, err)
	}
	grpcServer := grpc.NewServer()
	v1.RegisterMusicBackendServer(grpcServer, s) // Register your service implementation

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = v1.RegisterMusicBackendHandlerFromEndpoint(ctx, gwMux, grpcPort, opts)
	if err != nil {
		return fmt.Errorf("failed to register grpc-gateway: %w", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/api/v1/", gwMux)

	mux.HandleFunc("/files", s.songFile)

	fileServer := http.FileServer(http.Dir("gen/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fileServer))

	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 2)

	go func() {
		log.Printf("Starting gRPC server on %s", grpcPort)
		errChan <- grpcServer.Serve(grpcLis)
	}()

	go func() {
		log.Printf("Starting HTTP gateway server on %s", httpPort)
		errChan <- httpServer.ListenAndServe()
	}()

	select {
	case <-stop:
		log.Println("Shutting down servers...")
		grpcServer.GracefulStop()
		httpServer.Shutdown(ctx)
		return nil
	case err := <-errChan:
		return err
	}
}

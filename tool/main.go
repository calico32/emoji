package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var CACHE_FILE string
var CONFIG_FILE string
var LOG_FILE string

type Cache struct {
	UpdatedAt int64    `json:"updated_at"`
	Emojis    []string `json:"emojis"`
}

type Result struct {
	Status string   `json:"status"`
	Files  []string `json:"files"`
}

func init() {
	CONFIG_FILE = os.Getenv("EMOJITOOL_CONFIG_FILE")
	if CONFIG_FILE == "" {
		CONFIG_FILE = os.Getenv("HOME") + "/.config/emojitool"
	}
	godotenv.Load(CONFIG_FILE)

	CACHE_FILE = os.Getenv("EMOJITOOL_CACHE_FILE")
	if CACHE_FILE == "" {
		CACHE_FILE = os.Getenv("HOME") + "/.cache/emojitool"
	}

	LOG_FILE = os.Getenv("EMOJITOOL_LOG_FILE")
	if LOG_FILE == "" {
		LOG_FILE = os.Getenv("HOME") + "/.cache/emojitool.log"
	}

	f, err := os.OpenFile(LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}

	log.SetOutput(f)
}

func main() {
	key := os.Getenv("MANAGE_KEY")
	if key == "" {
		log.Println("MANAGE_KEY not set!")
		os.Exit(1)
	}

	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		log.Println("BASE_URL not set!")
	}

	var cache Cache

	if f, err := os.Open(CACHE_FILE); err == nil {
		defer f.Close()
		json.NewDecoder(f).Decode(&cache)
	}

	if len(cache.Emojis) != 0 {
		showPicker(baseUrl, cache)
	}

	if time.Now().Unix()-cache.UpdatedAt > 3600 {
		refreshCache(baseUrl, key)
	}
}

func showPicker(baseUrl string, cache Cache) {
	windowCmd := exec.Command("xdotool", "getactivewindow")
	window, err := windowCmd.Output()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("showing picker")
	rofiArgs := []string{}
	if os.Getenv("ROFI_ARGS") != "" {
		rofiArgs = append(rofiArgs, strings.Split(os.Getenv("ROFI_ARGS"), " ")...)
	}
	rofiArgs = append(rofiArgs, "-sep", "|", "-dmenu", "-p", "emoji", "-i", "-no-custom", "-w", strings.TrimSpace(string(window)))

	cmd := exec.Command("rofi", rofiArgs...)
	cmd.Stdin = strings.NewReader(strings.Join(cache.Emojis, "|"))

	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	name := strings.TrimSpace(string(out))

	if len(name) == 0 {
		return
	}

	// type with xdotool
	exec.Command("xdotool", "type", baseUrl+"/"+name).Run()
}

func refreshCache(baseUrl, key string) {
	log.Println("refreshing cache")
	req, err := http.NewRequest(http.MethodOptions, baseUrl, nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", "Bearer "+key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var result Result
	json.NewDecoder(res.Body).Decode(&result)

	cache := Cache{
		UpdatedAt: time.Now().Unix(),
		Emojis:    result.Files,
	}

	if f, err := os.Create(CACHE_FILE); err == nil {
		defer f.Close()
		json.NewEncoder(f).Encode(cache)
	} else {
		log.Println(err)
		os.Exit(1)
	}
}

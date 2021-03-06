package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"os"
	"os/exec"
	"runtime"
	"time"
	
	"github.com/spf13/cobra"
	"github.com/fsnotify/fsnotify"
)

var (
	state sync.Mutex
	_cmd *exec.Cmd
	appname string
	eventTime = make(map[string]int64)
	scheduleTime        time.Time
	watchExts = []string{
		".go",
		".conf",
	}
	ignoredFilesRegExps = []string{
		`.#(\w+).go`,
		`.(\w+).go.swp`,
		`(\w+).go~`,
		`(\w+).tmp`,
		`commentsRouter_controllers.go`,
		`\/features\/`,
	}
)
var started = make(chan bool)
var exit = make(chan bool)

var buildSource = "main.go"

func readAppDirectories(directory string, paths *[]string) {
	fileInfos, err := ioutil.ReadDir(directory)
	if err != nil {
		return
	}
	
	useDirectory := false
	for _, fileInfo := range fileInfos {
		if strings.HasSuffix(fileInfo.Name(), "docs") {
			continue
		}
		if strings.HasSuffix(fileInfo.Name(), "swagger") {
			continue
		}
		
		if strings.HasSuffix(fileInfo.Name(), "vendor") {
			continue
		}
		
		if fileInfo.IsDir() && fileInfo.Name()[0] != '.' {
			readAppDirectories(directory+"/"+fileInfo.Name(), paths)
			continue
		}
		
		if useDirectory {
			continue
		}
		
		if path.Ext(fileInfo.Name()) == ".go" || path.Ext(fileInfo.Name()) == ".conf" {
			*paths = append(*paths, directory)
			useDirectory = true
		}
	}
}

// getFileModTime returns unix timestamp of `os.File.ModTime` for the given path.
func getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		log.Printf("[ERROR] Failed to open file on '%s': %s", path, err)
		return time.Now().Unix()
	}
	defer f.Close()
	
	fi, err := f.Stat()
	if err != nil {
		log.Printf("Failed to get file stats: %s", err)
		return time.Now().Unix()
	}
	
	return fi.ModTime().Unix()
}

// NewWatcher starts an fsnotify Watcher on the specified paths
func newWatcher(paths []string, buildSource string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %s", err)
	}
	
	go func() {
		for {
			select {
			case e := <-watcher.Events:
				isBuild := true
				
				// Skip ignored files
				if shouldIgnoreFile(e.Name) {
					continue
				}
				
				if !shouldWatchFileWithExtension(e.Name) {
					continue
				}
				
				mt := getFileModTime(e.Name)
				if t := eventTime[e.Name]; mt == t {
					log.Printf("Skipping: %s", e.String())
					isBuild = false
				}
				
				eventTime[e.Name] = mt
				
				if isBuild {
					log.Printf("Event fired: %s", e)
					go func() {
						// Wait 1s before autobuild until there is no file change.
						scheduleTime = time.Now().Add(1 * time.Second)
						time.Sleep(time.Until(scheduleTime))
						autoBuild(buildSource)
					}()
				}
			case err := <-watcher.Errors:
				log.Printf("[ERROR] Watcher error: %s", err.Error()) // No need to exit here
			}
		}
	}()
	
	log.Println("Initializing watcher...")
	for _, path := range paths {
		log.Printf("Watching: %s", path)
		err = watcher.Add(path)
		if err != nil {
			log.Fatalf("Failed to watch directory: %s", err)
		}
	}
}

// shouldIgnoreFile ignores filenames generated by Emacs, Vim or SublimeText.
// It returns true if the file should be ignored, false otherwise.
func shouldIgnoreFile(filename string) bool {
	for _, regex := range ignoredFilesRegExps {
		r, err := regexp.Compile(regex)
		if err != nil {
			log.Fatalf("Could not compile regular expression: %s", err)
		}
		if r.MatchString(filename) {
			return true
		}
		continue
	}
	return false
}

// shouldWatchFileWithExtension returns true if the name of the file
// hash a suffix that should be watched.
func shouldWatchFileWithExtension(name string) bool {
	for _, s := range watchExts {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}

func autoBuild(file string) bool {
	state.Lock()
	defer state.Unlock()
	
	appName := appname
	if runtime.GOOS == "windows" {
		appName += ".exe"
	}
	
	args := []string{"build", "-v", "-mod=vendor"}
	args = append(args, "-o", appName)
	args = append(args, file)
	
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	bcmd := exec.Command("go", args...)
	bcmd.Env = append(os.Environ(), "GOGC=off")
	bcmd.Stderr = &stderr
	bcmd.Stdout = &stdout
	fmt.Println(fmt.Sprintf("run `%s` ...", bcmd.String()))
	err := bcmd.Run()
	if err != nil {
		fmt.Println(fmt.Errorf("[ERROR] Failed to build the application: \n%s", stderr.String()))
		return false
	} else {
		fmt.Println("[BUILD] SUCCESS: build success, now restart...")
		restart(appName)
		return true
	}
}

func restart(appname string) {
	log.Printf("restart %s", appname)
	kill()
	go start(appname)
}

func kill() {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("Kill recover: %s", e)
		}
	}()
	if _cmd != nil && _cmd.Process != nil {
		// Windows does not support Interrupt
		if runtime.GOOS == "windows" {
			_cmd.Process.Signal(os.Kill)
		} else {
			_cmd.Process.Signal(os.Interrupt)
		}
		
		ch := make(chan struct{}, 1)
		go func() {
			_cmd.Wait()
			ch <- struct{}{}
		}()
		
		select {
		case <-ch:
			return
		case <-time.After(10 * time.Second):
			log.Println("Timeout. Force kill cmd process")
			err := _cmd.Process.Kill()
			if err != nil {
				log.Println("[ERROR] Error while killing cmd process: %s", err)
			}
			return
		}
	}
}

func start(appname string) {
	log.Printf("Restarting '%s'...", appname)
	if !strings.Contains(appname, "./") {
		appname = "./" + appname
	}
	
	_cmd = exec.Command(appname)
	_cmd.Stdout = os.Stdout
	_cmd.Stderr = os.Stderr
	_cmd.Args = []string{appname}
	_cmd.Env = os.Environ()
	
	go _cmd.Run()
	log.Printf("[START] SUCCESS: '%s' is running...", appname)
	started <- true
}

var source string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the service by starting a local development server",
	
	Run: func(cmd *cobra.Command, args []string) {
		buildSource = "./main.go"
		if source != "" {
			buildSource = source
		}
		
		PrintBanner()
		curDir, _ := os.Getwd()
		
		//收集监控的文件列表
		var paths []string
		readAppDirectories(curDir, &paths)
		log.Printf("[WATCHER] collect total %d files", len(paths))
		
		// 监控文件
		newWatcher(paths, buildSource)
		
		//编译并启动
		if runtime.GOOS == "windows"{
			curDir = filepath.ToSlash(curDir)
		}
		service := path.Base(curDir)
		appname = service
		autoBuild(buildSource)
		for {
			<-exit
			runtime.Goexit()
		}
	},
}


func init() {
	runCmd.Flags().StringVarP(&source, "file", "f", "", "main file to be built")
	rootCmd.AddCommand(runCmd)
}
package overlap

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"path/filepath"
	"os"
)

func dirIsNotHidden(path string) bool {
	return len(path) > 1 && path[0:1] != "."
}

func IsCSV(filename string) bool {
	return filepath.Ext(filename) == ".csv"
}

func MonitorFolder(path string, done chan bool, files chan string) {
	fsnotify.NewWatcher()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panic(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op == fsnotify.Write && IsCSV(event.Name) {
					files <- event.Name
				}

			case err := <-watcher.Errors:
				log.Println("ошибка при мониторинге директории:", err)
			}
		}
	}()

	// if path is relative
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err == nil {
			err = watcher.Add(path)
		}

		if err != nil {
			log.Panic(err)
		}
	}


	err = filepath.Walk(path, func(filename string, f os.FileInfo, err error) error {
		if IsCSV(filename) {
			files <- filename
		}
		if f.IsDir() && dirIsNotHidden(filename) {
			if errWatcher := watcher.Add(filename); errWatcher != nil {
				return errWatcher
			}
		}
		return nil
	})

	if err != nil {
		log.Panicf("ошибка при добавлении директории для наблюдения: %v", err)
	}
	<-done
}

package main

import (
	detector "github.com/flameous/overlap-detection"
	"flag"
	"log"
	"os"
	"net/http"
	"os/signal"
	"syscall"
	"runtime/debug"
)

var (
	sig                chan os.Signal
	logFile            *os.File
	pathPtr            *string
	selfOverlappingPtr *bool
	done               chan bool
	files              chan string
	server             *http.Server
	db                 *detector.DB
)

func main() {
	defer func() {
		r := recover()
		if r != nil {
			log.Printf("непредвиденная ошибка: %#v. stacktrace: %s", r, debug.Stack())
		} else {
			log.Println("программа успешно остановлена")
		}
		logFile.Close()
	}()

	defer func() {
		server.Shutdown(nil)
		db.Close()
		close(sig)
		close(done)
		close(files)
	}()

	sig = make(chan os.Signal)
	done = make(chan bool)
	files = make(chan string, 1000)

	// мониторинг папок в отдельном потоке
	go detector.MonitorFolder(*pathPtr, done, files)
	go server.ListenAndServe()

	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case s := <-sig:
			done <- true
			log.Println("получен сигнал: " + s.String())
			return

		case filename := <-files:
			overlaps, err := detector.HandleCSV(filename, *selfOverlappingPtr)
			if err != nil {
				log.Printf("файл '%s' не обработан: %v", filename, err)
				continue
			}
			if err = db.SaveFileData(filename, overlaps); err != nil {
				log.Printf("файл '%s', ошибка при сохрании данных: %v", filename, err)
				continue
			}
			log.Printf("файл '%s' успешно обработан", filename)
		}
	}
}

func init() {
	selfOverlappingPtr = flag.Bool("self", false, "учитывать ли наложения на свой же путь?")
	pathPtr = flag.String("path", "/host_dir", "путь до директории для наблюдения")
	flag.Parse()

	var err error
	logFile, err = os.OpenFile("./logs.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Panicf("ошибка при создании лог-файла: %v", err)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(logFile)

	db = detector.NewDB("./db.file")

	http.HandleFunc("/get_files_list", GetProcessedFileList)
	http.HandleFunc("/get_overlaps", GetOverlapsList)
	server = &http.Server{Addr: ":8080"}
	log.Println("программа запущена")
}

// получение списка обработанных файлов
func GetProcessedFileList(w http.ResponseWriter, r *http.Request) {
	data, err := db.GetProcessedFileList()
	var response = data
	var respCode = http.StatusOK
	if err != nil {
		log.Printf("server: '/get_files_list', error: %v", err)
		respCode = http.StatusInternalServerError
		response = []byte(err.Error())
	}
	w.WriteHeader(respCode)
	w.Write(response)
}

// получение массива наложений у обработанного файла
func GetOverlapsList(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var respCode = http.StatusOK
	filename := r.FormValue("file")
	b, err := db.GetFileData(filename)
	if err != nil {
		log.Printf("server: '/get_overlaps', filename: %s, error: %v", filename, err)
		respCode = http.StatusInternalServerError
		b = []byte(err.Error())
	}
	if b == nil {
		respCode = http.StatusNotFound
		b = []byte("file: '" + filename + "' not found")
	}
	w.WriteHeader(respCode)
	w.Write(b)
}

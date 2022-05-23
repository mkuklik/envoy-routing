package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"go.opencensus.io/zpages"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func timeHandler(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Write([]byte("The time is: " + tm))
}

type BackendServer struct {
	mux     *http.ServeMux
	addr    string
	name    string
	port    int
	healthy int64
}

func NewBackendServer(port int, name string) *BackendServer {
	log.Info(fmt.Sprintf("New backend `%s` on :%d", name, port))
	s := BackendServer{
		addr:    fmt.Sprintf(":%d", port),
		mux:     http.NewServeMux(),
		name:    name,
		port:    port,
		healthy: 1,
	}
	return &s
}

func (s *BackendServer) Init() {

	s.mux.HandleFunc("/", s.main_handler)
	s.mux.HandleFunc("/headers", s.headers)
	s.mux.HandleFunc("/time", timeHandler)
	s.mux.Handle("/healthz", s.healthz())

	zpages.Handle(s.mux, "/debug")

}

func (s *BackendServer) main_handler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	fmt.Fprintf(w, "this is server %s\n", s.name)
}

func (s *BackendServer) headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func (s *BackendServer) healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt64(&s.healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

func (s *BackendServer) Start() {
	log.Infof("Starting server `%s` on %s", s.name, s.addr)
	if err := http.ListenAndServe(s.addr, NewRequestLogger(s.mux)); err != nil {
		panic(err.Error())
	}
}

type RequestLogger struct {
	handler http.Handler
}

func (rl RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Printf("Started %s %s", r.Method, r.URL.Path)
	defer func() {
		// requestID, ok := r.Context().Value(requestIDKey).(string)
		// if !ok {
		// 	requestID = "unknown"
		// }
		// requestID := "unknown"
		// s.logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	}()
	rl.handler.ServeHTTP(w, r)
}

func NewRequestLogger(handlerToWrap http.Handler) *RequestLogger {
	return &RequestLogger{handlerToWrap}
}

type Options struct {
	Port int
	Name string `validate:"required"`
}

func main() {
	flag.Int("port", 8000, "port")
	flag.String("name", "", "server name")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.BindEnv("port", "POD_PORT")
	viper.BindEnv("name", "POD_NAME")

	opts := Options{}

	err := viper.Unmarshal(&opts)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	log.Infof("after parse %s", opts.Name)
	backend := NewBackendServer(opts.Port, opts.Name)
	backend.Init()
	backend.Start()

	// router := http.NewServeMux()
	// router.HandleFunc("/headers", func(w http.ResponseWriter, req *http.Request) {
	// 	for name, headers := range req.Header {
	// 		for _, h := range headers {
	// 			fmt.Fprintf(w, "%v: %v\n", name, h)
	// 		}
	// 	}
	// })
	// s := http.Server{Addr: ":8000", Handler: router}
	// if err := s.ListenAndServe(); err != nil {
	// 	panic(err.Error())
	// }
}

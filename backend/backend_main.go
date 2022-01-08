package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func timeHandler(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Write([]byte("The time is: " + tm))
}

type BackendServer struct {
	http.Server
	name    string
	port    int
	healthy int64
	logger  log.Logger
	// mux  *http.ServeMux
}

func NewBackendServer(port int, name string) *BackendServer {
	glog.Info(fmt.Sprintf(":%d", port))
	s := BackendServer{
		Server: http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
		name: name,
		port: port,
	}
	return &s
}

func (s *BackendServer) Init() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.main_handler)
	mux.HandleFunc("/headers", s.headers)
	mux.HandleFunc("/time", timeHandler)
	mux.Handle("/healthz", s.healthz())
	s.Handler = s.logging(mux)
}

func (s *BackendServer) main_handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "this is server %s\n", s.name)
}

func (s *BackendServer) headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func (s *BackendServer) Start() {
	glog.Infof("Starting server `%s` on %s", s.name, s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err.Error())
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

func (s *BackendServer) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("%s %s %s %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		defer func() {
			// requestID, ok := r.Context().Value(requestIDKey).(string)
			// if !ok {
			// 	requestID = "unknown"
			// }
			// requestID := "unknown"
			// s.logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		}()
		next.ServeHTTP(w, r)
	})
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
		glog.Fatalf("unable to decode into struct, %v", err)
	}

	glog.Infof("after parse %s", opts.Name)
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

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	listenPort int
	endpoints  []string

	rootCmd = &cobra.Command{
		Use:   "api-router",
		Short: "Run the API router",
		Run: func(cmd *cobra.Command, args []string) {
			for _, cfg := range endpoints {
				endpoint := strings.Split(cfg, "->")
				if len(endpoint) != 2 {
					fmt.Fprintf(os.Stderr, "Invalid endpoint configuration \"%s\". Ignoring\n", cfg)
					continue
				}

				functionSegments := strings.SplitN(endpoint[1], ":", 2)
				if len(functionSegments) != 2 {
					fmt.Fprintf(os.Stderr, "Invalid function destination \"%s\". Ignoring\n", endpoint[1])
					continue
				}

				fmt.Printf("Registering function handler %s calling %s %s\n", endpoint[0], functionSegments[0], functionSegments[1])

				http.HandleFunc(endpoint[0], func(w http.ResponseWriter, r *http.Request) {
					req, err := http.NewRequest(functionSegments[0], functionSegments[1], nil)
					if err != nil {
						w.Write([]byte(err.Error()))
						return
					}

					client := &http.Client{
						Timeout: time.Minute * 5,
					}
					resp, err := client.Do(req)
					if err != nil {
						w.Write([]byte(err.Error()))
						return
					}

					defer resp.Body.Close()

					_, err = io.Copy(w, resp.Body)
					if err != nil {
						w.Write([]byte(err.Error()))
						return
					}
				})
			}

			if err := http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil); err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	flags := rootCmd.Flags()
	flags.StringSliceVar(&endpoints, "endpoint", []string{}, "List of API endpoints to route (/get:service2/foo)")
	flags.IntVar(&listenPort, "port", 80, "Port to listen on")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

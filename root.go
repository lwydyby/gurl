package gurl

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var curl Curl

var RootCmd = &cobra.Command{
	Use:   "gurl",
	Short: "gurl [options...] <url>",
	Long:  `gurl [options...] <url>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			for _, a := range args {
				if a != "\\" {
					curl.URL = a
					break
				}
			}
		}
		req := curl.Request()
		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		var content []byte
		content, err = io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		if curl.Include {
			fmt.Printf("%s %d\n", response.Proto, response.StatusCode)
			for header := range response.Header {
				for _, value := range response.Header.Values(header) {
					fmt.Printf("%s: %s\n", header, value)
				}
			}
			fmt.Print("\n")
		}

		fmt.Println(string(content))
	},
}

func init() {
	RootCmd.Flags().StringVarP(&curl.Method, "request", "X", "", "Specify request method to use")
	RootCmd.Flags().BoolVarP(&curl.Get, "get", "G", false, "Put the post data in the URL and use GET")
	RootCmd.Flags().StringArrayVarP(&curl.Header, "header", "H", []string{}, "Pass custom header(s) to server")
	RootCmd.Flags().StringVarP(&curl.Data, "data", "d", "", "HTTP POST data")
	RootCmd.Flags().StringVar(&curl.Data, "data-raw", "", "HTTP POST data, '@' allowed")
	RootCmd.Flags().StringArrayVarP(&curl.Form, "form", "f", []string{}, "Specify multipart MIME data")
	RootCmd.Flags().StringVar(&curl.URL, "url", "", "URL to work with")
	RootCmd.Flags().BoolVarP(&curl.Location, "location", "L", false, "Follow redirects")
	RootCmd.Flags().StringArrayVar(&curl.Form, "data-urlencode", []string{}, "HTTP POST data url encoded")
	RootCmd.Flags().BoolVar(&curl.Compressed, "compressed", false, "Request compressed response")
	RootCmd.Flags().BoolVarP(&curl.Include, "include", "i", false, "Include the HTTP response headers in the output. The HTTP response headers can include things like server name, cookies, date of the document, HTTP version and more.")
	RootCmd.Flags().BoolVarP(&curl.Insecure, "insecure", "k", false, "Allow insecure server connections when using SSL")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

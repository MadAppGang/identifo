package spa

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const DefaultIndexFilename = "index.html"

type SPASettings struct {
	Name       string
	Root       string
	FileSystem http.FileSystem
	Index      string
}

type SPAHandlerError struct {
	s string
	e error
}

// NewSpaHandlerFunc returns new handle to serve SPA apps
// highly inspired by echo spa static server
// https://github.com/labstack/echo/blob/master/middleware/static.go#L164
func NewSPAHandlerFunc(c SPASettings) http.HandlerFunc {
	// set default index file name if it not set
	if len(c.Index) == 0 {
		c.Index = DefaultIndexFilename
	}

	return func(w http.ResponseWriter, r *http.Request) {
		p, err := url.PathUnescape(r.URL.Path)
		if err != nil {
			return
		}
		name := filepath.Join(c.Root, filepath.Clean("/"+p)) // "/"+ for security. TODO: Jack: it is not clear why adding leading slash providing extra security

		// let's open the file
		file, err := openFile(c.FileSystem, name)
		if err != nil {
			// any error, but not found
			if !os.IsNotExist(err) {
				handleError(w, SPAHandlerError{s: "unable to read file", e: err})
				return
			}

			file, err = openFile(c.FileSystem, filepath.Join(c.Root, c.Index))
			if err != nil {
				handleError(w, SPAHandlerError{s: "unable to serve index file", e: err})
				return
			}
		}

		defer file.Close()

		// get information about the file or dir
		info, err := file.Stat()
		if err != nil {
			handleError(w, SPAHandlerError{s: "unable to fet file/dir stat", e: err})
			return
		}

		if info.IsDir() {
			// if the path is a dir, then we trying to get index file from dir
			index, err := openFile(c.FileSystem, filepath.Join(name, c.Index))
			if err != nil {
				handleError(w, SPAHandlerError{s: "unable to serve index file", e: err})
				return
			}

			defer index.Close()

			info, err = index.Stat()
			if err != nil {
				handleError(w, SPAHandlerError{s: "unable to fet file/dir stat", e: err})
			}

			buffer, err := ioutil.ReadAll(file)
			if err != nil {
				handleError(w, SPAHandlerError{s: "unable to fet file/dir stat", e: err})
			}

			// return the file
			http.ServeContent(w, r, info.Name(), info.ModTime(), bytes.NewReader(buffer))
			return
		}

		// return the file
		buffer, err := ioutil.ReadAll(file)
		if err != nil {
			handleError(w, SPAHandlerError{s: "unable to fet file/dir stat", e: err})
		}
		http.ServeContent(w, r, info.Name(), info.ModTime(), bytes.NewReader(buffer))
	}
}

func openFile(fs http.FileSystem, name string) (http.File, error) {
	pathWithSlashes := filepath.ToSlash(name)
	return fs.Open(pathWithSlashes)
}

// func serveFile(c echo.Context, file http.File, info os.FileInfo) error {
// 	http.ServeContent(c.Response(), c.Request(), info.Name(), info.ModTime(), file)
// 	return nil
// }

func handleError(w http.ResponseWriter, e SPAHandlerError) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, fmt.Sprintf(httpSPAError, e.Error()))
}

func (e SPAHandlerError) Error() string {
	return fmt.Sprintf("spa handler error with message: %s\n with underlying error: %v /n", e.s, e.e)
}

// if we failed to deliver error team - it is html error fallback screen to show the message
// something like BSoD for SPA
const httpSPAError = `
<!DOCTYPE html>
<html>
<head>
<title>Login SPA app error</title>
</head>
<body>

<h1>Error</h1>
<p>%s</p>

</body>
</html>
`

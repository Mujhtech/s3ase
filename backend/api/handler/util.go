package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const (
	QueryParamPage    = "page"
	QueryParamPerPage = "per_page"

	PerPageDefault = 10
	PerPageMax     = 100
)

type CookieOptions struct {
	Name     string
	Path     string
	State    string
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

func getCookie(r *http.Request, cookieName string) (string, bool) {
	cookie, err := r.Cookie(cookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return "", false
	} else if err != nil {
		// this should never happen - documentation and code only return `nil` or `http.ErrNoCookie`
		panic(fmt.Sprintf("unexpected error from request.Cookie(...) method: %s", err))
	}

	return cookie.Value, true
}

// getHeaderOrDefault returns the first non-empty value of the header with the given name.
// func getHeaderOrDefault(r *http.Request, headerName string, dflt string) string {
// 	val, ok := getHeader(r, headerName)
// 	if !ok {
// 		return dflt
// 	}

// 	return val
// }

// getHeader returns the first non-empty value of the header with the given name.
// func getHeader(r *http.Request, headerName string) (string, bool) {
// 	for _, val := range r.Header.Values(headerName) {
// 		if val != "" {
// 			return val, true
// 		}
// 	}

// 	return "", false
// }

// pathParamOrError returns the value of the path parameter with the given name or an error if the parameter is not found.
func pathParamOrError(r *http.Request, paramName string) (string, error) {
	val, ok := pathParam(r, paramName)
	if !ok {
		return "", fmt.Errorf("parameter '%s' not found in request path", paramName)
	}

	return val, nil
}

// pathParamOrDefault returns the value of the path parameter with the given name or the default value if the parameter is not found.
// func pathParamOrEmpty(r *http.Request, paramName string) string {
// 	val, ok := pathParam(r, paramName)
// 	if !ok {
// 		return ""
// 	}

// 	return val
// }

// pathParam returns the value of the path parameter with the given name.
func pathParam(r *http.Request, paramName string) (string, bool) {
	val := chi.URLParam(r, paramName)
	if val == "" {
		return "", false
	}

	return val, true
}

// queryParam returns the value of the query parameter with the given name.
func queryParam(r *http.Request, paramName string) (string, bool) {
	query := r.URL.Query()
	if !query.Has(paramName) {
		return "", false
	}

	return query.Get(paramName), true
}

// queryParamOrError returns the value of the query parameter with the given name or an error if the parameter is not found.
func queryParamOrError(r *http.Request, paramName string) (string, error) {
	val, ok := queryParam(r, paramName)
	if !ok {
		return "", fmt.Errorf("parameter '%s' not found in query", paramName)
	}

	return val, nil
}

// queryParamOrDefault returns the value of the query parameter with the given name or the default value if the parameter is not found.
func queryParamOrDefault(r *http.Request, paramName string, deflt string) string {
	val, ok := queryParam(r, paramName)
	if !ok {
		return deflt
	}

	return val
}

// queryParamList returns the list of values of the query parameter with the given name.
// func queryParamList(r *http.Request, paramName string) ([]string, bool) {
// 	query := r.URL.Query()
// 	if !query.Has(paramName) {
// 		return nil, false
// 	}

// 	return query[paramName], true
// }

// ParsePage returns the page number.
func ParsePage(r *http.Request) int {
	s := r.URL.Query().Get(QueryParamPage)
	i, _ := strconv.Atoi(s)
	if i <= 0 {
		i = 1
	}
	return i
}

// ParsePerPage returns the number of items per page.
func ParsePerPage(r *http.Request) int {
	s := r.URL.Query().Get(QueryParamPerPage)
	i, _ := strconv.Atoi(s)
	if i <= 0 {
		i = PerPageDefault
	} else if i > PerPageMax {
		i = PerPageMax
	}
	return i
}

// createCookie creates a cookie with the given options.
func createCookie(w http.ResponseWriter, opts CookieOptions) {
	cookie := &http.Cookie{
		Name:     opts.Name,
		Value:    opts.State,
		Path:     opts.Path,
		HttpOnly: opts.HttpOnly,
		Secure:   opts.Secure,
		MaxAge:   opts.MaxAge, // 5 minutes
		SameSite: opts.SameSite,
	}
	http.SetCookie(w, cookie)
}

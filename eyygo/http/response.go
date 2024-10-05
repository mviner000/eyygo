package http

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/project_name"
)

// Updated Response struct to include Template
type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       interface{}
	Template   string
}

// Updated HttpResponseOK function
func HttpResponseOK(body interface{}, headers map[string]string, template ...string) *Response {
	var tmpl string
	if len(template) > 0 {
		tmpl = template[0]
	}

	return &Response{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       body,
		Template:   tmpl,
	}
}

func (r *Response) Render(c *fiber.Ctx) error {
	for key, value := range r.Headers {
		c.Set(key, value)
	}
	if r.Template != "" {
		// Dynamically resolve the template path using AppSettings
		templatePath := project_name.AppSettings.TemplateBasePath + r.Template

		// Log the template path for debugging
		fmt.Println("Rendering template at path:", templatePath)

		return c.Render(templatePath, r.Body)
	} else if r.Body == nil {
		return c.SendStatus(r.StatusCode)
	} else {
		return c.Status(r.StatusCode).JSON(r.Body)
	}
}

func HttpResponse(body interface{}, status int, headers map[string]string) *Response {
	return &Response{
		StatusCode: status,
		Headers:    headers,
		Body:       body,
	}
}

func HttpResponseCreated(body interface{}, headers map[string]string) *Response {
	return HttpResponse(body, http.StatusCreated, headers)
}

func HttpResponseNoContent(headers map[string]string) *Response {
	return HttpResponse(nil, http.StatusNoContent, headers)
}

func HttpResponseBadRequest(body interface{}, headers map[string]string) *Response {
	return HttpResponse(body, http.StatusBadRequest, headers)
}

func HttpResponseUnauthorized(body interface{}, headers map[string]string) *Response {
	return HttpResponse(body, http.StatusUnauthorized, headers)
}

func HttpResponseForbidden(body interface{}, headers map[string]string) *Response {
	return HttpResponse(body, http.StatusForbidden, headers)
}

func HttpResponseNotFound(body interface{}, headers map[string]string) *Response {
	return HttpResponse(body, http.StatusNotFound, headers)
}

func HttpResponseMethodNotAllowed(body interface{}, headers map[string]string) *Response {
	return HttpResponse(body, http.StatusMethodNotAllowed, headers)
}

func HttpResponseServerError(body interface{}, headers map[string]string) *Response {
	return HttpResponse(body, http.StatusInternalServerError, headers)
}

func HttpResponseRedirect(url string, permanent bool) *Response {
	status := http.StatusFound
	if permanent {
		status = http.StatusMovedPermanently
	}
	return &Response{
		StatusCode: status,
		Headers:    map[string]string{"Location": url},
	}
}

func JsonResponse(data interface{}, status int, headers map[string]string) *Response {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	return HttpResponse(data, status, headers)
}

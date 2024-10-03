package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       interface{}
}

func (r *Response) Render(c *fiber.Ctx) error {
	for key, value := range r.Headers {
		c.Set(key, value)
	}
	if r.Body == nil {
		return c.SendStatus(r.StatusCode)
	}
	return c.Status(r.StatusCode).JSON(r.Body)
}

func HttpResponse(body interface{}, status int, headers map[string]string) *Response {
	return &Response{
		StatusCode: status,
		Headers:    headers,
		Body:       body,
	}
}

func HttpResponseOK(body interface{}, headers map[string]string) *Response {
	return HttpResponse(body, http.StatusOK, headers)
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

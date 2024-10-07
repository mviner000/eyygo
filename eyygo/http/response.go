package http

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/aymerick/raymond"
	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/project_name"
)

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       interface{}
	Template   string
	Layout     string
}

func HttpResponseOK(body interface{}, headers map[string]string, template string, layout ...string) *Response {
	var layoutTemplate string
	if len(layout) > 0 {
		layoutTemplate = layout[0]
	}

	return &Response{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       body,
		Template:   template,
		Layout:     layoutTemplate,
	}
}

func (r *Response) Render(c *fiber.Ctx) error {
	for key, value := range r.Headers {
		c.Set(key, value)
	}

	if r.StatusCode == http.StatusFound || r.StatusCode == http.StatusMovedPermanently {
		return c.Redirect(r.Headers["Location"], r.StatusCode)
	}

	if r.Template != "" {
		templatePath := filepath.Join(project_name.AppSettings.TemplateBasePath, r.Template)
		layoutPath := ""
		if r.Layout != "" {
			layoutPath = filepath.Join(project_name.AppSettings.TemplateBasePath, r.Layout)
		}

		content, err := renderTemplate(templatePath, layoutPath, r.Body)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Template rendering error: %v", err))
		}

		c.Set("Content-Type", "text/html")
		return c.Status(r.StatusCode).Send(content)
	} else if r.Body == nil {
		return c.SendStatus(r.StatusCode)
	} else {
		return c.Status(r.StatusCode).JSON(r.Body)
	}
}

func renderTemplate(templatePath, layoutPath string, data interface{}) ([]byte, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}

	var content bytes.Buffer
	if err := tmpl.Execute(&content, data); err != nil {
		return nil, fmt.Errorf("error executing template: %v", err)
	}

	if layoutPath != "" {
		layout, err := raymond.ParseFile(layoutPath)
		if err != nil {
			return nil, fmt.Errorf("error parsing layout: %v", err)
		}

		dataMap, ok := data.(fiber.Map)
		if !ok {
			dataMap = fiber.Map{}
		}

		layoutData := fiber.Map{
			"yield": raymond.SafeString(content.String()),
			"title": dataMap["MetaTitle"],
		}

		for k, v := range dataMap {
			if k != "MetaTitle" {
				layoutData[k] = v
			}
		}

		result, err := layout.Exec(layoutData)
		if err != nil {
			return nil, fmt.Errorf("error executing layout: %v", err)
		}

		return []byte(result), nil
	}

	return content.Bytes(), nil
}

func HttpResponseHTMX(body interface{}, template string, layout ...string) *Response {
	headers := map[string]string{"Content-Type": "text/html"}
	var layoutTemplate string
	if len(layout) > 0 {
		layoutTemplate = layout[0]
	}
	return HttpResponseOK(body, headers, template, layoutTemplate)
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
	status := http.StatusFound // 302
	if permanent {
		status = http.StatusMovedPermanently // 301
	}
	return &Response{
		StatusCode: status,
		Headers:    map[string]string{"Location": url},
	}
}

func WindowReload(url string) string {
	return `<script>window.location.href='` + url + `'</script>`
}

func JsonResponse(data interface{}, status int, headers map[string]string) *Response {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	return HttpResponse(data, status, headers)
}

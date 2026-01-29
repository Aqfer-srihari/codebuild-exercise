package handler

import (
	"encoding/json"
	"net/http"
	"question5updation/internal/response"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type MethodHandler func(*http.Request, map[string]string) response.APIResponse

var routes = map[string]map[string]MethodHandler{
	"/version": {
		http.MethodGet: VersionHandler,
	},
	"/employees": {
		http.MethodGet:  GetEmployees,
		http.MethodPost: CreateEmployee,
	},
	"/employees/{name}": {
		http.MethodGet:    GetEmployeeByName,
		http.MethodPut:    UpdateEmployee,
		http.MethodDelete: DeleteEmployee,
	},
}

func LambdaRouter(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	url := req.Path
	if len(req.QueryStringParameters) > 0 {
		q := make([]string, 0)
		for k, v := range req.QueryStringParameters {
			q = append(q, k+"="+v)
		}
		url = url + "?" + strings.Join(q, "&")
	}

	httpReq, _ := http.NewRequest(
		req.HTTPMethod,
		url,
		strings.NewReader(req.Body),
	)

	httpReq.Header.Set("Content-Type", "application/json")

	for route, methods := range routes {
		params, ok := matchRoute(route, req.Path)
		if !ok {
			continue
		}

		handler, exists := methods[req.HTTPMethod]
		if !exists {
			return convertResponse(response.APIResponse{
				Status: http.StatusMethodNotAllowed,
				Err:    response.ErrMethodNotAllowed,
			})
		}

		resp := handler(httpReq, params)
		return convertResponse(resp)
	}

	return convertResponse(response.APIResponse{
		Status: http.StatusNotFound,
		Err:    response.ErrRouteNotFound,
	})
}

func matchRoute(route, path string) (map[string]string, bool) {
	routeParts := strings.Split(strings.Trim(route, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(routeParts) != len(pathParts) {
		return nil, false
	}

	params := make(map[string]string)
	for i := range routeParts {
		if strings.HasPrefix(routeParts[i], "{") {
			key := strings.Trim(routeParts[i], "{}")
			params[key] = pathParts[i]
		} else if routeParts[i] != pathParts[i] {
			return nil, false
		}
	}
	return params, true
}

func convertResponse(resp response.APIResponse) events.APIGatewayProxyResponse {

	if resp.Err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: resp.Status,
			Body:       resp.Err.Error(),
		}
	}

	bodyBytes, err := json.Marshal(resp.Data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to marshal response data",
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: resp.Status,
		Body:       string(bodyBytes),
	}
}

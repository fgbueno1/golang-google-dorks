package golangGoogleDorks

import (
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const (
	SEARCH_URL = "https://customsearch.googleapis.com/customsearch/v1?key={api_key}&cx={cx}&q={query}&start={start}"
)

// Dork receive a query to search on google based on api-key and searchiengine provided.
func Dork(query string, apiKey string, cx string) (result []*GoogleSearchResults, err error) {
	cli := resty.New()
	params := make(map[string]string)
	params["api_key"] = apiKey
	params["cx"] = cx
	params["start"] = "1"
	params["query"] = query
	for {
		resp, err := cli.R().SetPathParams(params).SetResult(&GoogleSearchResults{}).Get(SEARCH_URL)
		if err != nil {
			return result, err
		}
		err = CheckError(resp)
		if err != nil {
			return result, err
		}
		output := resp.Result().(*GoogleSearchResults)
		result = append(result, output)
		if output.Queries.NextPage == nil {
			break
		}
		for _, nextPage := range output.Queries.NextPage {
			params["start"] = strconv.FormatInt(int64(nextPage.StartIndex), 10)
		}
	}
	return
}

// CheckError validate if an response is an error or not.
func CheckError(response *resty.Response) (resolvedError error) {
	if response.IsError() {
		apiError := response.Error().(*APIError)
		resolvedError = fmt.Errorf(
			"the application returned an error of type %v with message: %s",
			apiError.Error.Code, apiError.Error.Message,
		)
	}
	return resolvedError
}

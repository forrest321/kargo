package pantry

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Export(pantryID, basketName, contents string) (string, error) {
	url := fmt.Sprintf("https://getpantry.cloud/apiv1/pantry/%s/basket/%s", pantryID, basketName)
	method := "PUT"

	payload := strings.NewReader(contents)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(body), nil
}

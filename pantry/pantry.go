package pantry

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Export(pantryID, basketName, contents string) (string, error) {
	url := fmt.Sprintf("https://getpantry.cloud/apiv1/pantry/%s/basket/%s", pantryID, basketName)

	payload := strings.NewReader(fmt.Sprintf(`%s`, contents))

	res, err := http.Post(url, "application/json", payload)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(body), nil
}

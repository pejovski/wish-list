package catalog

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pejovski/wish-list/domain"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Gateway struct {
	client *retryablehttp.Client
	host   string
}

func NewGateway(c *retryablehttp.Client, host string) Gateway {
	return Gateway{client: c, host: host}
}

func (g Gateway) Product(id string) (*domain.Product, error) {

	url := g.host + fmt.Sprintf("/products/%s", id)

	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorln("Failed to create new request", err)
		return nil, err
	}

	res, err := g.client.Do(req)
	if err != nil {
		logrus.Errorln("Failed to Do request", err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		logrus.Errorln(ErrorNotOk, strconv.Itoa(res.StatusCode))
		return nil, ErrorNotOk
	}

	var p *Product
	err = json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
		logrus.Errorln("Failed to Decode", err)
		return nil, err
	}
	defer res.Body.Close()

	return g.mapProductToDomainProduct(p), nil
}

package fns_api

import (
	"ShoppingExpensesService/internal/config"
	"ShoppingExpensesService/internal/models/dbmodels"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func GetQRInfo(tgID int, qrCodeInfo string, cfg *config.Config) ([]dbmodels.Product, error) {
	body, err := json.Marshal(map[string]string{
		"token": cfg.App.API.Token,
		"qrraw": qrCodeInfo,
	})
	if err != nil {
		return nil, err
	}
	res, err := http.Post(cfg.App.API.URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	if body, err = io.ReadAll(res.Body); err != nil {
		return nil, err
	}
	var resBody FNSResponse
	err = json.Unmarshal(body, &resBody)
	if err != nil {
		return nil, err
	}

	return ProductsList(tgID, resBody)
}

func ProductsList(tgID int, res FNSResponse) ([]dbmodels.Product, error) {
	var products []dbmodels.Product
	for _, item := range res.Data.Json.Items {
		prod := dbmodels.Product{
			ID:       0,
			UserID:   tgID,
			Name:     item.Name,
			Market:   res.Data.Json.User,
			Price:    item.Price,
			Category: strconv.Itoa(item.ProductType),
			Count:    item.Quantity,
			Barcode:  item.ProductCodeNew.Ean13.Gtin,
		}
		products = append(products, prod)
	}
	return products, nil
}

// Package classification 商城API
//
// 商城API
//
//      Host: localhost:4006/api/v1
//      Version: 0.0.1
//
// SecurityDefinitions:
//  authorizationHeaderToken:
// 	 type: apiKey
//   name: Authorization
//   in: header
// swagger:meta
package main

// swagger:parameters GetGoodsListParam
type GetGoodsListParam struct {

}

// swagger:parameters GetGoodsStockParam
type GetGoodsStockParam struct {
}

// swagger:parameters BuyGoodsParam
type BuyGoodsParam struct {
	// in: body
	Body struct {
			Coin      string  `json:"coin"`
			GoodsId   float64 `json:"goods_id"`
			Price     float64 `json:"price"`
			GoodsName string  `json:"goods_name"`
			Memo      string  `json:"memo"`
			Amount    float64 `json:"amount"`
	}
}


// swagger:parameters InsertGoodsParam
type InsertGoodsParam struct {
	// in: body
	Body struct {
		// in: query
		GoodsName	string 		`json:"goods_name"`
		Price	float64 	`json:"price"`
		ImageUrl	string 		`json:"image_url"`
	}
}


// swagger:parameters UpdateGoodsOrderParam
type UpdateGoodsOrderParam struct {
	// in: body
	Body struct {
		// in: query
		Amount int64 `json:"amount"`
		GoodsName string `json:"goods_name"`
		GoodsId int64`json:"goods_id"`
		Memo string `json:"memo"`
		Status int64 `json:"status"`
		GoodsOrderId int64 `json:"goods_order_id"`
	}
}

// swagger:parameters GetGoodsOrder
type GetGoodsOrder struct {
}



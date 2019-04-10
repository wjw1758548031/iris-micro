package routers

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"rubik/server/common/middleware"
	"rubik/server/goods/controllers"
)


func SetRouter(router iris.Party, path string) iris.Party {
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})
	router = router.Party("/", crs).AllowMethods(iris.MethodOptions)
	r := router.Party(path)
	r.Get("/get_goods_list",controllers.GetGoodsList)//查询商品列表
	r.Post("/buy_goods",middleware.AuthToken,controllers.BuyGoods)//购买商品
	r.Post("/insert_goods",middleware.AuthToken,controllers.InsertGoods)//发布商品
	//修改商品订单流程
	//三个流程：1：待发货、2：发货中、3：已收货
	//根据用户权限校验
	r.Post("/update_goods_order",middleware.AuthToken,controllers.UpdateGoodsOrder)
	r.Post("/get_goods_order",controllers.GetGoodsOrder)//查询商品订单
	r.Post("/get_goods_stock",middleware.AuthToken,controllers.GetGoodsStock)//用户查询商品库存

	return router
}
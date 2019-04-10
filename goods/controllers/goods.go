package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/mongodb/mongo-go-driver/bson"
	"rubik/server/common"
	"rubik/server/common/consts"
	"rubik/server/common/rpc"
	"rubik/server/goods/mongo"
	"rubik/server/goods/services"
	"time"
)


func init(){
	//初始化用户权限表
	count, _ := mongo.Client.Count("goods_role",bson.M{})
	//没有用户权限，则进行默认admin
	if count == 0{
		goods_role := make(map[string]interface{})
		goods_role["uid"] = int64(1)
		goods_role["code"] = "001"
		goods_role["goodsRoleId"] = "1"
		goods_role["goodsRoleName"] = "超级管理员"
		goods_role["insertTime"] = time.Now()
		goods_role["updateTime"]= time.Now()
		err := mongo.Client.Insert("goods_role",bson.M(goods_role))
		fmt.Println("insert_goods_role：",err.Error())
	}
}


//查询商品列表
func GetGoodsList(ctx iris.Context) {
	// swagger:route GET /get_goods_list 查询商品列表 GetGoodsListParam
	//
	// 查询商品列表
	//
	// 查询商品列表
	//
	// responses:
	//	200: ResponseSuccess
	//	400: ResponseError
	var result []map[string]interface{}
	//mongodb商品列表
	err := mongo.Client.GoodsFind("goods",bson.M{},&result);if err!=nil{
		ctx.JSON(iris.Map{"err":  err.Error(),})
		return
	}

	//商品个数
	count, err := mongo.Client.Count("goods",bson.M{});if err!=nil{
		ctx.JSON(iris.Map{"err":  err.Error(),})
		return
	}

	ctx.JSON(iris.Map{
		"count":  count,
		"result":   result,
	})

}


//发布商品
func InsertGoods(ctx iris.Context) {
	// swagger:route POST /insert_goods 发布商品 InsertGoodsParam
	//
	// 发布商品
	//
	// 发布商品
	//
	// Security:
	// 	authorizationHeaderToken:
	// responses:
	//	200: ResponseSuccess
	//	400: ResponseError
	post := make(map[string]interface{})
	_ = ctx.ReadJSON(&post)
	//校验参数
	if services.Check(post["goods_name"],post["price"],post["image_url"]) == false{
		ctx.JSON(iris.Map{"err":  "参数校验失败",})
		return
	}
	uid := commons.ToInt64(ctx.Values().Get("users").(map[string]interface{})["uid"])
	goodsName := post["goods_name"].(string)
	price := post["price"].(float64)
	imageUrl := post["image_url"].(string)

	//查看uid权限
	role_len,err := services.UserGoodsRole(uid);if err != nil || role_len == 0{
		ctx.JSON(iris.Map{"err":  err.Error(),"data" : "该用户不是超级管理员"})
		return
	}

	//获取goodsId
	goodsId,err := services.FindMaxId("goods","goodsId");if err != nil{
		ctx.JSON(iris.Map{"err":  err.Error(),})
		return
	}

	//发布goods商品
	err = services.InsertGoods(goodsName,price,imageUrl,goodsId,time.Now(),time.Now());if err != nil{
		ctx.JSON(iris.Map{"err":  err.Error(),})
		return
	}
	ctx.JSON(iris.Map{"data":  "发布商品成功",})
}


func BuyGoods(ctx iris.Context){
	// swagger:route POST /buy_goods 购买商品 BuyGoodsParam
	//
	// 购买商品
	//
	// 购买商品
	//
	// Security:
	// 	authorizationHeaderToken:
	// responses:
	//	200: ResponseSuccess
	//	400: ResponseError
	post := make(map[string]interface{})
	_ = ctx.ReadJSON(&post)
	//校验参数
	if services.Check(post["coin"],post["goods_id"],post["price"],post["goods_name"],post["memo"],post["amount"]) == false{
		ctx.JSON(iris.Map{"err":  "参数校验失败",})
		return
	}
	uid := commons.ToInt64(ctx.Values().Get("users").(map[string]interface{})["uid"])
	coin := post["coin"].(string)
	goodsId := post["goods_id"].(float64)
	price := post["price"].(float64)
	goodsName := post["goods_name"].(string)
	memo := post["memo"].(string)
	amount := post["amount"].(float64)
	outUid := int64(0)
	putUid := uid

	//获取goodsId
	goodsOrderId,err := services.FindMaxId("goods_order","goodsOrderId");if err != nil{
		ctx.JSON(iris.Map{"err":  err.Error(),})
		return
	}

	//扣款
	err = rpc.BaseServer.ChangeMoney(uid, coin, -price * amount, consts.MONEY_METHOD_BANK, string(goodsOrderId))
	if err != nil {
		ctx.JSON(iris.Map{"err":  err.Error(),})
		return
	}

	//购买商品,新增订单流程
	err = services.InsertGoodsOrder(int64(goodsId) ,price ,goodsName ,memo ,int64(amount) ,outUid ,putUid ,goodsOrderId );if err!=nil{
		ctx.JSON(iris.Map{"err":  err.Error(),})
		return
	}

	ctx.JSON(iris.Map{"data":  "购买商品成功",})

}


func UpdateGoodsOrder(ctx iris.Context) {
	// swagger:route POST /update_goods_order 修改商品订单流程 UpdateGoodsOrderParam
	//
	// 修改商品订单流程
	//
	// 修改商品订单流程
	//
	// Security:
	// 	authorizationHeaderToken:
	// responses:
	//	200: ResponseSuccess
	//	400: ResponseError
	post := make(map[string]interface{})
	_ = ctx.ReadJSON(&post)
	//校验参数
	if services.Check(post["amount"],post["goods_name"],post["goods_id"],post["memo"],post["status"],post["goods_order_id"]) == false{
		ctx.JSON(iris.Map{"err":  "参数校验失败",})
		return
	}
	uid := commons.ToInt64(ctx.Values().Get("users").(map[string]interface{})["uid"])
	memo := post["memo"].(string)
	status := int64(post["status"].(float64))
	outUid := uid
	goodsOrderId := int64(post["goods_order_id"].(float64))
	goodsId := int64(post["goods_id"].(float64))
	goodsName :=  post["goods_name"].(string)
	amount :=int64(post["amount"].(float64))


	//判断状态
	if  status!= 2 && status != 3{
		ctx.JSON(iris.Map{"err":  "您所调整的状态不在数据中",})
		return
	}

	//状态一为待发货，二为已发货，三为已收货
	if status == 2{
		//查看uid权限
		role_len,err := services.UserGoodsRole(uid);if err != nil || role_len == 0{
			ctx.JSON(iris.Map{"err":  err.Error(),"data" : "该用户不是超级管理员,无法发货"})
			return
		}
		//修改订单状态
		err = services.UpdateGoodsOrder(goodsOrderId ,status ,memo ,outUid );if err!=nil{
			ctx.JSON(iris.Map{"err":  err.Error()})
			return
		}
		ctx.JSON(iris.Map{"data":  "发货成功"})

	}else{
		//修改订单状态
		err := services.UpdateGoodsOrder(goodsOrderId ,status ,memo ,-1 );if err!=nil{
			ctx.JSON(iris.Map{"err":  err.Error()})
			return
		}

		//查看用户库存
		var result []map[string]interface{}
		err = mongo.Client.GoodsFind("user_goods_stock",bson.M{"uid":uid,"status":1,"goodsId":goodsId},&result);if err!=nil{
			ctx.JSON(iris.Map{"err":  err.Error(),})
			return
		}
		//没有库存新增否则改变数量
		if len(result) == 0{
			services.InsertGoodsStock(uid ,goodsId ,goodsName ,1 ,amount)
		}else{
			services.UpdateGoodsStock(goodsId,uid,amount+result[0]["amount"].(int64))
		}
		ctx.JSON(iris.Map{"data":  "收货成功"})
	}

}



//查询商品订单列表
func GetGoodsStock(ctx iris.Context) {
	// swagger:route POST /get_goods_stock 用户查询商品库存 GetGoodsStockParam
	//
	// 用户查询商品库存
	//
	// 用户查询商品库存
	//
	// Security:
	// 	authorizationHeaderToken:
	// responses:
	//	200: ResponseSuccess
	//	400: ResponseError
	uid := commons.ToInt64(ctx.Values().Get("users").(map[string]interface{})["uid"])
	var result []map[string]interface{}
	//mongodb商品列表
	err := mongo.Client.GoodsFind("user_goods_stock",bson.M{"uid":uid},&result);if err!=nil{
		ctx.JSON(iris.Map{"err":  err.Error(),})
		return
	}

	//商品个数
	count, err := mongo.Client.Count("user_goods_stock",bson.M{"uid":uid});if err!=nil{
		ctx.JSON(iris.Map{"err":  err.Error(),})
		return
	}

	ctx.JSON(iris.Map{
		"count":  count,
		"result":   result,
	})

}

//查询商品订单
func GetGoodsOrder(ctx iris.Context) {
	// swagger:route POST /get_goods_order 查询商品订单 GetGoodsOrder
	//
	// 查询商品订单
	//
	// 查询商品订单
	//
	// responses:
	//	200: ResponseSuccess
	//	400: ResponseError
	var result []map[string]interface{}
	//mongodb商品列表
	err := mongo.Client.GoodsFind("goods_order",bson.M{},&result);if err!=nil{
		ctx.JSON(iris.Map{"err":  err.Error(),})
		return
	}

	//商品个数
	count, err := mongo.Client.Count("goods_order",bson.M{});if err!=nil{
		ctx.JSON(iris.Map{"err":  err.Error(),})
		return
	}

	ctx.JSON(iris.Map{
		"count":  count,
		"result":   result,
	})

}

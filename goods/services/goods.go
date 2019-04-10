package services

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"rubik/server/goods/mongo"
	"time"
)

//传入的参数不能为空
func Check(orther ...interface{}) bool{
	for  _,v := range orther{
		if v == nil || v == ""{
			return false
		}
	}
	return true
}


func UpdateGoodsOrder(goodsOrderId int64,status int64,memo string,outUid int64) error {
	goodsOrder := make(map[string]interface{})
	if outUid != -1{
		goodsOrder["outUid"] = outUid
	}
	goodsOrder["status"] = status
	goodsOrder["memo"] = memo
	goodsOrder["updateTime"] = time.Now()
	check := make(map[string]interface{})
	check["goodsOrderId"] = goodsOrderId
	err := mongo.Client.Update("goods_order",bson.M(check),bson.M{"$set":bson.M(goodsOrder)})
	return err
}

func InsertGoodsStock(uid int64,goodsId int64,goodsName string,status int64,amount int64) error {
	goodsOrder := make(map[string]interface{})
	goodsOrder["uid"] = uid
	goodsOrder["goodsId"] = goodsId
	goodsOrder["goodsName"] = goodsName
	goodsOrder["status"] = status
	goodsOrder["amount"] = amount
	goodsOrder["insertTime"] = time.Now()
	goodsOrder["updateTime"] = time.Now()
	err := mongo.Client.Insert("user_goods_stock",bson.M(goodsOrder))
	return err
}

func UpdateGoodsStock(goodsId,uid,amount int64) error {
	err := mongo.Client.Update("user_goods_stock",bson.M{"uid":uid,"goodsId":goodsId},bson.M{"$set": bson.M{"amount":amount,"updateTime":time.Now()}})
	return err
}

func InsertGoodsOrder(goodsId int64,price float64,goodsName string,memo string,amount int64,outUid int64,putUid int64,goodsOrderId int64) error {
	goodsOrder := make(map[string]interface{})
	goodsOrder["goodsId"] = goodsId
	goodsOrder["price"] = price
	goodsOrder["goodsName"] = goodsName
	goodsOrder["memo"] = memo
	goodsOrder["amount"] = amount
	goodsOrder["outUid"] = outUid
	goodsOrder["putUid"] = putUid
	goodsOrder["status"] = int64(1)
	goodsOrder["goodsOrderId"] = goodsOrderId
	goodsOrder["insertTime"] = time.Now()
	goodsOrder["updateTime"] = time.Now()
	goodsOrder["endTime"] = int64(time.Now().Unix()+ 864000)
	err := mongo.Client.Insert("goods_order",bson.M(goodsOrder))
	return err
}

func InsertGoods(goodsName string , price float64 , imageUrl string , goodsId int64 ,insertTime time.Time,updateTime time.Time) error {
	goods := make(map[string]interface{})
	goods["goodsName"] = goodsName
	goods["price"] = price
	goods["imageUrl"] = imageUrl
	goods["goodsId"] = goodsId
	goods["insertTime"] = time.Now()
	goods["updateTime"] = time.Now()
	err := mongo.Client.Insert("goods",bson.M(goods))
	return err
}


func UserGoodsRole(uid int64) (int,error) {
	//获取用户权限
	var result []map[string]interface{}
	err :=mongo.Client.Aggregate("goods_role",[]bson.M{{"$match":bson.M{"code":"001","uid":uid}}},&result)
	return len(result),err
}

func FindMaxId(table,id string) (int64,error) {
	//获取最大的id进行增加
	var result []map[string]interface{}
	err :=mongo.Client.Aggregate(table,[]bson.M{{"$sort":bson.M{id:-1}},{"$limit":1}},&result)
	rid := int64(1)
	if result != nil{
		rid = result[0][id].(int64)+1
	}
	return rid,err
}



/*
//获取当天支付数量
func GetTodayPayAmount(uid int64, coin string) float64 {
	zero := commons.GetZeroTime()
	var history []interface{}
	mongo.Client.Find("moneyhistories", bson.M{"uid":uid, "type":-1, "create":bson.M{"$gte":zero}, "method":11, "coin":coin}, &history)
	total := 0.0
	for _, v := range history{
		total += v.(map[string]float64)["amount"]
	}
	return total
}
//支付
func Pay(from int64, to int64, coin string, balance float64, address string) (err error) {
	//获取市场价格
	prices := caches.GetPrices()
	if prices == nil {
		return errors.New("获取市场价格出错")
	}
	price := prices.(map[string]float64)[coin]
	if price == 0 {
		return errors.New("价格异常")
	}
	//根据价格获取支付的币数量
	amount := balance / price
	//获取币种配信息
	coins := caches.GetCoin(coin)
	if coins == nil {
		return errors.New("ERROR")
	}
	fee := 0.0
	if coins["payFee"] != nil {
		fee = commons.GetFee(amount, t.Str(coins["payFee"]))
	}
	fmt.Println(fee)
	//数量加上手续费等于真实扣币
	amount += fee
	//var transacion services.Transacion
	//err = transacion.Convert(from, to, coin, "CNY", balance, amount)
	//if err != nil {
	//	panic("转帐失败")
	//}
	return
}
//生成新订单ID
func CreateOrderId() string {
	orderId := redis.CreateId("configs:weixin:orderid")
	return commons.ToStr(orderId)
}
//创建充值订单
func CreateOrder(orderId, uid int64, amount float64, _type int32, nonceStr string) error {
	history := models.PayHistory{ID_:primitive.NewObjectID(), OrderId:orderId, UID:uid, Create:commons.GetNowTime(), Coin:"FP", Amount:amount, Status:0, Type:_type, NonceStr:nonceStr}
	err := mongo.Client.Insert("payhistory", history)
	return err
}
//获取订单
func GetOrder(orderId int64) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := mongo.Client.FindOne("payhistory", bson.M{"orderId":orderId}, &result)
	return result, err
}
//查询订单
func QueryOrder(q interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := mongo.Client.FindOne("payhistory", q, &result)
	return result, err
}*/
package init

import (
	"fmt"
	"github.com/astaxie/beego/toolbox"
	"github.com/gogf/gf/g/encoding/gjson"
	"github.com/mongodb/mongo-go-driver/bson"
	"io/ioutil"
	"rubik/server/common/log"
	"rubik/server/common/middleware"
	"rubik/server/common/redis"
	"rubik/server/common/t"
	"rubik/server/goods/config"
	"rubik/server/goods/mongo"
	"rubik/server/goods/services"
	"time"
)


func init()  {
	if data, err := ioutil.ReadFile("goods.yml"); err != nil {
		panic("goods.yml 配置文件不存在")
	} else {
		dbUrl := gjson.New(string(data)).GetString("mongodb.url")
		redisMap := gjson.New(string(data)).GetMap("redis")
		rabbitMQUrl := gjson.New(string(data)).GetString("rabbitMQ.url")
		webChat := gjson.New(string(data)).GetMap("webchat")
		baiduAI := gjson.New(string(data)).GetMap("baiduAI")
		webMap := gjson.New(string(data)).GetMap("web")
		//初始化日志
		log.InitWarn("goods")
		//初始化数据库
		mongo.Init(dbUrl)
		//初始化redis
		redis.Init(redisMap)
		//获取rabbitmq
		config.RabbitMQUrl = rabbitMQUrl
		// webchat
		config.WebChat = webChat
		config.BaiduAI = baiduAI
		config.Web = webMap
		//设置中间件的安全码
		middleware.SetSecret(t.Str(webMap["secretKey"]))
	}

	//启动定时器
	Toolbox()
}

func Toolbox(){
	//秒钟：0-59、分钟：0-59、小时：1-23、日期：1-31、月份：1-12、星期：0-6（0 表示周日）
	tk := toolbox.NewTask("goodsCronTimer", "59 59 23 * * *", func() error {
		GoodsCronTimer(); return nil
	})
	toolbox.AddTask("goodsCronTimer", tk)
	toolbox.StartTask()
}

func GoodsCronTimer(){

	fmt.Println("--进入GoodsCronTimer--")
	var result []map[string]interface{}
	//mongodb商品列表

	err := mongo.Client.Aggregate("goods_order",[]bson.M{{"$match":bson.M{"status":int64(2),"endTime":bson.M{"$lt" : time.Now().Unix()}}}},&result);if err!=nil{
		fmt.Println("goods项目init方法err:",err)
		return
	}
	for _,v := range result{
		//修改订单状态
		err := services.UpdateGoodsOrder(v["goodsOrderId"].(int64) ,3 ,v["memo"].(string) ,-1 );if err!=nil{
			fmt.Println("err:",err.Error())
			return
		}

		//查看用户库存
		var result1 []map[string]interface{}
		err = mongo.Client.GoodsFind("user_goods_stock",bson.M{"uid":v["putUid"].(int64),"status":1,"goodsId":v["goodsId"].(int64)},&result1);if err!=nil{
			fmt.Println("err:",err.Error())
			return
		}
		//没有库存新增否则改变数量
		if len(result1) == 0{
			services.InsertGoodsStock(v["putUid"].(int64) ,v["goodsId"].(int64) ,v["goodsName"].(string) ,1 ,v["amount"].(int64))
		}else{
			services.UpdateGoodsStock(v["goodsId"].(int64),v["putUid"].(int64),v["amount"].(int64)+result1[0]["amount"].(int64))
		}
		fmt.Println("收货成功")
	}


}



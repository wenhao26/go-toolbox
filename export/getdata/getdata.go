package getdata

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"

	"toolbox/export/internal/mysql"
)

type Order struct {
	Id              int     `json:"id"`
	UserId          int     `json:"user_id"`
	OrderNo         string  `json:"order_no"`
	TradeNo         string  `json:"trade_no"`
	OriginalTradeNo string  `json:"original_trade_no"`
	ProductId       string  `json:"product_id"`
	ItemType        int     `json:"item_type"`
	PayTime         int     `json:"pay_time"`
	PayType         int     `json:"pay_type"`
	PayAmount       float32 `json:"pay_amount"`
	Coin            float32 `json:"coin"`
	GiveCoin        float32 `json:"give_coin"`
	Ip              string  `json:"ip"`
	CreateTime      int     `json:"create_time"`
}

type Orders struct {
	Data []*Order
	Page int
}

type Result struct {
	Page int
	Err  error
}

const sql = "select " +
	"user_id,order_no,trade_no,original_trade_no,product_id,item_type,pay_time,pay_type,pay_amount,coin,give_coin,ip,create_time " +
	"from iw_order where order_status=6 order by create_time limit ? offset ? "

type InputChan chan *Orders
type OutputChan chan *Result

type DataCMD func() InputChan
type DataPipeCMD func(input InputChan) OutputChan

// 管道
func Pipe(cmd1 DataCMD, cs ...DataPipeCMD) OutputChan {
	// cmd1执行完，获取到数据
	data := cmd1()
	out := make(OutputChan)
	wg := sync.WaitGroup{}
	for _, c := range cs {
		output := c(data) // 遍历上面的dataPipeCMD让每个管道函数去消化拿出来的数据，有可能是很耗时的
		// 上面消化完处理数据的同时。
		wg.Add(1)
		// 是异步执行的
		go func(outData OutputChan) { // 也需要开一个协程，去异步关闭
			defer wg.Done()
			for i := range outData { // 不断去获取是不是有数据，有的话就把他放入out中
				out <- i // 把执行的结果，放入out中，并返回
			}
		}(output)
	}
	go func() {
		defer close(out) // 在wait结束的时候，要关闭这个管道
		wg.Wait()        // 如果放到了外面，则会等待到所有的结果出来，才会return out，外面的进程就会卡死等待，所以需要放到协程中，异步进行等待
	}()
	return out // 执行的同时，需要先把out返回。因为执行是异步的，全部的通信全部在管道中
}

// 读取数据
func ReadData() InputChan {
	dbSvr := mysql.GetDb()
	page := 1
	limit := 1000

	result := make(InputChan)
	go func() {
		defer close(result)
		for {
			orders := &Orders{make([]*Order, 0), page}
			db := dbSvr.Raw(sql, limit, (page-1)*limit).Find(&orders.Data)
			if db.Error != nil || db.RowsAffected == 0 {
				break
			}
			result <- orders
			page++
		}
	}()
	return result
}

func ReadData1() InputChan {
	var sql string
	dbSvr := mysql.GetDb()
	finalId := 0
	field := "id,user_id,order_no,trade_no,original_trade_no,product_id,item_type,pay_time,pay_type,pay_amount,coin,give_coin,ip,create_time "
	times := 1
	limit := 10000

	result := make(InputChan)
	go func() {
		defer close(result)
		for {
			if finalId == 0 {
				sql = "select " + field + "from iw_order where order_status=6 order by create_time limit ? "
			} else {
				sql = "select " + field + "from iw_order where order_status=6 and id>" + strconv.Itoa(finalId) + " order by id limit ? "
			}

			orders := &Orders{make([]*Order, 0), times}
			db := dbSvr.Raw(sql, limit).Find(&orders.Data)
			if db.Error != nil || db.RowsAffected == 0 {
				break
			}
			finalId = orders.Data[len(orders.Data)-1].Id
			result <- orders
			times++
		}
	}()
	return result
}

// 写入数据
func WriteData(input InputChan) OutputChan {
	out := make(OutputChan)
	go func() {
		defer close(out)
		for data := range input {
			out <- &Result{
				Page: data.Page,
				Err:  SaveData(data),
			}
		}
	}()
	return out
}

// 保存到CSV
func SaveData(data *Orders) error {
	filename := fmt.Sprintf("csv/%d.csv", data.Page)
	csvFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	w := csv.NewWriter(csvFile)
	header := []string{"用户ID", "订单ID", " 交易号", "原始交易号", "商品ID", "商品类型", "支付时间", "支付类型", "支付金额", "金币", "赠送金币", "创建时间"}
	export := [][]string{
		header,
	}

	for _, d := range data.Data {
		row := []string{
			strconv.Itoa(d.UserId),
			d.OrderNo,
			d.TradeNo,
			d.OriginalTradeNo,
			d.ProductId,
			strconv.Itoa(d.ItemType),
			strconv.Itoa(d.PayTime),
			strconv.Itoa(d.PayType),
			strconv.FormatFloat(float64(d.PayAmount), 'f', -1, 32),
			strconv.FormatFloat(float64(d.Coin), 'f', -1, 32),
			strconv.FormatFloat(float64(d.GiveCoin), 'f', -1, 32),
			strconv.Itoa(d.CreateTime),
		}
		export = append(export, row)
	}

	err = w.WriteAll(export)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func Export() {
	// out := Pipe(ReadData, WriteData, WriteData) // 可以写多个WriteData来消化readData产生的数据
	out := Pipe(ReadData1, WriteData, WriteData)
	for res := range out {
		fmt.Printf("%d.csv文件执行完成，结果%v\n", res.Page, res.Err)
	}
}

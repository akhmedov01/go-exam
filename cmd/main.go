package main

import (
	"app/config"
	"app/controller"
	"app/storage/jsondb"
)

func main() {
	cfg := config.Load()
	strg, err := jsondb.NewConnectionJSON(&cfg)
	if err != nil {
		panic("Failed connect to json:" + err.Error())
	}
	con := controller.NewController(&cfg, strg)

	/* con.OrderGetList(&models.OrderGetListRequest{
		Offset: 1,
		Limit:  10,
	}) */

	//con.SearchOrdersFromDate("2022-07-16", "2022-08-16")

	//con.UserTotal("080b6453-d424-4362-b1e0-18b80caa6fce")

	//con.ProductSaleCount("a7ddbf9d-10b7-4429-93f8-712bd4074ca3")

	//con.UserHistory("e6ded598-675b-4de2-a1e9-00a876b8e719")

	//con.ReportCategory()

	//con.ActiveUser()

	//con.CreateBranch("Akhmedov")

	con.DateTopSales()

}

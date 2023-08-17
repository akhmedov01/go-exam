package controller

import (
	"app/models"
	"app/pkg/convert"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"
)

// 10-gacha bo'lgan tasklarni shop_cart.json ni order.json ga o'zgartirib ishlaganman
// Yani  "product_id": "7df81816-0b37-4922-bd87-d5dbe0f47c56",
//		 "user_id": "ebea6d88-820e-4863-8f69-e91f891b92b0",
//		 "count": 5,
//		 "status": false,
//		 "time": "2022-05-27 01:17:38"
// fieldlar order.json fieldlari

// 1. Order boyicha default holati time sort bolishi kerak. DESC

func (c *Controller) SortOrders() {
	resp, err := c.Strg.Order().GetList(&models.OrderGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting OrderList")
		return
	}

	sort.Slice(resp.Orders, func(i, j int) bool {

		time1, err := time.Parse("2006-01-02 15:04:05", resp.Orders[i].DateTime)

		if err != nil {
			fmt.Println("Error while parsing date")
			return false
		}

		time2, err := time.Parse("2006-01-02 15:04:05", resp.Orders[j].DateTime)

		if err != nil {
			fmt.Println("Error while parsing date")
			return false
		}

		return time1.After(time2)
	})

	for _, v := range resp.Orders {

		fmt.Printf("ProductID: %s UserID: %s Count: %d DateTime: %s Status: %t\n",
			v.ProductId, v.UserId, v.Count, v.DateTime, v.Status)

	}

}

// 2. Order Date boyicha filter qoyish kerak

func (c *Controller) SearchOrdersFromDate(fromDate, toDate string) {

	resp, err := c.Strg.Order().GetList(&models.OrderGetListRequest{})
	result := []*models.Order{}

	if err != nil {
		fmt.Println("Error while getting OrderList")
		return
	}

	parsingFromDate, _ := time.Parse("2006-01-02", fromDate)

	parsingToDate, _ := time.Parse("2006-01-02", toDate)

	for _, val := range resp.Orders {

		parsingSaleDate, _ := time.Parse("2006-01-02", val.DateTime[:10])

		if (parsingToDate.After(parsingSaleDate) && parsingFromDate.Before(parsingSaleDate)) || parsingFromDate == parsingSaleDate || parsingToDate == parsingSaleDate {

			order := val

			result = append(result, order)

		}

	}

	sort.Slice(result, func(i, j int) bool {
		time1, _ := time.Parse("2006-01-02 15:04:05", result[i].DateTime)
		time2, _ := time.Parse("2006-01-02 15:04:05", result[j].DateTime)
		return time1.After(time2)
	})

	for _, v := range result {

		fmt.Printf("ProductID: %s UserID: %s Count: %d DateTime: %s Status: %t\n",
			v.ProductId, v.UserId, v.Count, v.DateTime, v.Status)

	}
}

// 3. User history chiqish kerak. Ya'ni sotib olgan mahsulotlari korsatish kerak

func (c *Controller) UserHistory(id string) {

	resp, err := c.Strg.Order().GetList(&models.OrderGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting OrderList")
		return
	}

	userResp, err := c.Strg.User().GetById(&models.UserPrimaryKey{Id: id})

	if err != nil {
		fmt.Println("Error while getting UserList")
		return
	}

	respProduct, err := c.Strg.Product().GetList(&models.ProductGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting ProductList")
		return
	}

	resultMap := make(map[string]map[string]int)
	productMap := make(map[string]*models.Product)
	testMap := make(map[string]int)
	orderMap := make(map[string]*models.Order)

	for _, v := range respProduct.Products {
		productMap[v.Id] = v
	}

	for _, v := range resp.Orders {

		if v.Status {

			testMap[v.ProductId] = v.Count
			resultMap[v.Id] = testMap

		}

	}

	for _, v := range resp.Orders {

		orderMap[v.Id] = v

	}

	for id, value := range resultMap {

		if orderMap[id].UserId == userResp.Id {

			fmt.Printf("User Name: %s\n", userResp.Name)

			for i, v := range value {

				fmt.Printf("Product Name: %s\tPrice: %d\tCount: %d\tTotal: %d\t Time: %s\n",
					productMap[i].Name, productMap[i].Price, v, productMap[i].Price*v, orderMap[id].DateTime)

			}
		}
	}

}

// 4. User qancha pul mahsulot sotib olganligi haqida hisobot.

func (c *Controller) UserTotal(id string) {

	respOrder, err := c.Strg.Order().GetList(&models.OrderGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting OrderList")
		return
	}

	respProduct, err := c.Strg.Product().GetList(&models.ProductGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting ProductList")
		return
	}

	respUser, err := c.Strg.User().GetById(&models.UserPrimaryKey{Id: id})

	if err != nil {
		fmt.Println("Error while getting UserList")
		return
	}

	productMap := make(map[string]*models.Product)
	counterMap := make(map[string]int)

	for _, v := range respProduct.Products {
		productMap[v.Id] = v
	}

	for _, v := range respOrder.Orders {

		if v.Status {

			counterMap[v.UserId] += productMap[v.ProductId].Price * v.Count

		}

	}

	fmt.Printf("Name: %s\tTotal Buy Price: %d\n", respUser.Name, counterMap[id])

}

// 5. Productlarni Qancha sotilgan boyicha hisobot

func (c *Controller) ProductSaleCount(id string) {

	respProduct, err := c.Strg.Product().GetById(&models.ProductPrimaryKey{Id: id})

	if err != nil {
		fmt.Println("Error while getting ProductList")
		return
	}

	respOrder, err := c.Strg.Order().GetList(&models.OrderGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting OrderList")
		return
	}

	var counter int

	for _, v := range respOrder.Orders {
		if respProduct.Id == v.ProductId && v.Status {
			counter += v.Count
		}
	}

	fmt.Printf("Product Name: %s\t Count: %d\n", respProduct.Name, counter)

}

// 6. Top 10 ta sotilayotgan mahsulotlarni royxati.

func (c *Controller) TopHighSaleProducts() {

	respProduct, err := c.Strg.Product().GetList(&models.ProductGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting ProductList")
		return
	}

	respOrder, err := c.Strg.Order().GetList(&models.OrderGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting OrderList")
		return
	}

	productMap := make(map[string]*models.Product)
	counterMap := make(map[string]int)

	for _, v := range respProduct.Products {

		productMap[v.Id] = v

	}

	for _, v := range respOrder.Orders {

		if v.Status {

			counterMap[v.ProductId] += v.Count

		}
	}

	type Count struct {
		Key   string
		Value int
	}

	var counts = []Count{}

	for i, v := range counterMap {
		counts = append(counts, Count{i, v})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Value > counts[j].Value
	})

	counts = counts[:10]

	for _, v := range counts {
		fmt.Printf("Product Name: %s -> Count: %d \n", productMap[v.Key].Name, v.Value)
	}

}

// 7. TOP 10 ta Eng past sotilayotgan mahsulotlar royxati

func (c *Controller) TopLowSaleProducts() {

	respProduct, err := c.Strg.Product().GetList(&models.ProductGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting ProductList")
		return
	}

	respOrder, err := c.Strg.Order().GetList(&models.OrderGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting OrderList")
		return
	}

	productMap := make(map[string]*models.Product)
	counterMap := make(map[string]int)

	for _, v := range respProduct.Products {

		productMap[v.Id] = v

	}

	for _, v := range respOrder.Orders {

		if v.Status {

			counterMap[v.ProductId] += v.Count

		}

	}

	type Count struct {
		Key   string
		Value int
	}

	var counts = []Count{}

	for i, v := range counterMap {
		counts = append(counts, Count{i, v})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Value < counts[j].Value
	})

	counts = counts[:10]

	for _, v := range counts {
		fmt.Printf("Product Name: %s -> Count: %d \n", productMap[v.Key].Name, v.Value)
	}

}

// 8. Qaysi Sanada eng kop mahsulot sotilganligi boyicha jadval

func (c *Controller) DateTopSales() {

	respOrder, err := c.Strg.Order().GetList(&models.OrderGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting OrderList")
		return
	}

	respProduct, err := c.Strg.Product().GetList(&models.ProductGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting ProductList")
		return
	}

	dateMap := make(map[string]int)
	productMap := make(map[string]*models.Product)
	resultMap := make(map[string]int)

	for _, v := range respProduct.Products {

		productMap[v.Id] = v

	}

	// dateMap ning key iga sanalar soatlardan ajratilib, value siga esa shu sanada sotilgan maxsulotlar saqlanadi

	for _, v := range respOrder.Orders {

		if v.Status {

			parsingDate := v.DateTime[:10]

			dateMap[parsingDate] += v.Count

		}
	}

	type Count struct {
		Key   string
		Value int
	}

	var counts = []Count{}

	for i, v := range dateMap {
		counts = append(counts, Count{i, v})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Value > counts[j].Value
	})

	// Eng ko'p savdo bo'lgan kun counts ga sana va shu sanada nechta maxsulot sotilganni key, value tarzida joylandi

	counts = counts[:1]

	for _, v := range counts {
		for _, p := range respOrder.Orders {

			parsingDate := p.DateTime[:10]

			if v.Key == parsingDate {

				resultMap[p.DateTime] = p.Count

			}
		}
	}

	// Shu sananing soatlarida sotilgan maxsulotlar count bo'yicha DESC sort lanadi

	type ProductCount struct {
		Key   string
		Value int
	}

	var productCounts = []ProductCount{}

	for i, v := range resultMap {
		productCounts = append(productCounts, ProductCount{i, v})
	}

	sort.Slice(productCounts, func(i, j int) bool {
		return productCounts[i].Value > productCounts[j].Value
	})

	for _, v := range productCounts {
		for _, orderV := range respOrder.Orders {

			if v.Key == orderV.DateTime {

				fmt.Printf("Product Name: %s   \t Date: %s \t Count: %d\n",
					productMap[orderV.ProductId].Name, orderV.DateTime, v.Value)

			}

		}
	}

}

// 9. Qaysi category larda qancha mahsulot sotilgan boyicha jadval

func (c *Controller) ReportCategory() {

	respProduct, err := c.Strg.Product().GetList(&models.ProductGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting ProductList")
		return
	}

	respCategory, err := c.Strg.Category().GetList(&models.CategoryGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting ProductList")
		return
	}

	respOrder, err := c.Strg.Order().GetList(&models.OrderGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting OrderList")
		return
	}

	categoryMap := make(map[string]*models.Category)
	productMap := make(map[string]*models.Product)
	orderMap := make(map[string]int)

	for _, v := range respCategory.Categorys {

		categoryMap[v.Id] = v

	}

	for _, v := range respProduct.Products {

		productMap[v.Id] = v

	}

	for _, v := range respOrder.Orders {

		if v.Status {

			orderMap[productMap[v.ProductId].CategoryID] += v.Count

		}

	}

	type Count struct {
		Key   string
		Value int
	}

	var counts = []Count{}

	for i, v := range orderMap {
		counts = append(counts, Count{i, v})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Value > counts[j].Value
	})

	for _, v := range counts {
		fmt.Printf("Category Name: %s \t Count: %d\n", categoryMap[v.Key].Name, v.Value)
	}

}

// 10. Qaysi User eng Active xaridor. Bitta ma'lumot chiqsa yetarli.

func (c *Controller) ActiveUser() {

	respUser, err := c.Strg.User().GetList(&models.UserGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting UserList")
		return
	}

	respOrder, err := c.Strg.Order().GetList(&models.OrderGetListRequest{})

	if err != nil {
		fmt.Println("Error while getting OrderList")
		return
	}

	userMap := make(map[string]*models.User)
	resultMap := make(map[string]int)

	for _, v := range respUser.Users {

		userMap[v.Id] = v

	}

	for _, v := range respOrder.Orders {

		if v.Status {

			resultMap[v.UserId] += v.Count

		}

	}

	type Count struct {
		Key   string
		Value int
	}

	var counts = []Count{}

	for i, v := range resultMap {
		counts = append(counts, Count{i, v})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Value > counts[j].Value
	})

	counts = counts[:1]

	for _, v := range counts {
		fmt.Printf("User Name: %s \tCount: %d\n", userMap[v.Key].Name, v.Value)
	}

}

// 11. Agar User 9 dan kop mahuslot sotib olgan bolsa, 1 tasi tekinga beriladi va 9 ta uchun pul hisoblanadi. 1 tasi eng arzon mahsulotni pulini hisoblamaysiz.

// Faqat shu taskda eski order bo'yicha logika yozilgan.
// Qolgan tasklarda siz tashlagan data.zip dagi shop_cart.json ni order.json ga rename qilib json-DB.zip ning data sidagi
// order.json ni o'rniga qo'yganman. O'rniga ko'ygan datamdagi field lar bilan bu masalani yechib bo'lmas ekan.
// Shuning uchun eski order.json dagi fieldlarni ishlatdim.

func (c *Controller) OrderPayment(id string) error {

	order, err := c.Strg.Order().GetById(&models.OrderPrimaryKey{Id: id})
	if err != nil {
		log.Printf("error while Order => GetById: %+v\n", err)
		return err
	}

	user, err := c.Strg.User().GetById(&models.UserPrimaryKey{Id: order.UserId})
	if err != nil {
		log.Printf("error while User => GetById: %+v\n", err)
		return err
	}

	if len(order.OrderItems) > 9 {

		respProduct, err := c.Strg.Product().GetList(&models.ProductGetListRequest{})

		if err != nil {
			fmt.Println("Error while getting ProductList")
			return err
		}

		productMap := make(map[string]*models.Product)
		counterMap := make(map[string]int)

		for _, v := range respProduct.Products {

			productMap[v.Id] = v

		}

		for _, v := range order.OrderItems {

			counterMap[v.ProductId] = productMap[v.ProductId].Price

		}

		type Count struct {
			Key   string
			Value int
		}

		var counts = []Count{}

		for i, v := range counterMap {
			counts = append(counts, Count{i, v})
		}

		sort.Slice(counts, func(i, j int) bool {
			return counts[i].Value < counts[j].Value
		})

		counts = counts[:1]

		for _, v := range counts {
			order.Sum -= v.Value
		}

	}

	if order.Sum > user.Balance {
		return errors.New("Not enough balance " + user.Name + " " + user.Surname)
	}

	order.Status = true
	fmt.Println("suc")
	user.Balance -= order.Sum

	var updateOrder models.UpdateOrder
	err = convert.ConvertStructToStruct(order, &updateOrder)
	if err != nil {
		log.Printf("error while convertStructToStruct: %+v\n", err)
		return err
	}

	_, err = c.Strg.Order().Update(&updateOrder)
	if err != nil {
		log.Printf("error while order => Update: %+v\n", err)
		return err
	}

	var updateUser models.UpdateUser
	err = convert.ConvertStructToStruct(user, &updateUser)
	if err != nil {
		log.Printf("error while convertStructToStruct: %+v\n", err)
		return err
	}

	_, err = c.Strg.User().Update(&updateUser)
	if err != nil {
		log.Printf("error while User => Update: %+v\n", err)
		return err
	}

	return nil
}

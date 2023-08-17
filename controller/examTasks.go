package controller

import (
	"app/models"
	"fmt"
	"sort"
	"time"
)

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
	//orderMap := make(map[string]*models.Order)

	for _, v := range respProduct.Products {
		productMap[v.Id] = v
	}

	for _, v := range resp.Orders {

		if v.Status {

			testMap[v.ProductId] = v.Count
			resultMap[v.UserId] = testMap

		}

	}

	/* for _, v := range resp.Orders {

		orderMap[v.UserId] = v

	} */

	for userId, value := range resultMap {

		if userId == userResp.Id {

			fmt.Printf("User Name: %s\n", userResp.Name)

			for i, v := range value {

				fmt.Printf("Product Name: %s\tPrice: %d\tCount: %d\tTotal: %d\n",
					productMap[i].Name, productMap[i].Price, v, productMap[i].Price*v)

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
	resultMap := make(map[string]map[string]int)

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

	counts = counts[:1]

	for _, v := range respProduct.Products {

		productMap[v.Id] = v

	}

	for _, v := range respOrder.Orders {

		resultMap[v.DateTime][v.ProductId] += v.Count

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

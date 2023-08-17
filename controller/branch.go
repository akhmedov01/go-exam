package controller

import (
	"app/models"
	"fmt"
)

func (c *Controller) CreateBranch(name string) {

	resp, err := c.Strg.Branch().CreateBranch(models.CreateBranch{
		Name: name,
	})

	if err != nil {
		fmt.Println("error from CreateBranch:", err.Error())
		return
	}

	fmt.Println("created new branch with id:", resp)

}

func (c *Controller) UpdateBranch(id, name string) {

	resp, err := c.Strg.Branch().UpdateBranch(models.Branch{
		Id:   id,
		Name: name,
	})

	if err != nil {
		fmt.Println("error from UpdateBranch:", err.Error())
		return
	}

	fmt.Printf(" %+v\n", resp)

}

func (c *Controller) GetBranch(id string) {

	resp, err := c.Strg.Branch().GetBranch(models.IdRequest{Id: id})

	if err != nil {
		fmt.Println("error from GetBranch:", err.Error())
		return
	}

	fmt.Printf(" %+v\n", resp)

}

func (c *Controller) GetAllBranch(page, limit int) {

	resp, err := c.Strg.Branch().GetAllBranch(models.GetAllBranchRequest{
		Page:  page,
		Limit: limit,
	})

	if err != nil {
		fmt.Println("error from GetAllBranch:", err.Error())
		return
	}

	fmt.Printf(" %+v\n", resp)

}

func (c *Controller) DeleteBranch(id string) {

	resp, err := c.Strg.Branch().DeleteBranch(models.IdRequest{Id: id})

	if err != nil {
		fmt.Println("error from DeleteBranch:", err.Error())
		return
	}

	fmt.Println("id:", resp)

}

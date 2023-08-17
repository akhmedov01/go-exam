package models

type Branch struct {
	Id   string
	Name string
}

type CreateBranch struct {
	Name string
}

type IdRequest struct {
	Id string
}

type GetAllBranch struct {
	Branches []Branch
	Count    int
}

type GetAllBranchRequest struct {
	Page  int
	Limit int
}

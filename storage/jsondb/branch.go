package jsondb

import (
	"app/models"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/google/uuid"
)

type branchRepo struct {
	fileName string
}

func NewBranchRepo(fileName string) *branchRepo {

	return &branchRepo{fileName: fileName}

}

func (b *branchRepo) CreateBranch(req models.CreateBranch) (string, error) {

	branches, err := b.read()

	if err != nil {
		return "error", err
	}

	var id = uuid.NewString()

	branches = append(branches, models.Branch{
		Id:   id,
		Name: req.Name,
	})

	err = b.write(branches)

	if err != nil {
		return "error", err
	}

	return id, nil

}

func (b *branchRepo) UpdateBranch(req models.Branch) (string, error) {

	branches, err := b.read()

	if err != nil {
		return "error", err
	}

	for i, v := range branches {
		if v.Id == req.Id {

			branches[i].Name = req.Name

			err = b.write(branches)

			if err != nil {
				return "error", err
			}

			return "updated suc", nil
		}
	}

	return "", errors.New("Not found")
}

func (b *branchRepo) GetBranch(req models.IdRequest) (models.Branch, error) {

	branches, err := b.read()

	if err != nil {
		return models.Branch{}, err
	}

	for _, v := range branches {
		if v.Id == req.Id {
			return v, nil
		}
	}

	return models.Branch{}, errors.New("Not found")
}

func (b *branchRepo) GetAllBranch(req models.GetAllBranchRequest) (resp models.GetAllBranch, err error) {

	branches, err := b.read()

	if err != nil {
		return models.GetAllBranch{}, err
	}

	start := req.Limit * (req.Page - 1)
	end := start + req.Limit

	if start > len(branches) {

		resp.Branches = []models.Branch{}
		resp.Count = len(branches)
		return resp, nil

	} else if end > len(branches) {

		resp.Branches = branches[start:]
		resp.Count = len(branches)
		return resp, nil

	}

	return models.GetAllBranch{Branches: branches[start:end], Count: len(branches)}, nil
}

func (b *branchRepo) DeleteBranch(req models.IdRequest) (string, error) {

	branches, err := b.read()

	if err != nil {
		return "", err
	}

	for i, v := range branches {
		if v.Id == req.Id {
			if i == (len(branches) - 1) {
				branches = branches[:i]

				err = b.write(branches)

				if err != nil {
					return "error", err
				}

				return "deleted suc", nil
			} else {
				branches = append(branches[:i], branches[i+1:]...)

				err = b.write(branches)

				if err != nil {
					return "error", err
				}

				return "deleted suc", nil
			}
		}
	}

	return "", errors.New("Not found")
}

func (u *branchRepo) read() ([]models.Branch, error) {
	var (
		branches []models.Branch
	)

	data, err := os.ReadFile(u.fileName)
	if err != nil {
		log.Printf("Error while Read data: %+v", err)
		return nil, err
	}

	err = json.Unmarshal(data, &branches)
	if err != nil {
		log.Printf("Error while Unmarshal data: %+v", err)
		return nil, err
	}

	return branches, nil
}

func (u *branchRepo) write(branches []models.Branch) error {

	body, err := json.Marshal(branches)
	if err != nil {
		return err
	}

	err = os.WriteFile(u.fileName, body, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

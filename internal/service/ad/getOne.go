package ad

import (
	"ads-service/internal/models"
)

func (s *service) GetOne(id int) (*models.Ad, error) {
	ad, err := s.store.Ad.GetOne(id)
	if err != nil {
		s.logger.Error("unable to get ad", "err", err.Error())
		return nil, err
	}

	return ad, nil
}

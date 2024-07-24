package ad

import (
	"ads-service/internal/models"
)

func (s *service) GetAll(priceSort string, dateSort string, page int) ([]*models.Ad, error) {
	ads, err := s.store.Ad.GetAll(priceSort, dateSort, page)
	if err != nil {
		s.logger.Error("unable to get ads", "err", err.Error())
		return nil, err
	}

	return ads, nil
}

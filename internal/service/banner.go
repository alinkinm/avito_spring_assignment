package service

import (
	"avito2/internal/core"
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type BannerCache interface {
	Get(key string) (string, error)
	Set(key, value string, expiration time.Duration) error
}

type BannerRepository interface {
	GetUserBanner(ctx context.Context, banner *core.BannerRequest) (*core.Banner, error)
	GetAllBanners(ctx context.Context, banner *core.BannerRequest2) ([]*core.Banner, error)
}

type BannerService struct {
	BannerRepository BannerRepository
	Cache            BannerCache
}

func NewBannerService(bannerRepository BannerRepository, cache BannerCache) *BannerService {

	return &BannerService{
		BannerRepository: bannerRepository,
		Cache:            cache,
	}
}

func (service *BannerService) GetUserBanner(ctx context.Context, banner *core.BannerRequest, role string) (string, error) {

	if !banner.UseLastRevision {

		key := strconv.Itoa(banner.FeatureID) + strconv.Itoa(banner.TagID)
		cachedBanner, err := service.Cache.Get(key)
		if err != nil {
			log.Printf("Ошибка получения баннера из кэша: %v", err)
			return "", err
		}

		if cachedBanner != "" { //баннер найден в кэше
			result := core.Banner{}
			json.Unmarshal([]byte(cachedBanner), &result)
			if role == "user" && !result.IsActive {
				return "", core.NewErrAccessDenied()
			}
			return result.Banner, nil
		}

		banner, err := service.BannerRepository.GetUserBanner(ctx, banner)

		if err != nil {

			if err == sql.ErrNoRows {
				return banner.Banner, core.NewErrBannerDoesNotExist()
			} else {
				return banner.Banner, core.NewErrInternalServerError()
			}
		}

		err = service.Cache.Set(key, banner.Banner, 5*time.Minute)
		if err != nil {
			log.Printf("Ошибка кэширования баннера: %v", err)
		}

		return banner.Banner, nil

	} else {
		banner, err := service.BannerRepository.GetUserBanner(ctx, banner)

		if err != nil {

			if err == sql.ErrNoRows {
				return banner.Banner, core.NewErrBannerDoesNotExist()
			} else {
				return banner.Banner, core.NewErrInternalServerError()
			}
		}
	}

	return "", nil
}

func (service *BannerService) GetAllBanners(ctx context.Context, banner *core.BannerRequest2) ([]*core.Banner, error) {

	var banners []*core.Banner

	banners, err := service.BannerRepository.GetAllBanners(ctx, banner)
	if err != nil {
		return nil, err
	}

	banners = removeDuplicates(banners)

	return banners, nil
}

func removeDuplicates(arr []*core.Banner) []*core.Banner {

	unique := make(map[string]bool)
	var uniqueArr []*core.Banner

	for _, v := range arr {

		if !unique[v.Banner] {
			unique[v.Banner] = true
			uniqueArr = append(uniqueArr, v)
		}
	}

	return uniqueArr
}

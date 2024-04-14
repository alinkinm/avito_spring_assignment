package repository

import (
	"avito2/internal/core"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	selectUserBannerQuery = `SELECT doc FROM (banner join banner_feature_tag) where banner.feature_id = ($1) AND banner_feature_tag.tag_id = ($2)
	AND is_active = true;`
	selectAdminBannerQuery = `SELECT doc FROM (banner join banner_feature_tag) where banner.feature_id = ($1) AND banner_feature_tag.tag_id = ($2);`

	GetAllFeatureTag = `SELECT doc FROM banner
	WHERE id = (
		SELECT banner_id FROM banner_feature_tag
		WHERE feature_id = (SELECT id FROM feature WHERE system_id = ($1))
		AND tag_id = (SELECT id FROM tag WHERE system_id = ($2))
	)`

	GetAllTag1 = `SELECT b.system_id AS banner_system_id, 
	b.doc AS banner_doc, 
	f.system_id AS feature_system_id, 
	t.system_id AS tag_system_id, 
	b.created_at AS banner_created_at, 
	b.updated_at AS banner_updated_at
	FROM banner AS b
	JOIN banner_feature_tag AS bft ON bft.banner_id = b.id
	JOIN tag AS t ON bft.tag_id = t.id
	JOIN feature AS f ON bft.feature_id = f.id
	WHERE t.system_id = $1`

	selectAllFeatureQuery = `SELECT b.system_id AS banner_system_id, 
	b.doc AS banner_doc, 
	f.system_id AS feature_system_id, 
	t.system_id AS tag_system_id, 
	b.created_at AS banner_created_at, 
	b.updated_at AS banner_updated_at
	FROM banner AS b
	JOIN banner_feature AS bf ON bf.banner_id = b.id
	JOIN feature AS f ON bf.feature_id = f.id
	JOIN banner_feature_tag AS bft ON bft.banner_id = b.id AND bft.feature_id = f.id
	JOIN tag AS t ON bft.tag_id = t.id
	WHERE f.system_id = $1`
)

type BannerRepository struct {
	Db *sqlx.DB
}

func NewBannerRepository(db *sqlx.DB) *BannerRepository {
	return &BannerRepository{Db: db}
}

func (repository *BannerRepository) GetUserBanner(ctx context.Context, banner *core.BannerRequest) (*core.Banner, error) {
	var banner1 string

	err := repository.Db.QueryRowContext(ctx, banner1, selectUserBannerQuery).Scan(&banner)
	if err != nil {
		return nil, err
	}

	return &core.Banner{Banner: banner1}, nil
}

func (repository *BannerRepository) GetAdminBanner(ctx context.Context, banner *core.BannerRequest) (*core.Banner, error) {
	var banner1 string

	err := repository.Db.QueryRowContext(ctx, banner1, selectAdminBannerQuery).Scan(&banner)
	if err != nil {
		return nil, err
	}

	return &core.Banner{Banner: banner1}, nil
}

func (repository *BannerRepository) GetAllBanners(ctx context.Context, banner *core.BannerRequest2) ([]*core.Banner, error) {

	var banners []*core.Banner

	query := createQuery(&banner)

	rows, err := repository.Db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		banner := &core.Banner{}
		err = rows.Scan(&banner.Banner)
		if err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

func createQuery(banner core.BannerRequest2) string {

	query := ""

	if banner.FeatureID != 0 && banner.TagID != 0 {
		query = GetAllFeatureTag
	} else if banner.FeatureID != 0 {
		query = selectAllFeatureQuery
	} else {
		query = GetAllTag1
	}

	if banner.Offset != 0 {
		query += fmt.Sprintf(`OFFSET %d`, banner.Offset)
	}

	if banner.Limit != 0 {
		query += fmt.Sprintf(`LIMIT %d`, banner.Limit)
	}

	query += `;`

	return query

}

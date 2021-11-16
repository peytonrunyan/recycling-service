package sqlserver

import (
	"database/sql"
	"fmt"

	"recycling-service/pkg/models"
)

// Model for MaterialGuidelinesResults
type MGRModel struct {
	DB *sql.DB
}

func (m *MGRModel) Get(cID string) (*models.MaterialGuidelineResults, error) {
	return nil, nil
}

func (m *MGRModel) GetAll() (map[string]*[]models.MaterialGuidelineResults, error) {
	query := `
    SELECT top 100 [m_id]
          ,[community_id]
          ,[category]
          ,[yes_no]
          ,[category_type]
          ,[material]
      FROM [dbo].[material_guideline_results]
    `
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	mgResults := make(map[string]*[]models.MaterialGuidelineResults)
	for rows.Next() {
		m := models.MaterialGuidelineResults{}
		err = rows.Scan(m.MID, m.CommunityID, m.Category, m.YesNo, m.CategoryType, m.Material)
		fmt.Printf("m.CommunityID: %v\n", m.CommunityID)
		if err != nil {
			return nil, err
		}
		communityResults, ok := mgResults[m.CommunityID]
		if ok {
			*communityResults = append(*communityResults, m)
		} else {
			mgResults[m.CommunityID] = &[]models.MaterialGuidelineResults{m}
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return mgResults, nil
}

package dto

type ListRoleRequest struct {
	Limit        int    `query:"per_page"`
	Page         int    `query:"page"`
	SearchColumn string `query:"search_column"`
	SearchValue  string `query:"search_value"`
	Sort         string `query:"sort"`
	SortBy       string `query:"sort_by"`
}

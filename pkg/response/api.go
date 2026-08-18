package response

/*
Single Resource Response

	{
	  "success": true,
	  "data": {
	    "id": "usr_928310",
	    "firstName": "Alex",
	    "lastName": "Rivera",
	    "email": "alex@company.com"
	  },
	  "meta": {
	    "cached": false
	  },
	  "errors": null
	}

Collection (Paginated List) Response

	{
	  "success": true,
	  "data": [
	    { "id": "tr_102", "amount": 49.99 },
	    { "id": "tr_103", "amount": 120.00 }
	  ],
	  "meta": {
	    "pagination": {
	      "currentPage": 1,
	      "pageSize": 25,
	      "totalRecords": 240,
	      "totalPages": 10
	    }
	  },
	  "errors": null
	}

Error Response

	{
	  "success": false,
	  "data": null,
	  "meta": null,
	  "errors": [
	    {
	      "code": "INVALID_INPUT",
	      "message": "The email address domain provided is invalid.",
	      "field": "email"
	    }
	  ]
	}
*/

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta"`
	Errors  []ErrorItem `json:"errors"`
}

type ErrorItem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

type Meta struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}

func Success(data interface{}, meta *Meta) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
		Errors:  nil,
	}
}

func Error(code, message, field string) APIResponse {
	return APIResponse{
		Success: false,
		Data:    nil,
		Meta:    nil,
		Errors: []ErrorItem{
			{Code: code, Message: message, Field: field},
		},
	}
}

func Errors(errors []ErrorItem) APIResponse {
	return APIResponse{
		Success: false,
		Data:    nil,
		Meta:    nil,
		Errors:  errors,
	}
}

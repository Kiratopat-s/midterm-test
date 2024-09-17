package model

import "github.com/Kiratopat-s/workflow/internal/constant"

// Request to create a new item
type RequestCreateItem struct {
    Title    string  `json:"title" binding:"required"`
    Amount   int     `json:"amount" binding:"required"`
    Quantity int     `json:"quantity" binding:"required"`
}

// Request to update an existing item
type RequestUpdateItem struct {
    Title    *string `json:"title"`
    Amount   *int    `json:"amount"`
    Quantity *int    `json:"quantity"`
}

// Request to change item status
type RequestPatchItemStatus struct {
    Status constant.ItemStatus `json:"status" binding:"required"`
}

// Request to find items by status
type RequestFindItem struct {
    Status constant.ItemStatus `form:"status"`
    ItemID int                 `form:"item_id"`
}

// Request for user login
type RequestLogin struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// Response for item creation
type ResponseCreateItem struct {
    ID       int             `json:"id"`
    Amount   int             `json:"amount"`
    Quantity int             `json:"quantity"`
    Status   constant.ItemStatus `json:"status"`
    OwnerID  int             `json:"owner_id"`
}

// Response for getting a single item
type ResponseGetItem struct {
    ID       int             `json:"id"`
    Title    string          `json:"title"`
    Amount   int             `json:"amount"`
    Quantity int             `json:"quantity"`
    Status   constant.ItemStatus `json:"status"`
    OwnerID  int             `json:"owner_id"`
}

// Response for listing all items
type ResponseListItems struct {
    Items []ResponseGetItem `json:"items"`
}

// Response for login
type ResponseLogin struct {
    Message string `json:"message"`
}

// Response for logout
type ResponseLogout struct {
    Message string `json:"message"`
}

type RequestPatchManyItemStatus struct {
	IDs    []int    `json:"ids" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type RequestDeleteManyItems struct {
    IDs []int `json:"ids" binding:"required"`
}
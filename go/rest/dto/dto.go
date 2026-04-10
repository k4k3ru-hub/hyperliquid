//
// dto.go
//
package dto

import (

)



type RequestBody struct {
    Type string `json:"type"`
    User string `json:"user,omitempty"`
    DEX  string `json:"dex,omitempty"`
}


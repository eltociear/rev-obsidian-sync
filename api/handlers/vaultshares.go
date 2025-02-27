package handlers

import (
	"github.com/acheong08/obsidian-sync/database"
	"github.com/acheong08/obsidian-sync/utilities"
	"github.com/gin-gonic/gin"
)

func InviteVaultShare(c *gin.Context) {
	var req struct {
		Email   string `json:"email"`
		Token   string `json:"token" binding:"required"`
		VaultID string `json:"vault_uid"`
	}
	if c.BindJSON(&req) != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request",
		})
		return
	}
	email, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid token",
		})
		return
	}

	// Make invite
	db := c.MustGet("db").(*database.Database)

	if !db.HasAccessToVault(email, req.VaultID) {
		c.JSON(401, gin.H{
			"error": "You do not have access to this vault",
		})
	}

	user, err := db.UserInfo(req.Email)
	if err != nil {
		c.JSON(200, gin.H{
			"error": "User does not exist",
		})
	}

	err = db.ShareVaultInvite(req.Email, user.Name, req.VaultID)
	if err != nil {
		c.JSON(200, gin.H{
			"error": "Failed to share vault",
		})
	}

	c.JSON(200, gin.H{})
}

func ListVaultShares(c *gin.Context) {
	var req struct {
		Token   string `json:"token" binding:"required"`
		VaultID string `json:"vault_uid"`
	}
	if c.BindJSON(&req) != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request",
		})
		return
	}
	email, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid token",
		})
		return
	}
	db := c.MustGet("db").(*database.Database)
	if !db.HasAccessToVault(email, req.VaultID) {
		c.JSON(401, gin.H{
			"error": "You do not have access to this vault",
		})
		return
	}
	shares, err := db.GetVaultShares(req.VaultID)
	if err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"shares": shares,
	})
}

func RemoveVaultShare(c *gin.Context) {
	var req struct {
		ShareUID string `json:"share_uid"`
		Token    string `json:"token" binding:"required"`
		VaultUID string `json:"vault_uid" binding:"required"`
	}
	if c.BindJSON(&req) != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request",
		})
		return
	}
	email, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid token",
		})
		return
	}
	db := c.MustGet("db").(*database.Database)
	if req.ShareUID != "" {
		if !db.IsVaultOwner(req.VaultUID, email) {
			c.JSON(401, gin.H{
				"error": "You are not the owner of this vault",
			})
			return
		}
	} else {
		if !db.HasAccessToVault(email, req.VaultUID) {
			c.JSON(401, gin.H{
				"error": "You do not have access to this vault",
			})
		}
	}
	err = db.ShareVaultRevoke(req.ShareUID, req.VaultUID, email)
	if err != nil {
		c.JSON(200, gin.H{
			"error": "Failed to remove vault share",
		})
	}
	c.JSON(200, gin.H{})

}

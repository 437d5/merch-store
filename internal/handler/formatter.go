package handler

import (
	"github.com/437d5/merch-store/internal/inventory"
	"github.com/437d5/merch-store/internal/transactions"
	"github.com/gin-gonic/gin"
)

func formatTranscations(transactions []transactions.Transaction, userId int) gin.H {
	var received, sent []gin.H

	for _, t := range transactions {
		if t.FromUser == userId {
			sent = append(sent, gin.H{
				"toUser": t.ToUser,
				"amount": t.Amount,
			})
		} else {
			received = append(received, gin.H{
				"fromUser": t.FromUser,
				"amount": t.Amount,
			})
		}
	}

	return gin.H{
		"received": received,
		"sent": sent,  
	}
}

func formatInventory(inventory inventory.Inventory) []gin.H {
	var items []gin.H

	for _, i := range inventory.Items {
		items = append(items, gin.H{
			"type": i.ItemType,
			"quantity": i.Quantity,
		})
	}

	return items
}
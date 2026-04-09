package handlers

import (
	"io"
	"strconv"
	"strings"
	"time"

	"hookfy/config"
	"hookfy/models"

	"github.com/gin-gonic/gin"
)

const defaultTTL = 24 * time.Hour

func CreateWebhook(c *gin.Context) {
	hash := c.Param("hash")

	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	method := c.Request.Method
	url := c.Request.RequestURI
	contentType := c.ContentType()
	remoteAddr := c.ClientIP()

	query := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		query[k] = strings.Join(v, ", ")
	}

	rawBody, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	body := map[string]string{"raw": string(rawBody)}

	webhook := models.Webhook{
		Hash:        hash,
		Method:      method,
		URL:         url,
		Headers:     headers,
		Body:        body,
		QueryString: query,
		ContentType: contentType,
		RemoteAddr:  remoteAddr,
		ExpiresAt:   time.Now().Add(defaultTTL),
	}

	if err := config.DB.Create(&webhook).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to save webhook"})
		return
	}
	c.JSON(200, gin.H{"message": "webhook received", "id": webhook.ID})
}

func GetInbox(c *gin.Context) {
	format := c.DefaultQuery("type", "json")
	hash := c.Query("hash")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	psize, _ := strconv.Atoi(c.DefaultQuery("psize", "20"))
	if page < 1 {
		page = 1
	}
	if psize < 1 || psize > 100 {
		psize = 20
	}
	offset := (page - 1) * psize

	var total int64
	q := config.DB.Model(&models.Webhook{}).Where("expires_at > ?", time.Now())
	if hash != "" {
		q = q.Where("hash = ?", hash)
	}
	q.Count(&total)

	var webhooks []models.Webhook
	q.Order("created_at DESC").Limit(psize).Offset(offset).Find(&webhooks)

	if format == "html" {
		c.HTML(200, "inbox.html", gin.H{
			"webhooks": webhooks,
			"hash":     hash,
			"page":     page,
			"psize":    psize,
			"total":    total,
		})
		return
	}

	c.JSON(200, gin.H{
		"data":  webhooks,
		"page":  page,
		"psize": psize,
		"total": total,
	})
}

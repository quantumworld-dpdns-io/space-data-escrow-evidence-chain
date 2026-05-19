package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/domain"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/middleware"
	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/service"
)

func NewRouter(svc *service.Service, apiKey string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().UTC().Format(time.RFC3339)})
	})
	r.GET("/readyz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	v1 := r.Group("/v1")
	v1.Use(middleware.APIKey(apiKey))
	{
		v1.POST("/evidence", func(c *gin.Context) {
			var in service.CreateEvidenceInput
			if err := c.ShouldBindJSON(&in); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			rec, err := svc.CreateEvidence(in)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, rec)
		})

		v1.GET("/evidence/:id", func(c *gin.Context) {
			rec, ok := svc.GetEvidence(c.Param("id"))
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusOK, rec)
		})

		v1.POST("/custody", func(c *gin.Context) {
			var evt domain.CustodyEvent
			if err := c.ShouldBindJSON(&evt); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err := svc.AppendCustody(evt); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.Status(http.StatusAccepted)
		})

		v1.POST("/verify/:id", func(c *gin.Context) {
			report := svc.VerifyEvidence(c.Param("id"))
			if report.FailureReason == "not_found" {
				c.JSON(http.StatusNotFound, report)
				return
			}
			c.JSON(http.StatusOK, report)
		})

		v1.GET("/search", func(c *gin.Context) {
			c.JSON(http.StatusOK, svc.SearchEvidence(c.Query("q")))
		})

		v1.GET("/audit", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"entries": svc.AuditEntries()})
		})
	}

	return r
}

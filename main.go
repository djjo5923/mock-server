package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var database []LoadBalancer

type LoadBalancer struct {
	ID     string
	Name   string
	Status string
}

func initDatabase() {
	lb := LoadBalancer{
		ID:     "default-loadbalancer-id",
		Name:   "loadbalancer 01",
		Status: "Active",
	}

	database = append(database, lb)
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// List
	r.GET("/loadbalancers", func(c *gin.Context) {
		c.JSON(http.StatusOK, database)
	})

	// Get
	r.GET("/loadbalancers/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		loadbalancer, ok := getLoadBalancerByID(id)
		if ok {
			c.JSON(http.StatusOK, loadbalancer)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found - id"})
		}
	})

	// status
	r.GET("/loadbalancers/:id/status", func(c *gin.Context) {
		id := c.Params.ByName("id")
		loadbalancer, ok := getLoadBalancerByID(id)
		if ok {
			c.JSON(http.StatusOK, loadbalancer.Status)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found - id"})
		}
	})

	// create
	r.POST("/loadbalancers", func(c *gin.Context) {
		var loadbalancer LoadBalancer
		c.BindJSON(&loadbalancer)
		if loadbalancer.ID != "" && loadbalancer.Name != "" && loadbalancer.Status != "" {
			_, ok := getLoadBalancerByID(loadbalancer.ID)
			if !ok {
				database = append(database, loadbalancer)
				c.JSON(http.StatusCreated, loadbalancer)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "duplicate"})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "no value"})
		}
	})

	// delete
	r.DELETE("/loadbalancers/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		loadbalancer, ok := getLoadBalancerByID(id)
		if ok {
			ok := deleteLoadBalancerByID(id)
			if ok {
				c.JSON(http.StatusOK, loadbalancer)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "no value"})
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found - id"})
		}
	})

	return r
}

func deleteLoadBalancerByID(id string) bool {
	for k, v := range database {
		if v.ID == id {
			database = append(database[:k], database[k+1:]...)
			return true
		}
	}

	return false
}

func getLoadBalancerByID(id string) (LoadBalancer, bool) {
	for _, v := range database {
		if v.ID == id {
			return v, true
		}
	}

	return LoadBalancer{}, false
}

func main() {
	initDatabase()

	r := setupRouter()

	r.Run(":8080")
}

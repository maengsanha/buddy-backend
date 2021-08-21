// Package activity defines the router layer of the club activity of the Buddy System.
package activity

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/activity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create handles the activity creation request.
func Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(activity.Activity)
		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := activity.New(body.Start, body.End, body.Place, body.Description, body.Type, body.Participants, body.Private).
			Create(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Search handles the activity search request.
func Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			Query string `json:"query"`
		})
		resp := new(struct {
			Data struct {
				Activities []map[string]interface{} `json:"activities"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		activities, err := activity.Search(body.Query)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		resp.Data.Activities = activities.Public()
		c.JSON(http.StatusOK, resp)
	}
}

// Update handles the activity update request.
func Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			ID     string                 `json:"id"`
			Update map[string]interface{} `json:"update"`
		})
		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if objectID, err := primitive.ObjectIDFromHex(body.ID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else if err = (activity.Activity{ID: objectID}).Update(body.Update); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else {
			c.JSON(http.StatusOK, resp)
		}
	}
}

// Delete handles the activity deletion request.
func Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			ID string `json:"id"`
		})
		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if objectID, err := primitive.ObjectIDFromHex(body.ID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else if err = activity.Delete(objectID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else {
			c.JSON(http.StatusOK, resp)
		}
	}
}

// Upload handles the file upload request.
func Upload() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		id := c.Query("id")
		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		file, err := c.FormFile("file")
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		filename := filepath.Base(file.Filename)

		if err = c.SaveUploadedFile(file, activity.NewFile(filename).Absolute()); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		if objectID, err := primitive.ObjectIDFromHex(id); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else if err = (activity.Activity{ID: objectID}).Upload(filename); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else {
			c.JSON(http.StatusOK, resp)
		}
	}
}

// Download handles the file download request.
func Download() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			FileName string `json:"filename"`
		})
		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.File(activity.NewFile(body.FileName).Absolute())
	}
}

// DeleteFile handles the file deletion request.
func DeleteFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			ID       string `json:"id"`
			FileName string `json:"filename"`
		})
		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if objectID, err := primitive.ObjectIDFromHex(body.ID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else if err = (activity.Activity{ID: objectID}).DeleteFile(body.FileName); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else {
			c.JSON(http.StatusOK, resp)
		}
	}
}

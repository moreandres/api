// Copyright (c) 2021 Andres More

// Main

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	cfg "github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"

	"github.com/gin-gonic/gin"
	stats "github.com/semihalev/gin-stats"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// E set error message
func E(ctx *gin.Context, status int, title string) {
	ctx.JSON(status,
		gin.H{
			"errors": []map[string]interface{}{
				{
					"status": status,
					"title":  title,
				},
			}})
}

// HandleHealth handle health requests
func HandleHealth(ctx *gin.Context) {

	dbi, _ := DB.DB()

	ctx.JSON(http.StatusOK, // TODO: return OK only if DB can be reached
		gin.H{
			"data": map[string]interface{}{
				"type":        "health",
				"attributes":  stats.Report(),
				"attributes2": dbi.Stats(), // https://play.golang.org/p/8jlJUbEJKf
			},
		},
	)
}

// HandleGet handle resource GET
func HandleGet(ctx *gin.Context) {

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", cfg.GetString("QueryLimit")))
	if err != nil {
		E(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var objects []Object
	query := DB.Limit(limit)

	sort := ctx.Query("sort")
	if sort != "" {
		var isSort = regexp.MustCompile(`^[a-z-]+$`).MatchString
		if isSort(sort) {
			if sort[0] == '-' {
				query.Order(fmt.Sprint(sort[1:], " desc"))
			} else {
				query.Order(sort)
			}
		} else {
			E(ctx, http.StatusBadRequest, "unexpected sort field: "+sort)
			return
		}
	}

	result := query.Find(&objects)

	if result.Error != nil {
		E(ctx, http.StatusInternalServerError, result.Error.Error())
		return
	}

	log.Debugln(fmt.Sprintf("%+v", objects))

	data := []gin.H{}
	for _, object := range objects {
		item := gin.H{"ID": object.ID, "type": "objects", "attributes": object.Attributes}
		data = append(data, item)
	}

	log.Debugln(fmt.Sprintf("%+v", data))

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// HandleGetItem handle item GET
func HandleGetItem(ctx *gin.Context) {

	var object Object
	result := DB.Where("id = ?", ctx.Param("id")).First(&object)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			E(ctx, http.StatusNotFound, result.Error.Error())
		} else {
			E(ctx, http.StatusInternalServerError, result.Error.Error())
		}
		return
	}

	log.Debugln(fmt.Sprintf("%+v", object))

	ctx.JSON(http.StatusOK,
		gin.H{
			"data": map[string]interface{}{
				"type":       "objects",
				"ID":         object.ID,
				"attributes": object.Attributes,
			}})

}

// HandlePost handle POST
func HandlePost(ctx *gin.Context) {

	post := struct {
		Data Object
	}{}

	err := ctx.BindJSON(&post)
	if err != nil {
		E(ctx, http.StatusBadRequest, err.Error())
		return
	}

	c := context.Background()
	errs, err := Schema.ValidateBytes(c, post.Data.Attributes) // TODO: broken
	if err != nil {
		E(ctx, http.StatusBadRequest, errs[0].Error())
	}

	log.Debugln(fmt.Sprintf("%+v", post.Data.Attributes))

	object := &Object{Attributes: post.Data.Attributes}

	DB.Create(&object)

	log.Debugln(fmt.Sprintf("%+v", object))

	ctx.JSON(http.StatusCreated,
		gin.H{
			"data": map[string]interface{}{
				"type":       "objects",
				"ID":         object.ID,
				"attributes": object.Attributes,
			}})
}

// HandlePatch handle PATCH
func HandlePatch(context *gin.Context) {

	post := struct {
		Data Object
	}{}

	err := context.BindJSON(&post)
	if err != nil {
		E(context, http.StatusBadRequest, err.Error())
		return
	}

	var object Object
	result := DB.Where("id = ?", context.Param("id")).First(&object)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			E(context, http.StatusNotFound, "PATCH could not find first")
		} else {
			E(context, http.StatusInternalServerError, "TBD")
		}
		return
	}

	object.Attributes = post.Data.Attributes

	log.Debugln(fmt.Sprintf("%+v", object))

	DB.Model(&object).Updates(object)

	context.JSON(http.StatusOK,
		gin.H{
			"data": map[string]interface{}{
				"type":       "objects",
				"ID":         object.ID,
				"attributes": object.Attributes,
			}})
}

// HandleDelete handle DELETE
func HandleDelete(context *gin.Context) {
	var object Object

	result := DB.Where("id = ?", context.Param("id")).First(&object)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			E(context, http.StatusNotFound, result.Error.Error())
		} else {
			E(context, http.StatusInternalServerError, result.Error.Error())
		}
		return
	}

	log.Debugln(fmt.Sprintf("%+v", object))

	DB.Delete(&object)

	context.JSON(http.StatusNoContent, gin.H{})

}

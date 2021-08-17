// Package main runs the API server of the Buddy System.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/web/api/v1/activity"
	"github.com/kmu-kcc/buddy-backend/web/api/v1/fee"
	"github.com/kmu-kcc/buddy-backend/web/api/v1/member"
)

func main() {
	parser := argparse.NewParser("buddy", "API server of the Buddy System")

	// parse port number from command line arguments
	//
	// See https://github.com/akamensky/argparse#readme
	port := parser.Int("p", "port", &argparse.Options{Required: true, Help: "Port to run the server"})

	if err := parser.Parse(os.Args); err != nil {
		log.Fatalln(parser.Usage(err))
	}

	gin.SetMode(gin.DebugMode)

	engine := gin.Default()

	api := engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			mgroup := v1.Group("/member")
			{
				mgroup.POST("/signin", member.SignIn())
				mgroup.POST("/signup", member.SignUp())
				mgroup.GET("/signups", member.SignUps())
				mgroup.PUT("/approve", member.Approve())
				mgroup.DELETE("/delete", member.Delete())
				mgroup.PUT("/exit", member.Exit())
				mgroup.GET("/exits", member.Exits())
				mgroup.GET("/search", member.Search())
				mgroup.PUT("/update", member.Update())
				mgroup.PUT("/activate", nil)
				mgroup.PUT("/deactivate", nil)
				mgroup.GET("/graduates", member.Graduates())
			}
			agroup := v1.Group("/activity")
			{
				agroup.POST("/create", activity.Create())
				agroup.GET("/search", activity.Search())
				agroup.GET("/private", nil)
				agroup.PUT("/update", activity.Update())
				agroup.DELETE("/delete", activity.Delete())
				agroup.PUT("/addpicture", nil)
			}
			fgroup := v1.Group("/fee")
			{
				fgroup.POST("/create", fee.Create())
				fgroup.POST("/amount", fee.Amount())
				fgroup.POST("/payers", fee.Payers())
				fgroup.POST("/deptors", fee.Deptors())
				fgroup.POST("/search", fee.Search())
				fgroup.POST("/pay", fee.Pay())
				fgroup.POST("/deposit", fee.Deposit())
				fgroup.POST("/exempt", nil)
			}
		}
	}

	log.Fatalln(engine.Run(fmt.Sprintf(":%d", *port)))
}

package routes

import "github.com/gin-gonic/gin"

func WebRoutes(r *gin.Engine) {
	r.Static("/assets", "/home/ec2-user/My_Portfolio/frontend/assets")
	r.LoadHTMLGlob("/home/ec2-user/My_Portfolio/frontend/*.html")

	r.GET("/index", func(c *gin.Context) { c.HTML(200, "index.html", nil) })
	r.GET("/about", func(c *gin.Context) { c.HTML(200, "about.html", nil) })
	r.GET("/projects", func(c *gin.Context) { c.HTML(200, "projects.html", nil) })
	r.GET("/skills", func(c *gin.Context) { c.HTML(200, "skills.html", nil) })
	r.GET("/contact", func(c *gin.Context) { c.HTML(200, "contact.html", nil) })

	r.GET("/auth", func(c *gin.Context) { c.HTML(200, "auth.html", nil) })
	r.GET("/dashboard", func(c *gin.Context) { c.HTML(200, "dashboard.html", nil) })
}

package server

import (
	"embed"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kubebadges/kubebadges"
	"github.com/kubebadges/kubebadges/internal/server/controller"
	"github.com/kubebadges/kubebadges/internal/server/middleware"
)

var contentTypeMap = map[string]string{
	".js":   "application/javascript",
	".html": "text/html",
	".css":  "text/css",
	".json": "application/json",
	".png":  "image/png",
	".otf":  "font/otf",
	".ttf":  "font/ttf",
}

func getContentType(path string) (string, bool) {
	ext := filepath.Ext(path)
	contentType, ok := contentTypeMap[ext]
	return contentType, ok
}

func registerStaticFiles(router *gin.Engine, fs embed.FS, root string) {
	entries, _ := fs.ReadDir(root)
	for _, entry := range entries {
		if entry.IsDir() {
			registerStaticFiles(router, fs, root+"/"+entry.Name())
		} else {
			registerFile(router, fs, root, entry.Name())
		}
	}
}

func registerFile(router *gin.Engine, fs embed.FS, root, fileName string) {
	filePath := root + "/" + fileName
	mineType, mineTypeOK := getContentType(filePath)

	if filePath == "web/index.html" {
		router.GET("/", func(c *gin.Context) {
			serveFile(c, fs, filePath, mineType, mineTypeOK)
		})
	}

	router.GET("/"+strings.TrimPrefix(filePath, "web/"), func(c *gin.Context) {
		c.Header("Cache-Control", "private, max-age=0, no-cache")
		serveFile(c, fs, filePath, mineType, mineTypeOK)
	})
}

func serveFile(c *gin.Context, fs embed.FS, filePath, mineType string, mineTypeOK bool) {
	data, _ := fs.ReadFile(filePath)
	if mineTypeOK {
		c.Data(http.StatusOK, mineType, data)
	} else {
		c.Data(http.StatusOK, http.DetectContentType(data), data)
	}
}

func (s *Server) initRouter() {
	baseCtrl := &controller.BaseController{
		ServerContext: s.svcCtx,
	}
	kubeController := controller.NewKubeController(s.svcCtx)
	badgesController := controller.NewBadgesController(baseCtrl)

	registerStaticFiles(s.internalEngine, kubebadges.WebFiles, "web")

	s.internalEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Push-Id", "App", "App-Version", "X-Device-Id", "Content-Type", "Content-Length", "Authorization", "X-App-Name"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1")
		},
	}))

	// admin routes
	api := s.internalEngine.Group("/api")
	{
		api.GET("/nodes", kubeController.ListNodes)
		api.GET("/namespaces", kubeController.ListNamespaces)
		api.GET("/deployments/:namespace", kubeController.ListDeployments)
		api.POST("/badge", kubeController.UpdateBadge)
		api.GET("/config", kubeController.GetConfig)
		api.POST("/config", kubeController.UpdateConfig)

		// List Kustomizations (optional)
		api.GET("/kustomizations/:namespace", kubeController.ListKustomizations)
		api.GET("/postgresqls/:namespace", kubeController.ListPostgresqls)
		api.GET("/jobs/:namespace", kubeController.ListJobs)
	}

	badges := s.internalEngine.Group("/badges")
	{
		// badges routes
		badges.GET("/kube/node/:node", badgesController.Node)
		badges.GET("/kube/namespace/:namespace", badgesController.Namespace)
		badges.GET("/kube/deployment/:namespace/:deployment", badgesController.Deployment)
		badges.GET("/kube/pod/:namespace/:pod", badgesController.Pod)

		badges.GET("/kube/kustomization/:namespace/:kustomization", badgesController.Kustomization)
		badges.GET("/kube/postgresql/:namespace/:postgresql", badgesController.Postgresql)
		badges.GET("/kube/job/:namespace/:job", badgesController.Job)
	}

	// for external api
	s.externalEngine.NoRoute(func(ctx *gin.Context) {
		baseCtrl.NotFound(ctx)
	})
	s.externalEngine.Use(middleware.BadgeApiAccessMiddleware(s.svcCtx.KubeBadgesService))
	exBadges := s.externalEngine.Group("/badges")
	{
		exBadges.GET("/kube/node/:node", badgesController.Node)
		exBadges.GET("/kube/namespace/:namespace", badgesController.Namespace)
		exBadges.GET("/kube/deployment/:namespace/:deployment", badgesController.Deployment)
		exBadges.GET("/kube/pod/:namespace/:pod", badgesController.Pod)
		exBadges.GET("/kube/pod/:namespace/:pod/status", badgesController.Pod)

		exBadges.GET("/kube/kustomization/:namespace/:kustomization", badgesController.Kustomization)
		exBadges.GET("/kube/job/:namespace/:job", badgesController.Job)
		exBadges.GET("/kube/postgresql/:namespace/:postgresql", badgesController.Postgresql)
	}
}

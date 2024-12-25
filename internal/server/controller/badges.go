package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kubebadges/kubebadges/internal/badges"
	"github.com/kubebadges/kubebadges/internal/cache"
	"github.com/kubebadges/kubebadges/internal/server/svc"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// BadgeMessage structure
type BadgeMessage struct {
	Key          string
	Label        string
	Message      string
	MessageColor string
}

type BaseController struct {
	*svc.ServerContext
}

var notFoundSvg = `
<svg xmlns="http://www.w3.org/2000/svg" ...> ... </svg>
`

func (b *BaseController) NotFound(c *gin.Context) {
	resultType := c.DefaultQuery("type", "svg")
	if resultType == "json" {
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(404)
		return
	}

	c.Header("Content-Type", "image/svg+xml")
	c.String(404, notFoundSvg)
}

func (b *BaseController) Success(c *gin.Context, kubeBadgesService *svc.KubeBadgesService, badgeMessage BadgeMessage) {
	if kubeBadge, err := kubeBadgesService.GetKubeBadge(badgeMessage.Key, false); err == nil {
		if len(kubeBadge.Spec.DisplayName) > 0 {
			badgeMessage.Label = kubeBadge.Spec.DisplayName
		}
	}

	badge := badges.NewBadgeBuilder().
		SetLabel(badgeMessage.Label).
		SetMessage(badgeMessage.Message).
		SetMessageColor(badgeMessage.MessageColor).
		SetStyle(c.Query("style")).
		Build()

	resultType := c.DefaultQuery("type", "svg")
	if resultType == "json" {
		c.Header("Content-Type", "application/json")
		c.JSON(200, gin.H{
			"label":   badge.Label,
			"message": badge.Message,
			"color":   badge.MessageColor,
		})
		return
	}

	b.BadgesHelper.CreateBadgeProxy(badge, c)
}

// =============================================================
// BadgesController
// =============================================================
type BadgesController struct {
	BaseController
	namespaceCache  *cache.Cache[string, BadgeMessage]
	deploymentCache *cache.Cache[string, BadgeMessage]
	nodeCache       *cache.Cache[string, BadgeMessage]
	podCache        *cache.Cache[string, BadgeMessage]
	// ADDED for Kustomization
	kustomizationCache *cache.Cache[string, BadgeMessage]
}

func NewBadgesController(svc *svc.ServerContext) *BadgesController {
	return &BadgesController{
		BaseController: BaseController{
			ServerContext: svc,
		},
		namespaceCache:     cache.NewCache[string, BadgeMessage](),
		deploymentCache:    cache.NewCache[string, BadgeMessage](),
		nodeCache:          cache.NewCache[string, BadgeMessage](),
		podCache:           cache.NewCache[string, BadgeMessage](),
		kustomizationCache: cache.NewCache[string, BadgeMessage](), // ADDED
	}
}

func (s *BadgesController) getCacheDuration() time.Duration {
	return time.Duration(s.Config.CacheTime) * time.Second
}

// Node badge
func (s *BadgesController) Node(c *gin.Context) {
	name := c.Param("node")
	badgeMessage, ok := s.nodeCache.Get(name)
	if !ok {
		node, err := s.KubeHelper.GetNode(name)
		if err != nil {
			s.NotFound(c)
			return
		}

		badgeMessage = BadgeMessage{
			Key:   fmt.Sprintf("/kube/node/%s", name),
			Label: name,
		}

		isNodeReady := false
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
				isNodeReady = true
				badgeMessage.MessageColor = badges.Green
				badgeMessage.Message = string(condition.Type)
				break
			}
		}

		if isNodeReady {
			for _, condition := range node.Status.Conditions {
				if condition.Type != corev1.NodeReady && condition.Status == corev1.ConditionTrue {
					badgeMessage.MessageColor = badges.Yellow
					badgeMessage.Message = string(condition.Type)
					break
				}
			}
		} else {
			badgeMessage.MessageColor = badges.Red
			badgeMessage.Message = "NotReady"
		}

		s.nodeCache.Set(name, badgeMessage, s.getCacheDuration())
	}

	s.Success(c, s.KubeBadgesService, badgeMessage)
}

// Namespace badge
func (s *BadgesController) Namespace(c *gin.Context) {
	name := c.Param("namespace")
	badgeMessage, ok := s.namespaceCache.Get(name)
	if !ok {
		namespace, err := s.KubeHelper.GetNamespace(name)
		if err != nil {
			s.NotFound(c)
			return
		}

		badgeMessage = BadgeMessage{
			Key:     fmt.Sprintf("/kube/namespace/%s", name),
			Label:   name,
			Message: string(namespace.Status.Phase),
		}

		switch badgeMessage.Message {
		case string(corev1.NamespaceActive):
			badgeMessage.MessageColor = badges.Green
		case string(corev1.NamespaceTerminating):
			badgeMessage.MessageColor = badges.Red
		default:
			badgeMessage.MessageColor = badges.Blue
		}
		s.namespaceCache.Set(name, badgeMessage, s.getCacheDuration())
	}

	s.Success(c, s.KubeBadgesService, badgeMessage)
}

// Deployment badge
func (s *BadgesController) Deployment(c *gin.Context) {
	namespace := c.Param("namespace")
	deploymentName := c.Param("deployment")

	badgeMessage, ok := s.deploymentCache.Get(fmt.Sprintf("%s_%s", namespace, deploymentName))
	if !ok {
		deployment, err := s.KubeHelper.GetDeployment(namespace, deploymentName)
		if err != nil {
			s.NotFound(c)
			return
		}
		badgeMessage = BadgeMessage{
			Key:   fmt.Sprintf("/kube/deployment/%s/%s", namespace, deploymentName),
			Label: deploymentName,
		}
		statusMessage := ""
		available := true
		replicaFailure := false
		for _, condition := range deployment.Status.Conditions {
			if condition.Type == v1.DeploymentAvailable {
				available = condition.Status == corev1.ConditionTrue
			} else if condition.Type == v1.DeploymentReplicaFailure {
				replicaFailure = condition.Status == corev1.ConditionTrue
			}
		}

		if available && !replicaFailure {
			statusMessage = "Available"
		} else if available && replicaFailure {
			statusMessage = "Warning"
		} else if !available && !replicaFailure {
			statusMessage = "Unavailable"
		} else if !available && replicaFailure {
			statusMessage = "Failed"
		}

		switch statusMessage {
		case "Available":
			badgeMessage.MessageColor = badges.Green
		case "Warning":
			badgeMessage.MessageColor = badges.Yellow
		case "Unavailable":
			badgeMessage.MessageColor = badges.Red
		case "Failed":
			badgeMessage.MessageColor = badges.Red
		default:
			badgeMessage.MessageColor = badges.Blue
		}

		if deployment.Status.AvailableReplicas != deployment.Status.Replicas {
			badgeMessage.MessageColor = badges.Yellow
		}

		badgeMessage.Message = fmt.Sprintf("%d/%d %s", deployment.Status.AvailableReplicas, deployment.Status.Replicas, statusMessage)
		s.deploymentCache.Set(fmt.Sprintf("%s_%s", namespace, deploymentName), badgeMessage, s.getCacheDuration())
	}

	s.Success(c, s.KubeBadgesService, badgeMessage)
}

// Pod badge
func (s *BadgesController) Pod(c *gin.Context) {
	namespace := c.Param("namespace")
	podName := c.Param("pod")

	badgeMessage, ok := s.podCache.Get(fmt.Sprintf("%s_%s", namespace, podName))
	if !ok {
		pod, err := s.KubeHelper.GetPod(namespace, podName)
		if err != nil {
			s.NotFound(c)
			return
		}
		badgeMessage = BadgeMessage{
			Key:     fmt.Sprintf("/kube/pod/%s/%s", namespace, podName),
			Label:   podName,
			Message: string(pod.Status.Phase),
		}

		switch badgeMessage.Message {
		case string(corev1.PodRunning):
			badgeMessage.MessageColor = badges.Green
		case string(corev1.PodPending):
			badgeMessage.MessageColor = badges.Yellow
		case string(corev1.PodSucceeded):
			badgeMessage.MessageColor = badges.Green
		case string(corev1.PodFailed):
			badgeMessage.MessageColor = badges.Red
		case string(corev1.PodUnknown):
			badgeMessage.MessageColor = badges.Blue
		default:
			badgeMessage.MessageColor = badges.Blue
		}

		s.podCache.Set(fmt.Sprintf("%s_%s", namespace, podName), badgeMessage, s.getCacheDuration())
	}

	s.Success(c, s.KubeBadgesService, badgeMessage)
}

// ADDED for Kustomization
func (s *BadgesController) Kustomization(c *gin.Context) {
	namespace := c.Param("namespace")
	kustomizationName := c.Param("kustomization")

	key := fmt.Sprintf("/kube/kustomization/%s/%s", namespace, kustomizationName)
	badgeMessage, ok := s.kustomizationCache.Get(key)
	if !ok {
		kustomization, err := s.KubeHelper.GetKustomization(namespace, kustomizationName)
		if err != nil {
			s.NotFound(c)
			return
		}

		// On va parser le .status.conditions pour voir si c'est "Ready"
		label := kustomizationName
		message := "Unknown"
		messageColor := badges.Blue

		statusObj, hasStatus := kustomization["status"].(map[string]interface{})
		if hasStatus {
			if conditions, ok := statusObj["conditions"].([]interface{}); ok {
				for _, cnd := range conditions {
					cMap, ok := cnd.(map[string]interface{})
					if !ok {
						continue
					}
					cType, _ := cMap["type"].(string)
					cStatus, _ := cMap["status"].(string)
					if cType == "Ready" {
						if cStatus == "True" {
							message = "Ready"
							messageColor = badges.Green
						} else {
							message = "NotReady"
							messageColor = badges.Red
						}
						break
					}
				}
			}
		}

		badgeMessage = BadgeMessage{
			Key:          key,
			Label:        label,
			Message:      message,
			MessageColor: messageColor,
		}

		s.kustomizationCache.Set(key, badgeMessage, s.getCacheDuration())
	}

	s.Success(c, s.KubeBadgesService, badgeMessage)
}

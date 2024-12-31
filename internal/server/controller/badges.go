package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kubebadges/kubebadges/internal/badges"
	"github.com/kubebadges/kubebadges/internal/cache"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// =============================================================
// BadgesController
// =============================================================
type BadgesController struct {
	BaseController
	namespaceCache  *cache.Cache[string, BadgeMessage]
	deploymentCache *cache.Cache[string, BadgeMessage]
	nodeCache       *cache.Cache[string, BadgeMessage]
	podCache           *cache.Cache[string, BadgeMessage]
	kustomizationCache *cache.Cache[string, BadgeMessage]
	postgresqlCache    *cache.Cache[string, BadgeMessage]
	jobCache           *cache.Cache[string, BadgeMessage]
}

func NewBadgesController(base *BaseController) *BadgesController {
	return &BadgesController{
		BaseController: *base,
		namespaceCache:     cache.NewCache[string, BadgeMessage](),
		deploymentCache:    cache.NewCache[string, BadgeMessage](),
		nodeCache:          cache.NewCache[string, BadgeMessage](),
		podCache:           cache.NewCache[string, BadgeMessage](),
		kustomizationCache: cache.NewCache[string, BadgeMessage](),
		postgresqlCache:    cache.NewCache[string, BadgeMessage](),
		jobCache:           cache.NewCache[string, BadgeMessage](),
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

	s.Success(c, badgeMessage)
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

	s.Success(c, badgeMessage)
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

	s.Success(c, badgeMessage)
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

	s.Success(c, badgeMessage)
}

func (s *BadgesController) Job(c *gin.Context) {
	namespace := c.Param("namespace")
	jobName := c.Param("job")

	key := fmt.Sprintf("/kube/job/%s/%s", namespace, jobName)
	badgeMessage, ok := s.jobCache.Get(key)
	if !ok {
		job, err := s.KubeHelper.GetJob(namespace, jobName)
		if err != nil {
			s.NotFound(c)
			return
		}

		label := jobName
		message := "Unknown"
		messageColor := badges.Blue

		if job.Status.Succeeded > 0 {
			message = "Succeeded"
			messageColor = badges.Green
		} else if job.Status.Failed > 0 {
			message = "Failed"
			messageColor = badges.Red
		} else if job.Status.Active > 0 {
			message = "Active"
			messageColor = badges.Yellow
		}

		badgeMessage = BadgeMessage{
			Key:          key,
			Label:        label,
			Message:      message,
			MessageColor: messageColor,
		}

		s.jobCache.Set(key, badgeMessage, s.getCacheDuration())
	}

	s.Success(c, badgeMessage)
}

func (s *BadgesController) Postgresql(c *gin.Context) {
	namespace := c.Param("namespace")
	postgresqlName := c.Param("postgresql")

	key := fmt.Sprintf("/kube/postgresql/%s/%s", namespace, postgresqlName)
	badgeMessage, ok := s.postgresqlCache.Get(key)
	if !ok {
		postgresql, err := s.KubeHelper.GetPostgresql(namespace, postgresqlName)
		if err != nil {
			s.NotFound(c)
			return
		}

		label := postgresqlName
		message := "Unknown"
		messageColor := badges.Blue

		if status, ok := postgresql["status"].(map[string]interface{}); ok {
			if clusterStatus, exists := status["PostgresClusterStatus"].(string); exists {
				message = clusterStatus
				switch clusterStatus {
				case "Running":
					messageColor = badges.Green
				case "Creating":
					messageColor = badges.Yellow
				default:
					messageColor = badges.Red
				}
			}
		}

		badgeMessage = BadgeMessage{
			Key:          key,
			Label:        label,
			Message:      message,
			MessageColor: messageColor,
		}

		s.postgresqlCache.Set(key, badgeMessage, s.getCacheDuration())
	}

	s.Success(c, badgeMessage)
}

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

		// Parse .status.conditions to check if it's "Ready"
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

	s.Success(c, badgeMessage)
}

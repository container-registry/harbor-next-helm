package migrate

import (
	"fmt"
	"strconv"
	"strings"
)

// mapDirectFields maps simple fields that have the same meaning in both charts
func mapDirectFields(legacy *LegacyValues, result *MigrationResult) {
	result.Values.ExternalURL = legacy.ExternalURL
	result.Values.HarborAdminPassword = legacy.HarborAdminPassword
	result.Values.SecretKey = legacy.SecretKey
	result.Values.LogLevel = legacy.LogLevel

	if legacy.ImagePullPolicy != "" {
		result.Values.Image.PullPolicy = legacy.ImagePullPolicy
	}

	// Transform imagePullSecrets from [{name: x}] to [x]
	for _, secret := range legacy.ImagePullSecrets {
		if secret.Name != "" {
			result.Values.ImagePullSecrets = append(result.Values.ImagePullSecrets, secret.Name)
		}
	}
}

// mapIngress transforms expose configuration to ingress/gateway
func mapIngress(legacy *LegacyValues, result *MigrationResult) {
	// Check expose type
	exposeType := legacy.Expose.Type
	if exposeType == "" {
		exposeType = "ingress"
	}

	if exposeType != "ingress" {
		result.Warnings = append(result.Warnings, Warning{
			Level:   WarningLevelError,
			Path:    "expose.type",
			Message: fmt.Sprintf("expose.type=%q is not supported in harbor-next-helm. Only 'ingress' is supported.", exposeType),
			Action:  "Migrate to ingress or configure gateway API manually",
		})
	}

	// Enable ingress by default
	result.Values.Ingress.Enabled = true
	result.Values.Ingress.ClassName = legacy.Expose.Ingress.ClassName
	result.Values.Ingress.Annotations = legacy.Expose.Ingress.Annotations

	// Transform hosts: expose.ingress.hosts.core -> ingress.hosts[0].host
	if legacy.Expose.Ingress.Hosts.Core != "" {
		result.Values.Ingress.Hosts = []NewIngressHost{
			{
				Host: legacy.Expose.Ingress.Hosts.Core,
				Paths: []NewIngressPath{
					{Path: "/", PathType: "Prefix"},
				},
			},
		}
	}

	// Transform TLS
	if legacy.Expose.TLS.Enabled {
		switch legacy.Expose.TLS.CertSource {
		case "secret":
			if legacy.Expose.TLS.Secret.SecretName != "" {
				hostName := legacy.Expose.Ingress.Hosts.Core
				result.Values.Ingress.TLS = []NewIngressTLS{
					{
						SecretName: legacy.Expose.TLS.Secret.SecretName,
						Hosts:      []string{hostName},
					},
				}
			}
		case "auto":
			result.Warnings = append(result.Warnings, Warning{
				Level:   WarningLevelWarning,
				Path:    "expose.tls.certSource",
				Message: "certSource='auto' is not supported. Use cert-manager or provide TLS secret manually.",
				Action:  "Configure tls.certManager or provide ingress.tls with existing secret",
			})
		case "none", "":
			// No TLS configuration needed
		}
	}

	// Check for Gateway API route configuration
	if len(legacy.Expose.Route) > 0 {
		result.Warnings = append(result.Warnings, Warning{
			Level:   WarningLevelInfo,
			Path:    "expose.route",
			Message: "Route configuration detected. Use gateway.* for Gateway API in harbor-next-helm.",
			Action:  "Configure gateway.enabled, gateway.parentRefs, and gateway.hostnames",
		})
	}
}

// mapStorage transforms persistence.imageChartStorage to registry.storage
func mapStorage(legacy *LegacyValues, result *MigrationResult) {
	storage := &result.Values.Registry.Storage

	storageType := legacy.Persistence.ImageChartStorage.Type
	if storageType == "" {
		storageType = "filesystem"
	}
	storage.Type = storageType

	switch storageType {
	case "filesystem":
		rootDir := legacy.Persistence.ImageChartStorage.Filesystem.RootDirectory
		if rootDir == "" {
			rootDir = "/storage"
		}
		storage.Filesystem = map[string]interface{}{
			"rootdirectory": rootDir,
		}
	case "s3":
		storage.S3 = legacy.Persistence.ImageChartStorage.S3
	case "azure":
		storage.Azure = legacy.Persistence.ImageChartStorage.Azure
	case "gcs":
		storage.GCS = legacy.Persistence.ImageChartStorage.GCS
	case "oss":
		storage.OSS = legacy.Persistence.ImageChartStorage.OSS
	case "swift":
		// Warning already added in checkUnsupportedFeatures
	}

	// Map persistence PVC settings
	pvc := legacy.Persistence.PersistentVolumeClaim.Registry
	result.Values.Registry.Persistence.Enabled = legacy.Persistence.Enabled
	result.Values.Registry.Persistence.Size = pvc.Size
	if result.Values.Registry.Persistence.Size == "" {
		result.Values.Registry.Persistence.Size = "50Gi" // New default
	}
	result.Values.Registry.Persistence.StorageClass = pvc.StorageClass
	result.Values.Registry.Persistence.ExistingClaim = pvc.ExistingClaim
	result.Values.Registry.Persistence.Annotations = pvc.Annotations

	if pvc.AccessMode != "" {
		result.Values.Registry.Persistence.AccessModes = []string{pvc.AccessMode}
	}
}

// mapDatabase transforms database configuration
func mapDatabase(legacy *LegacyValues, result *MigrationResult) {
	dbType := legacy.Database.Type
	if dbType == "" {
		dbType = "internal"
	}

	if dbType == "internal" {
		result.Warnings = append(result.Warnings, Warning{
			Level:   WarningLevelError,
			Path:    "database.type",
			Message: "Internal database is not supported in harbor-next-helm. You must use an external PostgreSQL database.",
			Action:  "Deploy an external PostgreSQL instance (e.g., CloudNativePG, AWS RDS) and configure database.host",
		})
		return
	}

	ext := legacy.Database.External
	result.Values.Database.Host = ext.Host

	// Parse port string to int
	if ext.Port != "" {
		if port, err := strconv.Atoi(ext.Port); err == nil {
			result.Values.Database.Port = port
		}
	}
	if result.Values.Database.Port == 0 {
		result.Values.Database.Port = 5432
	}

	result.Values.Database.Username = ext.Username
	if result.Values.Database.Username == "" {
		result.Values.Database.Username = "postgres"
	}

	result.Values.Database.Password = ext.Password
	result.Values.Database.Database = ext.CoreDatabase
	if result.Values.Database.Database == "" {
		result.Values.Database.Database = "registry"
	}
	result.Values.Database.SSLMode = ext.SSLMode
	if result.Values.Database.SSLMode == "" {
		result.Values.Database.SSLMode = "disable"
	}
	result.Values.Database.ExistingSecret = ext.ExistingSecret

	// Connection pool settings
	if legacy.Database.MaxIdleConns > 0 {
		result.Values.Database.MaxIdleConns = legacy.Database.MaxIdleConns
	}
	if legacy.Database.MaxOpenConns > 0 {
		result.Values.Database.MaxOpenConns = legacy.Database.MaxOpenConns
	}
}

// mapRedis transforms redis configuration to valkey or externalRedis
func mapRedis(legacy *LegacyValues, result *MigrationResult) {
	redisType := legacy.Redis.Type
	if redisType == "" {
		redisType = "internal"
	}

	if redisType == "internal" {
		// Map internal redis to valkey subchart
		result.Values.Valkey.Enabled = true
		result.Values.Valkey.Architecture = "standalone"
		result.Values.Valkey.Auth.Enabled = true

		result.Warnings = append(result.Warnings, Warning{
			Level:   WarningLevelInfo,
			Path:    "redis.type",
			Message: "Internal Redis has been migrated to Valkey subchart (Redis-compatible)",
			Action:  "Review valkey.* settings in output",
		})
		return
	}

	// Map external redis
	result.Values.Valkey.Enabled = false
	ext := legacy.Redis.External

	// Parse addr "host:port" format
	if ext.Addr != "" {
		parts := strings.Split(ext.Addr, ":")
		result.Values.ExternalRedis.Host = parts[0]
		if len(parts) > 1 {
			if port, err := strconv.Atoi(parts[1]); err == nil {
				result.Values.ExternalRedis.Port = port
			}
		}
	}
	if result.Values.ExternalRedis.Port == 0 {
		result.Values.ExternalRedis.Port = 6379
	}

	result.Values.ExternalRedis.Password = ext.Password
	result.Values.ExternalRedis.SentinelMasterSet = ext.SentinelMasterSet
	result.Values.ExternalRedis.ExistingSecret = ext.ExistingSecret
}

// mapCore transforms core component configuration
func mapCore(legacy *LegacyValues, result *MigrationResult) {
	c := &result.Values.Core
	l := legacy.Core

	c.Image.Repository = l.Image.Repository
	c.Image.Tag = l.Image.Tag
	c.Replicas = l.Replicas
	c.Resources = l.Resources
	c.NodeSelector = l.NodeSelector
	c.Tolerations = l.Tolerations
	c.Affinity = l.Affinity
	c.PodAnnotations = l.PodAnnotations
	c.PodLabels = l.PodLabels

	// Rename extraEnvVars -> extraEnv
	c.ExtraEnv = l.ExtraEnvVars

	// Map secrets to secret map
	if l.Secret != "" || l.XSRFKey != "" {
		c.Secret = map[string]interface{}{}
		if l.Secret != "" {
			c.Secret["secret"] = l.Secret
		}
		if l.XSRFKey != "" {
			c.Secret["csrf_key"] = l.XSRFKey
		}
	}

	// Map core-specific settings to config (toEnvVars pattern)
	if l.QuotaUpdateProvider != "" {
		if c.Config == nil {
			c.Config = map[string]interface{}{}
		}
		c.Config["quota"] = map[string]interface{}{
			"update_provider": l.QuotaUpdateProvider,
		}
	}
}

// mapPortal transforms portal component configuration
func mapPortal(legacy *LegacyValues, result *MigrationResult) {
	p := &result.Values.Portal
	l := legacy.Portal

	p.Image.Repository = l.Image.Repository
	p.Image.Tag = l.Image.Tag
	p.Replicas = l.Replicas
	p.Resources = l.Resources
	p.NodeSelector = l.NodeSelector
	p.Tolerations = l.Tolerations
	p.Affinity = l.Affinity
	p.PodAnnotations = l.PodAnnotations
	p.PodLabels = l.PodLabels
	p.ExtraEnv = l.ExtraEnvVars
}

// mapRegistry transforms registry component configuration
func mapRegistry(legacy *LegacyValues, result *MigrationResult) {
	reg := &result.Values.Registry
	l := legacy.Registry

	// Flatten registry.registry.image -> registry.image
	reg.Image.Repository = l.Registry.Image.Repository
	reg.Image.Tag = l.Registry.Image.Tag

	// Controller stays nested
	reg.Controller.Image.Repository = l.Controller.Image.Repository
	reg.Controller.Image.Tag = l.Controller.Image.Tag

	reg.Replicas = l.Replicas
	reg.NodeSelector = l.NodeSelector
	reg.Tolerations = l.Tolerations
	reg.Affinity = l.Affinity
	reg.PodAnnotations = l.PodAnnotations
	reg.PodLabels = l.PodLabels

	// Merge extraEnvVars from both containers
	reg.ExtraEnv = append(l.Registry.ExtraEnvVars, l.Controller.ExtraEnvVars...)

	// Map registry secret if provided
	if l.Secret != "" {
		if reg.Secret == nil {
			reg.Secret = map[string]interface{}{}
		}
		reg.Secret["http_secret"] = l.Secret
	}
}

// mapJobservice transforms jobservice component configuration
func mapJobservice(legacy *LegacyValues, result *MigrationResult) {
	js := &result.Values.Jobservice
	l := legacy.Jobservice

	js.Image.Repository = l.Image.Repository
	js.Image.Tag = l.Image.Tag
	js.Replicas = l.Replicas
	js.Resources = l.Resources
	js.NodeSelector = l.NodeSelector
	js.Tolerations = l.Tolerations
	js.Affinity = l.Affinity
	js.PodAnnotations = l.PodAnnotations
	js.PodLabels = l.PodLabels
	js.ExtraEnv = l.ExtraEnvVars

	// Map jobservice-specific settings to config map (toEnvVars pattern)
	js.Config = map[string]interface{}{}
	if l.MaxJobWorkers > 0 {
		js.Config["max_job_workers"] = l.MaxJobWorkers
	}
	if len(l.JobLoggers) > 0 {
		js.Config["job_loggers"] = strings.Join(l.JobLoggers, ",")
	}
	if l.LoggerSweeperDuration > 0 {
		js.Config["logger_sweeper_duration"] = l.LoggerSweeperDuration
	}

	// Only include config if non-empty
	if len(js.Config) == 0 {
		js.Config = nil
	}

	if l.Secret != "" {
		js.Secret = map[string]interface{}{
			"secret": l.Secret,
		}
	}
}

// mapExporter transforms exporter component configuration
func mapExporter(legacy *LegacyValues, result *MigrationResult) {
	exp := &result.Values.Exporter
	l := legacy.Exporter

	// Exporter is always enabled in new chart if it has any config
	exp.Enabled = true
	exp.Image.Repository = l.Image.Repository
	exp.Image.Tag = l.Image.Tag
	exp.Replicas = l.Replicas
	exp.Resources = l.Resources
	exp.NodeSelector = l.NodeSelector
	exp.Tolerations = l.Tolerations
	exp.Affinity = l.Affinity
	exp.PodAnnotations = l.PodAnnotations
	exp.PodLabels = l.PodLabels

	// Map exporter-specific settings to config
	if l.CacheDuration > 0 || l.CacheCleanInterval > 0 {
		exp.Config = map[string]interface{}{}
		if l.CacheDuration > 0 {
			exp.Config["cache_time"] = l.CacheDuration
		}
		if l.CacheCleanInterval > 0 {
			exp.Config["cache_clean_interval"] = l.CacheCleanInterval
		}
	}
}

// mapTrivy transforms trivy configuration
func mapTrivy(legacy *LegacyValues, result *MigrationResult) {
	result.Values.Trivy.Enabled = legacy.Trivy.Enabled

	if legacy.Trivy.Enabled {
		// Check if there are detailed Trivy settings that can't be directly mapped
		hasDetailedConfig := legacy.Trivy.VulnType != "" ||
			legacy.Trivy.Severity != "" ||
			legacy.Trivy.IgnoreUnfixed ||
			legacy.Trivy.DebugMode ||
			legacy.Trivy.GitHubToken != "" ||
			legacy.Trivy.SkipUpdate ||
			legacy.Trivy.OfflineScan

		if hasDetailedConfig {
			result.Warnings = append(result.Warnings, Warning{
				Level:   WarningLevelWarning,
				Path:    "trivy",
				Message: "Trivy is now provided as a subchart (harbor-scanner-trivy). Detailed Trivy configuration must be done through the subchart values.",
				Action:  "Configure Trivy options under trivy.* using the harbor-scanner-trivy subchart format. See: https://github.com/aquasecurity/harbor-scanner-trivy",
			})
		}
	}
}

// mapMetrics transforms metrics configuration
func mapMetrics(legacy *LegacyValues, result *MigrationResult) {
	if !legacy.Metrics.Enabled {
		return
	}

	sm := legacy.Metrics.ServiceMonitor
	result.Values.Metrics.ServiceMonitor.Enabled = sm.Enabled
	result.Values.Metrics.ServiceMonitor.Labels = sm.AdditionalLabels
	result.Values.Metrics.ServiceMonitor.Interval = sm.Interval
}

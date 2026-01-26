package migrate

import (
	"fmt"
	"io"
	"sort"
)

// WarningLevel indicates the severity of a warning
type WarningLevel string

const (
	WarningLevelInfo    WarningLevel = "INFO"
	WarningLevelWarning WarningLevel = "WARNING"
	WarningLevelError   WarningLevel = "ERROR"
)

// Warning represents a migration warning
type Warning struct {
	Level   WarningLevel
	Path    string
	Message string
	Action  string
}

func (w Warning) String() string {
	return fmt.Sprintf("[%s] %s: %s\n  Action: %s", w.Level, w.Path, w.Message, w.Action)
}

// PrintWarnings outputs warnings to the given writer
func PrintWarnings(w io.Writer, warnings []Warning) {
	if len(warnings) == 0 {
		fmt.Fprintln(w, "No warnings generated. Migration completed successfully.")
		return
	}

	// Sort by level (ERROR first, then WARNING, then INFO)
	sort.Slice(warnings, func(i, j int) bool {
		levelOrder := map[WarningLevel]int{
			WarningLevelError:   0,
			WarningLevelWarning: 1,
			WarningLevelInfo:    2,
		}
		return levelOrder[warnings[i].Level] < levelOrder[warnings[j].Level]
	})

	errorCount := 0
	warningCount := 0
	infoCount := 0

	for _, warn := range warnings {
		switch warn.Level {
		case WarningLevelError:
			errorCount++
		case WarningLevelWarning:
			warningCount++
		case WarningLevelInfo:
			infoCount++
		}
	}

	fmt.Fprintf(w, "\n=== Migration Report ===\n")
	fmt.Fprintf(w, "Errors: %d, Warnings: %d, Info: %d\n\n", errorCount, warningCount, infoCount)

	for _, warn := range warnings {
		fmt.Fprintln(w, warn.String())
		fmt.Fprintln(w)
	}

	if errorCount > 0 {
		fmt.Fprintln(w, "!!! Migration has ERRORS that require manual intervention !!!")
	}
}

// HasErrors returns true if any warnings are at ERROR level
func HasErrors(warnings []Warning) bool {
	for _, w := range warnings {
		if w.Level == WarningLevelError {
			return true
		}
	}
	return false
}

// HasWarnings returns true if any warnings are at WARNING level or higher
func HasWarnings(warnings []Warning) bool {
	for _, w := range warnings {
		if w.Level == WarningLevelError || w.Level == WarningLevelWarning {
			return true
		}
	}
	return false
}

// checkUnsupportedFeatures examines legacy values for unsupported features
func checkUnsupportedFeatures(legacy *LegacyValues, warnings *[]Warning) {
	// Check for nginx configuration
	if legacy.Nginx.Image.Repository != "" || legacy.Nginx.Replicas > 0 {
		*warnings = append(*warnings, Warning{
			Level:   WarningLevelError,
			Path:    "nginx",
			Message: "Nginx component is not supported in harbor-next-helm. The new chart uses Kubernetes Ingress directly.",
			Action:  "Remove nginx configuration and configure ingress.* instead",
		})
	}

	// Check for notary
	if len(legacy.Notary) > 0 {
		if enabled, ok := legacy.Notary["enabled"].(bool); ok && enabled {
			*warnings = append(*warnings, Warning{
				Level:   WarningLevelError,
				Path:    "notary",
				Message: "Notary is deprecated and not supported in harbor-next-helm",
				Action:  "Migrate to Sigstore/cosign for container signing",
			})
		}
	}

	// Check for chartmuseum
	if len(legacy.Chartmuseum) > 0 {
		if enabled, ok := legacy.Chartmuseum["enabled"].(bool); ok && enabled {
			*warnings = append(*warnings, Warning{
				Level:   WarningLevelError,
				Path:    "chartmuseum",
				Message: "ChartMuseum is deprecated and not supported in harbor-next-helm",
				Action:  "Use Harbor's built-in OCI artifact support for Helm charts",
			})
		}
	}

	// Check for internalTLS
	if legacy.InternalTLS.Enabled {
		*warnings = append(*warnings, Warning{
			Level:   WarningLevelWarning,
			Path:    "internalTLS",
			Message: "Internal TLS between components is not yet implemented in harbor-next-helm",
			Action:  "Consider using a service mesh for internal encryption if required",
		})
	}

	// Check for proxy settings
	if len(legacy.Proxy) > 0 {
		*warnings = append(*warnings, Warning{
			Level:   WarningLevelWarning,
			Path:    "proxy",
			Message: "Proxy settings must be configured via component extraEnv",
			Action:  "Add HTTP_PROXY, HTTPS_PROXY, NO_PROXY to component extraEnv as needed",
		})
	}

	// Check for trace settings
	if len(legacy.Trace) > 0 {
		if enabled, ok := legacy.Trace["enabled"].(bool); ok && enabled {
			*warnings = append(*warnings, Warning{
				Level:   WarningLevelWarning,
				Path:    "trace",
				Message: "Tracing configuration must be done via core.config",
				Action:  "Configure tracing via core.config using the trace.* keys from core.yaml",
			})
		}
	}

	// Check for cache settings
	if len(legacy.Cache) > 0 {
		if enabled, ok := legacy.Cache["enabled"].(bool); ok && enabled {
			*warnings = append(*warnings, Warning{
				Level:   WarningLevelWarning,
				Path:    "cache",
				Message: "Cache configuration must be done via core.config",
				Action:  "Add cache.enabled and cache.expire_hours to core.config",
			})
		}
	}

	// Check for clusterIP/nodePort/loadBalancer expose types
	if len(legacy.Expose.ClusterIP) > 0 {
		*warnings = append(*warnings, Warning{
			Level:   WarningLevelError,
			Path:    "expose.clusterIP",
			Message: "ClusterIP expose type is not supported. Use ingress or Gateway API.",
			Action:  "Configure ingress.* or gateway.* instead",
		})
	}
	if len(legacy.Expose.NodePort) > 0 {
		*warnings = append(*warnings, Warning{
			Level:   WarningLevelError,
			Path:    "expose.nodePort",
			Message: "NodePort expose type is not supported. Use ingress or Gateway API.",
			Action:  "Configure ingress.* or gateway.* instead",
		})
	}
	if len(legacy.Expose.LoadBalancer) > 0 {
		*warnings = append(*warnings, Warning{
			Level:   WarningLevelError,
			Path:    "expose.loadBalancer",
			Message: "LoadBalancer expose type is not supported. Use ingress or Gateway API.",
			Action:  "Configure ingress.* or gateway.* instead",
		})
	}

	// Check for ipFamily settings
	if len(legacy.IPFamily) > 0 {
		*warnings = append(*warnings, Warning{
			Level:   WarningLevelWarning,
			Path:    "ipFamily",
			Message: "IP family configuration is not directly supported",
			Action:  "Configure dual-stack via Kubernetes cluster settings",
		})
	}

	// Check for Swift storage
	if len(legacy.Persistence.ImageChartStorage.Swift) > 0 {
		*warnings = append(*warnings, Warning{
			Level:   WarningLevelError,
			Path:    "persistence.imageChartStorage.swift",
			Message: "Swift storage backend is not supported in harbor-next-helm",
			Action:  "Migrate to a supported storage backend (s3, azure, gcs, oss, filesystem)",
		})
	}

	// Check for registry middleware (CloudFront)
	if len(legacy.Registry.Middleware) > 0 {
		if enabled, ok := legacy.Registry.Middleware["enabled"].(bool); ok && enabled {
			*warnings = append(*warnings, Warning{
				Level:   WarningLevelWarning,
				Path:    "registry.middleware",
				Message: "Registry middleware (CloudFront) configuration requires manual setup",
				Action:  "Configure middleware settings via registry.config or extraEnv",
			})
		}
	}

	// Check for registry credentials
	if legacy.Registry.Credentials.Username != "" || legacy.Registry.Credentials.Password != "" {
		*warnings = append(*warnings, Warning{
			Level:   WarningLevelInfo,
			Path:    "registry.credentials",
			Message: "Registry credentials handling has changed in harbor-next-helm",
			Action:  "Review registry authentication configuration - credentials are now auto-generated",
		})
	}
}

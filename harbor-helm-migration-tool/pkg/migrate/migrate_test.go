package migrate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrateDirectFields(t *testing.T) {
	legacy := &LegacyValues{
		ExternalURL:         "https://harbor.example.com",
		HarborAdminPassword: "admin123",
		SecretKey:           "secret-key-16chr",
		LogLevel:            "debug",
		ImagePullPolicy:     "Always",
		ImagePullSecrets: []LegacyImagePullSecret{
			{Name: "my-secret"},
			{Name: "another-secret"},
		},
	}

	result, err := Migrate(legacy)
	require.NoError(t, err)

	assert.Equal(t, "https://harbor.example.com", result.Values.ExternalURL)
	assert.Equal(t, "admin123", result.Values.HarborAdminPassword)
	assert.Equal(t, "secret-key-16chr", result.Values.SecretKey)
	assert.Equal(t, "debug", result.Values.LogLevel)
	assert.Equal(t, "Always", result.Values.Image.PullPolicy)
	assert.Equal(t, []string{"my-secret", "another-secret"}, result.Values.ImagePullSecrets)
}

func TestMigrateIngress(t *testing.T) {
	legacy := &LegacyValues{
		Expose: LegacyExpose{
			Type: "ingress",
			Ingress: LegacyExposeIngress{
				Hosts: LegacyIngressHosts{
					Core: "harbor.example.com",
				},
				ClassName: "nginx",
				Annotations: map[string]string{
					"cert-manager.io/cluster-issuer": "letsencrypt",
				},
			},
			TLS: LegacyExposeTLS{
				Enabled:    true,
				CertSource: "secret",
				Secret: LegacyTLSSecret{
					SecretName: "harbor-tls",
				},
			},
		},
	}

	result, err := Migrate(legacy)
	require.NoError(t, err)

	assert.True(t, result.Values.Ingress.Enabled)
	assert.Equal(t, "nginx", result.Values.Ingress.ClassName)
	assert.Equal(t, "letsencrypt", result.Values.Ingress.Annotations["cert-manager.io/cluster-issuer"])
	require.Len(t, result.Values.Ingress.Hosts, 1)
	assert.Equal(t, "harbor.example.com", result.Values.Ingress.Hosts[0].Host)
	require.Len(t, result.Values.Ingress.TLS, 1)
	assert.Equal(t, "harbor-tls", result.Values.Ingress.TLS[0].SecretName)
	assert.Contains(t, result.Values.Ingress.TLS[0].Hosts, "harbor.example.com")
}

func TestMigrateIngressNonIngressType(t *testing.T) {
	legacy := &LegacyValues{
		Expose: LegacyExpose{
			Type: "nodePort",
		},
	}

	result, err := Migrate(legacy)
	require.NoError(t, err)

	assert.True(t, HasErrors(result.Warnings))
	found := false
	for _, w := range result.Warnings {
		if w.Path == "expose.type" && w.Level == WarningLevelError {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected error warning for expose.type")
}

func TestMigrateExternalDatabase(t *testing.T) {
	legacy := &LegacyValues{
		Database: LegacyDatabase{
			Type: "external",
			External: LegacyExternalDB{
				Host:         "postgres.example.com",
				Port:         "5432",
				Username:     "harbor",
				Password:     "dbpassword",
				CoreDatabase: "harbor_db",
				SSLMode:      "require",
			},
			MaxIdleConns: 50,
			MaxOpenConns: 100,
		},
	}

	result, err := Migrate(legacy)
	require.NoError(t, err)

	assert.Equal(t, "postgres.example.com", result.Values.Database.Host)
	assert.Equal(t, 5432, result.Values.Database.Port)
	assert.Equal(t, "harbor", result.Values.Database.Username)
	assert.Equal(t, "dbpassword", result.Values.Database.Password)
	assert.Equal(t, "harbor_db", result.Values.Database.Database)
	assert.Equal(t, "require", result.Values.Database.SSLMode)
	assert.Equal(t, 50, result.Values.Database.MaxIdleConns)
	assert.Equal(t, 100, result.Values.Database.MaxOpenConns)
}

func TestMigrateInternalDatabase(t *testing.T) {
	legacy := &LegacyValues{
		Database: LegacyDatabase{
			Type: "internal",
		},
	}

	result, err := Migrate(legacy)
	require.NoError(t, err)

	assert.True(t, HasErrors(result.Warnings))
	found := false
	for _, w := range result.Warnings {
		if w.Path == "database.type" && w.Level == WarningLevelError {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected error warning for internal database")
}

func TestMigrateInternalRedis(t *testing.T) {
	legacy := &LegacyValues{
		Redis: LegacyRedis{
			Type: "internal",
		},
	}

	result, err := Migrate(legacy)
	require.NoError(t, err)

	assert.True(t, result.Values.Valkey.Enabled)
	assert.Equal(t, "standalone", result.Values.Valkey.Architecture)
	assert.True(t, result.Values.Valkey.Auth.Enabled)
}

func TestMigrateExternalRedis(t *testing.T) {
	legacy := &LegacyValues{
		Redis: LegacyRedis{
			Type: "external",
			External: LegacyExternalRedis{
				Addr:              "redis.example.com:6379",
				Password:          "redispassword",
				SentinelMasterSet: "mymaster",
			},
		},
	}

	result, err := Migrate(legacy)
	require.NoError(t, err)

	assert.False(t, result.Values.Valkey.Enabled)
	assert.Equal(t, "redis.example.com", result.Values.ExternalRedis.Host)
	assert.Equal(t, 6379, result.Values.ExternalRedis.Port)
	assert.Equal(t, "redispassword", result.Values.ExternalRedis.Password)
	assert.Equal(t, "mymaster", result.Values.ExternalRedis.SentinelMasterSet)
}

func TestMigrateStorage(t *testing.T) {
	t.Run("S3 storage", func(t *testing.T) {
		legacy := &LegacyValues{
			Persistence: LegacyPersistence{
				Enabled: true,
				ImageChartStorage: LegacyImageChartStorage{
					Type: "s3",
					S3: map[string]interface{}{
						"bucket":         "harbor-bucket",
						"region":         "us-east-1",
						"accesskey":      "AKIAIOSFODNN7EXAMPLE",
						"secretkey":      "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
						"regionendpoint": "https://s3.us-east-1.amazonaws.com",
					},
				},
				PersistentVolumeClaim: LegacyPVCConfig{
					Registry: LegacyPVC{
						Size:         "100Gi",
						StorageClass: "gp2",
					},
				},
			},
		}

		result, err := Migrate(legacy)
		require.NoError(t, err)

		assert.Equal(t, "s3", result.Values.Registry.Storage.Type)
		assert.Equal(t, "harbor-bucket", result.Values.Registry.Storage.S3["bucket"])
		assert.True(t, result.Values.Registry.Persistence.Enabled)
		assert.Equal(t, "100Gi", result.Values.Registry.Persistence.Size)
		assert.Equal(t, "gp2", result.Values.Registry.Persistence.StorageClass)
	})

	t.Run("filesystem storage", func(t *testing.T) {
		legacy := &LegacyValues{
			Persistence: LegacyPersistence{
				Enabled: true,
				ImageChartStorage: LegacyImageChartStorage{
					Type: "filesystem",
					Filesystem: LegacyFilesystemStorage{
						RootDirectory: "/data/registry",
					},
				},
			},
		}

		result, err := Migrate(legacy)
		require.NoError(t, err)

		assert.Equal(t, "filesystem", result.Values.Registry.Storage.Type)
		assert.Equal(t, "/data/registry", result.Values.Registry.Storage.Filesystem["rootdirectory"])
	})
}

func TestMigrateCore(t *testing.T) {
	legacy := &LegacyValues{
		Core: LegacyCore{
			Image: LegacyImage{
				Repository: "goharbor/harbor-core",
				Tag:        "v2.10.0",
			},
			Replicas: 3,
			Secret:   "core-secret",
			XSRFKey:  "xsrf-key",
			Resources: map[string]interface{}{
				"limits": map[string]interface{}{
					"cpu":    "1",
					"memory": "1Gi",
				},
			},
			NodeSelector: map[string]string{
				"node-type": "harbor",
			},
			ExtraEnvVars: []interface{}{
				map[string]interface{}{
					"name":  "CUSTOM_VAR",
					"value": "custom-value",
				},
			},
			QuotaUpdateProvider: "db",
		},
	}

	result, err := Migrate(legacy)
	require.NoError(t, err)

	assert.Equal(t, "goharbor/harbor-core", result.Values.Core.Image.Repository)
	assert.Equal(t, "v2.10.0", result.Values.Core.Image.Tag)
	assert.Equal(t, 3, result.Values.Core.Replicas)
	assert.Equal(t, "core-secret", result.Values.Core.Secret["secret"])
	assert.Equal(t, "xsrf-key", result.Values.Core.Secret["csrf_key"])
	assert.Equal(t, "harbor", result.Values.Core.NodeSelector["node-type"])
	assert.NotEmpty(t, result.Values.Core.ExtraEnv)
	assert.Equal(t, "db", result.Values.Core.Config["quota"].(map[string]interface{})["update_provider"])
}

func TestMigrateJobservice(t *testing.T) {
	legacy := &LegacyValues{
		Jobservice: LegacyJobservice{
			Image: LegacyImage{
				Repository: "goharbor/harbor-jobservice",
				Tag:        "v2.10.0",
			},
			Replicas:              2,
			MaxJobWorkers:         10,
			JobLoggers:            []string{"file", "database"},
			LoggerSweeperDuration: 14,
			Secret:                "jobservice-secret",
		},
	}

	result, err := Migrate(legacy)
	require.NoError(t, err)

	assert.Equal(t, "goharbor/harbor-jobservice", result.Values.Jobservice.Image.Repository)
	assert.Equal(t, 2, result.Values.Jobservice.Replicas)
	assert.Equal(t, 10, result.Values.Jobservice.Config["max_job_workers"])
	assert.Equal(t, "file,database", result.Values.Jobservice.Config["job_loggers"])
	assert.Equal(t, 14, result.Values.Jobservice.Config["logger_sweeper_duration"])
	assert.Equal(t, "jobservice-secret", result.Values.Jobservice.Secret["secret"])
}

func TestMigrateRegistry(t *testing.T) {
	legacy := &LegacyValues{
		Registry: LegacyRegistry{
			Registry: LegacyRegistryContainer{
				Image: LegacyImage{
					Repository: "goharbor/registry-photon",
					Tag:        "v2.10.0",
				},
				ExtraEnvVars: []interface{}{
					map[string]interface{}{"name": "REG_VAR", "value": "reg-value"},
				},
			},
			Controller: LegacyRegistryContainer{
				Image: LegacyImage{
					Repository: "goharbor/harbor-registryctl",
					Tag:        "v2.10.0",
				},
				ExtraEnvVars: []interface{}{
					map[string]interface{}{"name": "CTL_VAR", "value": "ctl-value"},
				},
			},
			Replicas: 2,
			Secret:   "registry-secret",
		},
	}

	result, err := Migrate(legacy)
	require.NoError(t, err)

	assert.Equal(t, "goharbor/registry-photon", result.Values.Registry.Image.Repository)
	assert.Equal(t, "v2.10.0", result.Values.Registry.Image.Tag)
	assert.Equal(t, "goharbor/harbor-registryctl", result.Values.Registry.Controller.Image.Repository)
	assert.Equal(t, 2, result.Values.Registry.Replicas)
	assert.Equal(t, "registry-secret", result.Values.Registry.Secret["http_secret"])
	assert.Len(t, result.Values.Registry.ExtraEnv, 2)
}

func TestMigrateTrivy(t *testing.T) {
	t.Run("Trivy enabled with detailed config", func(t *testing.T) {
		legacy := &LegacyValues{
			Trivy: LegacyTrivy{
				Enabled:       true,
				VulnType:      "os,library",
				Severity:      "CRITICAL,HIGH",
				IgnoreUnfixed: true,
			},
		}

		result, err := Migrate(legacy)
		require.NoError(t, err)

		assert.True(t, result.Values.Trivy.Enabled)
		// Should have a warning about Trivy subchart format
		found := false
		for _, w := range result.Warnings {
			if w.Path == "trivy" && w.Level == WarningLevelWarning {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected warning for Trivy subchart format")
	})

	t.Run("Trivy disabled", func(t *testing.T) {
		legacy := &LegacyValues{
			Trivy: LegacyTrivy{
				Enabled: false,
			},
		}

		result, err := Migrate(legacy)
		require.NoError(t, err)

		assert.False(t, result.Values.Trivy.Enabled)
	})
}

func TestMigrateMetrics(t *testing.T) {
	legacy := &LegacyValues{
		Metrics: LegacyMetrics{
			Enabled: true,
			ServiceMonitor: LegacyServiceMonitor{
				Enabled:  true,
				Interval: "30s",
				AdditionalLabels: map[string]string{
					"release": "prometheus",
				},
			},
		},
	}

	result, err := Migrate(legacy)
	require.NoError(t, err)

	assert.True(t, result.Values.Metrics.ServiceMonitor.Enabled)
	assert.Equal(t, "30s", result.Values.Metrics.ServiceMonitor.Interval)
	assert.Equal(t, "prometheus", result.Values.Metrics.ServiceMonitor.Labels["release"])
}

func TestCheckUnsupportedFeatures(t *testing.T) {
	t.Run("nginx configuration", func(t *testing.T) {
		legacy := &LegacyValues{
			Nginx: LegacyNginx{
				Image: LegacyImage{
					Repository: "goharbor/nginx-photon",
				},
			},
		}

		result, err := Migrate(legacy)
		require.NoError(t, err)

		found := false
		for _, w := range result.Warnings {
			if w.Path == "nginx" && w.Level == WarningLevelError {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected error for nginx configuration")
	})

	t.Run("notary enabled", func(t *testing.T) {
		legacy := &LegacyValues{
			Notary: map[string]interface{}{
				"enabled": true,
			},
		}

		result, err := Migrate(legacy)
		require.NoError(t, err)

		found := false
		for _, w := range result.Warnings {
			if w.Path == "notary" && w.Level == WarningLevelError {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected error for notary")
	})

	t.Run("chartmuseum enabled", func(t *testing.T) {
		legacy := &LegacyValues{
			Chartmuseum: map[string]interface{}{
				"enabled": true,
			},
		}

		result, err := Migrate(legacy)
		require.NoError(t, err)

		found := false
		for _, w := range result.Warnings {
			if w.Path == "chartmuseum" && w.Level == WarningLevelError {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected error for chartmuseum")
	})

	t.Run("internalTLS enabled", func(t *testing.T) {
		legacy := &LegacyValues{
			InternalTLS: LegacyInternalTLS{
				Enabled: true,
			},
		}

		result, err := Migrate(legacy)
		require.NoError(t, err)

		found := false
		for _, w := range result.Warnings {
			if w.Path == "internalTLS" && w.Level == WarningLevelWarning {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected warning for internalTLS")
	})

	t.Run("swift storage", func(t *testing.T) {
		legacy := &LegacyValues{
			Persistence: LegacyPersistence{
				ImageChartStorage: LegacyImageChartStorage{
					Swift: map[string]interface{}{
						"container": "harbor",
					},
				},
			},
		}

		result, err := Migrate(legacy)
		require.NoError(t, err)

		found := false
		for _, w := range result.Warnings {
			if w.Path == "persistence.imageChartStorage.swift" && w.Level == WarningLevelError {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected error for swift storage")
	})
}

func TestMigrateFromYAML(t *testing.T) {
	yamlData := `
externalURL: https://harbor.example.com
harborAdminPassword: Harbor12345
expose:
  type: ingress
  ingress:
    hosts:
      core: harbor.example.com
    className: nginx
  tls:
    enabled: true
    certSource: secret
    secret:
      secretName: harbor-tls
database:
  type: external
  external:
    host: postgres.example.com
    port: "5432"
    username: harbor
    password: dbpassword
    coreDatabase: harbor
redis:
  type: internal
`

	result, err := MigrateFromYAML([]byte(yamlData))
	require.NoError(t, err)

	assert.Equal(t, "https://harbor.example.com", result.Values.ExternalURL)
	assert.Equal(t, "Harbor12345", result.Values.HarborAdminPassword)
	assert.True(t, result.Values.Ingress.Enabled)
	assert.Equal(t, "harbor.example.com", result.Values.Ingress.Hosts[0].Host)
	assert.Equal(t, "postgres.example.com", result.Values.Database.Host)
	assert.True(t, result.Values.Valkey.Enabled)
}

func TestToYAML(t *testing.T) {
	result := &MigrationResult{
		Values: &NewValues{
			ExternalURL:         "https://harbor.example.com",
			HarborAdminPassword: "admin123",
		},
	}

	output, err := result.ToYAML()
	require.NoError(t, err)

	assert.Contains(t, string(output), "externalURL: https://harbor.example.com")
	assert.Contains(t, string(output), "harborAdminPassword: admin123")
}

func TestToYAMLWithHeader(t *testing.T) {
	result := &MigrationResult{
		Values: &NewValues{
			ExternalURL: "https://harbor.example.com",
		},
	}

	output, err := result.ToYAMLWithHeader("/path/to/legacy.yaml")
	require.NoError(t, err)

	assert.Contains(t, string(output), "MIGRATED VALUES")
	assert.Contains(t, string(output), "Original file: /path/to/legacy.yaml")
}

func TestWarningHelpers(t *testing.T) {
	t.Run("HasErrors", func(t *testing.T) {
		warnings := []Warning{
			{Level: WarningLevelInfo, Path: "test", Message: "info"},
			{Level: WarningLevelWarning, Path: "test", Message: "warning"},
		}
		assert.False(t, HasErrors(warnings))

		warnings = append(warnings, Warning{Level: WarningLevelError, Path: "test", Message: "error"})
		assert.True(t, HasErrors(warnings))
	})

	t.Run("HasWarnings", func(t *testing.T) {
		warnings := []Warning{
			{Level: WarningLevelInfo, Path: "test", Message: "info"},
		}
		assert.False(t, HasWarnings(warnings))

		warnings = append(warnings, Warning{Level: WarningLevelWarning, Path: "test", Message: "warning"})
		assert.True(t, HasWarnings(warnings))
	})
}

// Helper function
func getKey(m map[string]string, key string) string {
	if m == nil {
		return ""
	}
	return m[key]
}

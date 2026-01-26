package migrate

// LegacyValues represents the structure of harbor-helm values.yaml
type LegacyValues struct {
	Expose              LegacyExpose            `yaml:"expose,omitempty"`
	ExternalURL         string                  `yaml:"externalURL,omitempty"`
	Persistence         LegacyPersistence       `yaml:"persistence,omitempty"`
	HarborAdminPassword string                  `yaml:"harborAdminPassword,omitempty"`
	SecretKey           string                  `yaml:"secretKey,omitempty"`
	LogLevel            string                  `yaml:"logLevel,omitempty"`
	ImagePullPolicy     string                  `yaml:"imagePullPolicy,omitempty"`
	ImagePullSecrets    []LegacyImagePullSecret `yaml:"imagePullSecrets,omitempty"`

	// Components
	Core       LegacyCore       `yaml:"core,omitempty"`
	Portal     LegacyPortal     `yaml:"portal,omitempty"`
	Registry   LegacyRegistry   `yaml:"registry,omitempty"`
	Jobservice LegacyJobservice `yaml:"jobservice,omitempty"`
	Trivy      LegacyTrivy      `yaml:"trivy,omitempty"`
	Exporter   LegacyExporter   `yaml:"exporter,omitempty"`
	Nginx      LegacyNginx      `yaml:"nginx,omitempty"`

	// Dependencies
	Database LegacyDatabase `yaml:"database,omitempty"`
	Redis    LegacyRedis    `yaml:"redis,omitempty"`

	// Deprecated/Unsupported
	Notary      map[string]interface{} `yaml:"notary,omitempty"`
	Chartmuseum map[string]interface{} `yaml:"chartmuseum,omitempty"`

	// Other
	InternalTLS LegacyInternalTLS      `yaml:"internalTLS,omitempty"`
	Metrics     LegacyMetrics          `yaml:"metrics,omitempty"`
	Proxy       map[string]interface{} `yaml:"proxy,omitempty"`
	Trace       map[string]interface{} `yaml:"trace,omitempty"`
	Cache       map[string]interface{} `yaml:"cache,omitempty"`
	IPFamily    map[string]interface{} `yaml:"ipFamily,omitempty"`
}

type LegacyImagePullSecret struct {
	Name string `yaml:"name,omitempty"`
}

// Expose configuration
type LegacyExpose struct {
	Type         string                 `yaml:"type,omitempty"`
	TLS          LegacyExposeTLS        `yaml:"tls,omitempty"`
	Ingress      LegacyExposeIngress    `yaml:"ingress,omitempty"`
	Route        map[string]interface{} `yaml:"route,omitempty"`
	ClusterIP    map[string]interface{} `yaml:"clusterIP,omitempty"`
	NodePort     map[string]interface{} `yaml:"nodePort,omitempty"`
	LoadBalancer map[string]interface{} `yaml:"loadBalancer,omitempty"`
}

type LegacyExposeTLS struct {
	Enabled    bool              `yaml:"enabled,omitempty"`
	CertSource string            `yaml:"certSource,omitempty"`
	Auto       map[string]string `yaml:"auto,omitempty"`
	Secret     LegacyTLSSecret   `yaml:"secret,omitempty"`
}

type LegacyTLSSecret struct {
	SecretName string `yaml:"secretName,omitempty"`
}

type LegacyExposeIngress struct {
	Hosts       LegacyIngressHosts `yaml:"hosts,omitempty"`
	Controller  string             `yaml:"controller,omitempty"`
	ClassName   string             `yaml:"className,omitempty"`
	Annotations map[string]string  `yaml:"annotations,omitempty"`
	Labels      map[string]string  `yaml:"labels,omitempty"`
}

type LegacyIngressHosts struct {
	Core string `yaml:"core,omitempty"`
}

// Persistence configuration
type LegacyPersistence struct {
	Enabled              bool                    `yaml:"enabled,omitempty"`
	ResourcePolicy       string                  `yaml:"resourcePolicy,omitempty"`
	PersistentVolumeClaim LegacyPVCConfig        `yaml:"persistentVolumeClaim,omitempty"`
	ImageChartStorage    LegacyImageChartStorage `yaml:"imageChartStorage,omitempty"`
}

type LegacyPVCConfig struct {
	Registry   LegacyPVC              `yaml:"registry,omitempty"`
	Jobservice map[string]interface{} `yaml:"jobservice,omitempty"`
	Database   map[string]interface{} `yaml:"database,omitempty"`
	Redis      map[string]interface{} `yaml:"redis,omitempty"`
	Trivy      map[string]interface{} `yaml:"trivy,omitempty"`
}

type LegacyPVC struct {
	ExistingClaim string            `yaml:"existingClaim,omitempty"`
	StorageClass  string            `yaml:"storageClass,omitempty"`
	SubPath       string            `yaml:"subPath,omitempty"`
	AccessMode    string            `yaml:"accessMode,omitempty"`
	Size          string            `yaml:"size,omitempty"`
	Annotations   map[string]string `yaml:"annotations,omitempty"`
}

type LegacyImageChartStorage struct {
	Type            string                  `yaml:"type,omitempty"`
	DisableRedirect bool                    `yaml:"disableredirect,omitempty"`
	Filesystem      LegacyFilesystemStorage `yaml:"filesystem,omitempty"`
	S3              map[string]interface{}  `yaml:"s3,omitempty"`
	Azure           map[string]interface{}  `yaml:"azure,omitempty"`
	GCS             map[string]interface{}  `yaml:"gcs,omitempty"`
	OSS             map[string]interface{}  `yaml:"oss,omitempty"`
	Swift           map[string]interface{}  `yaml:"swift,omitempty"`
}

type LegacyFilesystemStorage struct {
	RootDirectory string `yaml:"rootdirectory,omitempty"`
	MaxThreads    int    `yaml:"maxthreads,omitempty"`
}

// Database configuration
type LegacyDatabase struct {
	Type         string                 `yaml:"type,omitempty"`
	Internal     map[string]interface{} `yaml:"internal,omitempty"`
	External     LegacyExternalDB       `yaml:"external,omitempty"`
	MaxIdleConns int                    `yaml:"maxIdleConns,omitempty"`
	MaxOpenConns int                    `yaml:"maxOpenConns,omitempty"`
	PodAnnotations map[string]string    `yaml:"podAnnotations,omitempty"`
	PodLabels      map[string]string    `yaml:"podLabels,omitempty"`
}

type LegacyExternalDB struct {
	Host           string `yaml:"host,omitempty"`
	Port           string `yaml:"port,omitempty"`
	Username       string `yaml:"username,omitempty"`
	Password       string `yaml:"password,omitempty"`
	CoreDatabase   string `yaml:"coreDatabase,omitempty"`
	SSLMode        string `yaml:"sslmode,omitempty"`
	ExistingSecret string `yaml:"existingSecret,omitempty"`
}

// Redis configuration
type LegacyRedis struct {
	Type           string                 `yaml:"type,omitempty"`
	Internal       map[string]interface{} `yaml:"internal,omitempty"`
	External       LegacyExternalRedis    `yaml:"external,omitempty"`
	PodAnnotations map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels      map[string]string      `yaml:"podLabels,omitempty"`
}

type LegacyExternalRedis struct {
	Addr              string `yaml:"addr,omitempty"`
	SentinelMasterSet string `yaml:"sentinelMasterSet,omitempty"`
	Username          string `yaml:"username,omitempty"`
	Password          string `yaml:"password,omitempty"`
	ExistingSecret    string `yaml:"existingSecret,omitempty"`
}

// Component configurations
type LegacyImage struct {
	Repository string `yaml:"repository,omitempty"`
	Tag        string `yaml:"tag,omitempty"`
}

type LegacyCore struct {
	Image          LegacyImage            `yaml:"image,omitempty"`
	Replicas       int                    `yaml:"replicas,omitempty"`
	Resources      map[string]interface{} `yaml:"resources,omitempty"`
	NodeSelector   map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations    []interface{}          `yaml:"tolerations,omitempty"`
	Affinity       map[string]interface{} `yaml:"affinity,omitempty"`
	PodAnnotations map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels      map[string]string      `yaml:"podLabels,omitempty"`
	ExtraEnvVars   []interface{}          `yaml:"extraEnvVars,omitempty"`
	Secret         string                 `yaml:"secret,omitempty"`
	XSRFKey        string                 `yaml:"xsrfKey,omitempty"`
	ServiceAccountName string             `yaml:"serviceAccountName,omitempty"`
	StartupProbe   map[string]interface{} `yaml:"startupProbe,omitempty"`
	ConfigureUserSettings string          `yaml:"configureUserSettings,omitempty"`
	QuotaUpdateProvider string            `yaml:"quotaUpdateProvider,omitempty"`
}

type LegacyPortal struct {
	Image          LegacyImage            `yaml:"image,omitempty"`
	Replicas       int                    `yaml:"replicas,omitempty"`
	Resources      map[string]interface{} `yaml:"resources,omitempty"`
	NodeSelector   map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations    []interface{}          `yaml:"tolerations,omitempty"`
	Affinity       map[string]interface{} `yaml:"affinity,omitempty"`
	PodAnnotations map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels      map[string]string      `yaml:"podLabels,omitempty"`
	ExtraEnvVars   []interface{}          `yaml:"extraEnvVars,omitempty"`
	ServiceAccountName string             `yaml:"serviceAccountName,omitempty"`
}

type LegacyRegistry struct {
	Registry       LegacyRegistryContainer `yaml:"registry,omitempty"`
	Controller     LegacyRegistryContainer `yaml:"controller,omitempty"`
	Replicas       int                     `yaml:"replicas,omitempty"`
	Credentials    LegacyCredentials       `yaml:"credentials,omitempty"`
	NodeSelector   map[string]string       `yaml:"nodeSelector,omitempty"`
	Tolerations    []interface{}           `yaml:"tolerations,omitempty"`
	Affinity       map[string]interface{}  `yaml:"affinity,omitempty"`
	PodAnnotations map[string]string       `yaml:"podAnnotations,omitempty"`
	PodLabels      map[string]string       `yaml:"podLabels,omitempty"`
	ServiceAccountName string              `yaml:"serviceAccountName,omitempty"`
	Secret         string                  `yaml:"secret,omitempty"`
	RelativeURLs   bool                    `yaml:"relativeurls,omitempty"`
	Middleware     map[string]interface{}  `yaml:"middleware,omitempty"`
	UploadPurging  map[string]interface{}  `yaml:"upload_purging,omitempty"`
}

type LegacyRegistryContainer struct {
	Image        LegacyImage            `yaml:"image,omitempty"`
	Resources    map[string]interface{} `yaml:"resources,omitempty"`
	ExtraEnvVars []interface{}          `yaml:"extraEnvVars,omitempty"`
}

type LegacyCredentials struct {
	Username       string `yaml:"username,omitempty"`
	Password       string `yaml:"password,omitempty"`
	ExistingSecret string `yaml:"existingSecret,omitempty"`
	HtpasswdString string `yaml:"htpasswdString,omitempty"`
}

type LegacyJobservice struct {
	Image                 LegacyImage            `yaml:"image,omitempty"`
	Replicas              int                    `yaml:"replicas,omitempty"`
	MaxJobWorkers         int                    `yaml:"maxJobWorkers,omitempty"`
	JobLoggers            []string               `yaml:"jobLoggers,omitempty"`
	LoggerSweeperDuration int                    `yaml:"loggerSweeperDuration,omitempty"`
	Secret                string                 `yaml:"secret,omitempty"`
	Resources             map[string]interface{} `yaml:"resources,omitempty"`
	NodeSelector          map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations           []interface{}          `yaml:"tolerations,omitempty"`
	Affinity              map[string]interface{} `yaml:"affinity,omitempty"`
	PodAnnotations        map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels             map[string]string      `yaml:"podLabels,omitempty"`
	ExtraEnvVars          []interface{}          `yaml:"extraEnvVars,omitempty"`
	ServiceAccountName    string                 `yaml:"serviceAccountName,omitempty"`
	Notification          map[string]interface{} `yaml:"notification,omitempty"`
	Reaper                map[string]interface{} `yaml:"reaper,omitempty"`
}

type LegacyTrivy struct {
	Enabled            bool                   `yaml:"enabled,omitempty"`
	Image              LegacyImage            `yaml:"image,omitempty"`
	Replicas           int                    `yaml:"replicas,omitempty"`
	Resources          map[string]interface{} `yaml:"resources,omitempty"`
	NodeSelector       map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations        []interface{}          `yaml:"tolerations,omitempty"`
	Affinity           map[string]interface{} `yaml:"affinity,omitempty"`
	PodAnnotations     map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels          map[string]string      `yaml:"podLabels,omitempty"`
	ExtraEnvVars       []interface{}          `yaml:"extraEnvVars,omitempty"`
	ServiceAccountName string                 `yaml:"serviceAccountName,omitempty"`
	DebugMode          bool                   `yaml:"debugMode,omitempty"`
	VulnType           string                 `yaml:"vulnType,omitempty"`
	Severity           string                 `yaml:"severity,omitempty"`
	IgnoreUnfixed      bool                   `yaml:"ignoreUnfixed,omitempty"`
	Insecure           bool                   `yaml:"insecure,omitempty"`
	GitHubToken        string                 `yaml:"gitHubToken,omitempty"`
	SkipUpdate         bool                   `yaml:"skipUpdate,omitempty"`
	SkipJavaDBUpdate   bool                   `yaml:"skipJavaDBUpdate,omitempty"`
	OfflineScan        bool                   `yaml:"offlineScan,omitempty"`
	SecurityCheck      string                 `yaml:"securityCheck,omitempty"`
	Timeout            string                 `yaml:"timeout,omitempty"`
}

type LegacyExporter struct {
	Image              LegacyImage            `yaml:"image,omitempty"`
	Replicas           int                    `yaml:"replicas,omitempty"`
	Resources          map[string]interface{} `yaml:"resources,omitempty"`
	NodeSelector       map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations        []interface{}          `yaml:"tolerations,omitempty"`
	Affinity           map[string]interface{} `yaml:"affinity,omitempty"`
	PodAnnotations     map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels          map[string]string      `yaml:"podLabels,omitempty"`
	ServiceAccountName string                 `yaml:"serviceAccountName,omitempty"`
	CacheDuration      int                    `yaml:"cacheDuration,omitempty"`
	CacheCleanInterval int                    `yaml:"cacheCleanInterval,omitempty"`
}

type LegacyNginx struct {
	Image          LegacyImage            `yaml:"image,omitempty"`
	Replicas       int                    `yaml:"replicas,omitempty"`
	Resources      map[string]interface{} `yaml:"resources,omitempty"`
	NodeSelector   map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations    []interface{}          `yaml:"tolerations,omitempty"`
	Affinity       map[string]interface{} `yaml:"affinity,omitempty"`
	PodAnnotations map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels      map[string]string      `yaml:"podLabels,omitempty"`
}

type LegacyInternalTLS struct {
	Enabled          bool   `yaml:"enabled,omitempty"`
	StrongSSLCiphers bool   `yaml:"strong_ssl_ciphers,omitempty"`
	CertSource       string `yaml:"certSource,omitempty"`
}

type LegacyMetrics struct {
	Enabled        bool                   `yaml:"enabled,omitempty"`
	Core           map[string]interface{} `yaml:"core,omitempty"`
	Registry       map[string]interface{} `yaml:"registry,omitempty"`
	Jobservice     map[string]interface{} `yaml:"jobservice,omitempty"`
	Exporter       map[string]interface{} `yaml:"exporter,omitempty"`
	ServiceMonitor LegacyServiceMonitor   `yaml:"serviceMonitor,omitempty"`
}

type LegacyServiceMonitor struct {
	Enabled          bool              `yaml:"enabled,omitempty"`
	AdditionalLabels map[string]string `yaml:"additionalLabels,omitempty"`
	Interval         string            `yaml:"interval,omitempty"`
}

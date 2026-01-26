package migrate

// NewValues represents the structure of harbor-next-helm values.yaml
type NewValues struct {
	// Global settings
	Image               NewGlobalImage `yaml:"image,omitempty"`
	ImagePullSecrets    []string       `yaml:"imagePullSecrets,omitempty"`
	NameOverride        string         `yaml:"nameOverride,omitempty"`
	FullnameOverride    string         `yaml:"fullnameOverride,omitempty"`
	ExternalURL         string         `yaml:"externalURL,omitempty"`
	LogLevel            string         `yaml:"logLevel,omitempty"`
	HarborAdminPassword string         `yaml:"harborAdminPassword,omitempty"`
	SecretKey           string         `yaml:"secretKey,omitempty"`

	// Ingress
	Ingress NewIngress `yaml:"ingress,omitempty"`

	// Gateway API
	Gateway NewGateway `yaml:"gateway,omitempty"`

	// Components
	Core       NewCore       `yaml:"core,omitempty"`
	Registry   NewRegistry   `yaml:"registry,omitempty"`
	Portal     NewPortal     `yaml:"portal,omitempty"`
	Jobservice NewJobservice `yaml:"jobservice,omitempty"`
	Exporter   NewExporter   `yaml:"exporter,omitempty"`

	// Dependencies
	Database      NewDatabase      `yaml:"database,omitempty"`
	Valkey        NewValkey        `yaml:"valkey,omitempty"`
	ExternalRedis NewExternalRedis `yaml:"externalRedis,omitempty"`
	Trivy         NewTrivy         `yaml:"trivy,omitempty"`

	// TLS
	TLS NewTLS `yaml:"tls,omitempty"`

	// Metrics
	Metrics NewMetrics `yaml:"metrics,omitempty"`

	// Escape hatches
	ExtraManifests         []interface{} `yaml:"extraManifests,omitempty"`
	ExtraTemplateManifests []string      `yaml:"extraTemplateManifests,omitempty"`
}

type NewGlobalImage struct {
	PullPolicy string `yaml:"pullPolicy,omitempty"`
}

type NewIngress struct {
	Enabled     bool              `yaml:"enabled"`
	ClassName   string            `yaml:"className,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Hosts       []NewIngressHost  `yaml:"hosts,omitempty"`
	TLS         []NewIngressTLS   `yaml:"tls,omitempty"`
}

type NewIngressHost struct {
	Host  string           `yaml:"host"`
	Paths []NewIngressPath `yaml:"paths,omitempty"`
}

type NewIngressPath struct {
	Path     string `yaml:"path"`
	PathType string `yaml:"pathType"`
}

type NewIngressTLS struct {
	SecretName string   `yaml:"secretName"`
	Hosts      []string `yaml:"hosts"`
}

type NewGateway struct {
	Enabled    bool          `yaml:"enabled"`
	ParentRefs []interface{} `yaml:"parentRefs,omitempty"`
	Hostnames  []string      `yaml:"hostnames,omitempty"`
}

// NewComponentBase contains common fields for all components
type NewComponentBase struct {
	Config             map[string]interface{} `yaml:"config,omitempty"`
	Secret             map[string]interface{} `yaml:"secret,omitempty"`
	ExtraEnv           []interface{}          `yaml:"extraEnv,omitempty"`
	Image              NewImage               `yaml:"image,omitempty"`
	Replicas           int                    `yaml:"replicas,omitempty"`
	Resources          map[string]interface{} `yaml:"resources,omitempty"`
	NodeSelector       map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations        []interface{}          `yaml:"tolerations,omitempty"`
	Affinity           map[string]interface{} `yaml:"affinity,omitempty"`
	PodAnnotations     map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels          map[string]string      `yaml:"podLabels,omitempty"`
	ServiceAccount     NewServiceAccount      `yaml:"serviceAccount,omitempty"`
	SecurityContext    map[string]interface{} `yaml:"securityContext,omitempty"`
	PodSecurityContext map[string]interface{} `yaml:"podSecurityContext,omitempty"`
}

type NewImage struct {
	Repository string `yaml:"repository,omitempty"`
	Tag        string `yaml:"tag,omitempty"`
}

type NewServiceAccount struct {
	Create      bool              `yaml:"create"`
	Name        string            `yaml:"name,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type NewCore struct {
	Config             map[string]interface{} `yaml:"config,omitempty"`
	Secret             map[string]interface{} `yaml:"secret,omitempty"`
	ExtraEnv           []interface{}          `yaml:"extraEnv,omitempty"`
	Image              NewImage               `yaml:"image,omitempty"`
	Replicas           int                    `yaml:"replicas,omitempty"`
	Resources          map[string]interface{} `yaml:"resources,omitempty"`
	NodeSelector       map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations        []interface{}          `yaml:"tolerations,omitempty"`
	Affinity           map[string]interface{} `yaml:"affinity,omitempty"`
	PodAnnotations     map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels          map[string]string      `yaml:"podLabels,omitempty"`
	ServiceAccount     NewServiceAccount      `yaml:"serviceAccount,omitempty"`
	SecurityContext    map[string]interface{} `yaml:"securityContext,omitempty"`
	PodSecurityContext map[string]interface{} `yaml:"podSecurityContext,omitempty"`
}

type NewRegistry struct {
	Config             map[string]interface{} `yaml:"config,omitempty"`
	Secret             map[string]interface{} `yaml:"secret,omitempty"`
	ExtraEnv           []interface{}          `yaml:"extraEnv,omitempty"`
	Image              NewImage               `yaml:"image,omitempty"`
	Replicas           int                    `yaml:"replicas,omitempty"`
	Resources          map[string]interface{} `yaml:"resources,omitempty"`
	NodeSelector       map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations        []interface{}          `yaml:"tolerations,omitempty"`
	Affinity           map[string]interface{} `yaml:"affinity,omitempty"`
	PodAnnotations     map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels          map[string]string      `yaml:"podLabels,omitempty"`
	ServiceAccount     NewServiceAccount      `yaml:"serviceAccount,omitempty"`
	SecurityContext    map[string]interface{} `yaml:"securityContext,omitempty"`
	PodSecurityContext map[string]interface{} `yaml:"podSecurityContext,omitempty"`
	Storage            NewStorage             `yaml:"storage,omitempty"`
	Controller         NewController          `yaml:"controller,omitempty"`
	Persistence        NewPersistence         `yaml:"persistence,omitempty"`
}

type NewStorage struct {
	Type       string                 `yaml:"type,omitempty"`
	Filesystem map[string]interface{} `yaml:"filesystem,omitempty"`
	S3         map[string]interface{} `yaml:"s3,omitempty"`
	Azure      map[string]interface{} `yaml:"azure,omitempty"`
	GCS        map[string]interface{} `yaml:"gcs,omitempty"`
	OSS        map[string]interface{} `yaml:"oss,omitempty"`
}

type NewController struct {
	Image NewImage `yaml:"image,omitempty"`
}

type NewPersistence struct {
	Enabled       bool              `yaml:"enabled"`
	StorageClass  string            `yaml:"storageClass,omitempty"`
	AccessModes   []string          `yaml:"accessModes,omitempty"`
	Size          string            `yaml:"size,omitempty"`
	ExistingClaim string            `yaml:"existingClaim,omitempty"`
	Annotations   map[string]string `yaml:"annotations,omitempty"`
}

type NewPortal struct {
	Config             map[string]interface{} `yaml:"config,omitempty"`
	Secret             map[string]interface{} `yaml:"secret,omitempty"`
	ExtraEnv           []interface{}          `yaml:"extraEnv,omitempty"`
	Image              NewImage               `yaml:"image,omitempty"`
	Replicas           int                    `yaml:"replicas,omitempty"`
	Resources          map[string]interface{} `yaml:"resources,omitempty"`
	NodeSelector       map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations        []interface{}          `yaml:"tolerations,omitempty"`
	Affinity           map[string]interface{} `yaml:"affinity,omitempty"`
	PodAnnotations     map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels          map[string]string      `yaml:"podLabels,omitempty"`
	ServiceAccount     NewServiceAccount      `yaml:"serviceAccount,omitempty"`
	SecurityContext    map[string]interface{} `yaml:"securityContext,omitempty"`
	PodSecurityContext map[string]interface{} `yaml:"podSecurityContext,omitempty"`
}

type NewJobservice struct {
	Config             map[string]interface{} `yaml:"config,omitempty"`
	Secret             map[string]interface{} `yaml:"secret,omitempty"`
	ExtraEnv           []interface{}          `yaml:"extraEnv,omitempty"`
	Image              NewImage               `yaml:"image,omitempty"`
	Replicas           int                    `yaml:"replicas,omitempty"`
	Resources          map[string]interface{} `yaml:"resources,omitempty"`
	NodeSelector       map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations        []interface{}          `yaml:"tolerations,omitempty"`
	Affinity           map[string]interface{} `yaml:"affinity,omitempty"`
	PodAnnotations     map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels          map[string]string      `yaml:"podLabels,omitempty"`
	ServiceAccount     NewServiceAccount      `yaml:"serviceAccount,omitempty"`
	SecurityContext    map[string]interface{} `yaml:"securityContext,omitempty"`
	PodSecurityContext map[string]interface{} `yaml:"podSecurityContext,omitempty"`
}

type NewExporter struct {
	Enabled            bool                   `yaml:"enabled"`
	Config             map[string]interface{} `yaml:"config,omitempty"`
	Secret             map[string]interface{} `yaml:"secret,omitempty"`
	ExtraEnv           []interface{}          `yaml:"extraEnv,omitempty"`
	Image              NewImage               `yaml:"image,omitempty"`
	Replicas           int                    `yaml:"replicas,omitempty"`
	Resources          map[string]interface{} `yaml:"resources,omitempty"`
	NodeSelector       map[string]string      `yaml:"nodeSelector,omitempty"`
	Tolerations        []interface{}          `yaml:"tolerations,omitempty"`
	Affinity           map[string]interface{} `yaml:"affinity,omitempty"`
	PodAnnotations     map[string]string      `yaml:"podAnnotations,omitempty"`
	PodLabels          map[string]string      `yaml:"podLabels,omitempty"`
	ServiceAccount     NewServiceAccount      `yaml:"serviceAccount,omitempty"`
	SecurityContext    map[string]interface{} `yaml:"securityContext,omitempty"`
	PodSecurityContext map[string]interface{} `yaml:"podSecurityContext,omitempty"`
}

type NewDatabase struct {
	Host           string `yaml:"host,omitempty"`
	Port           int    `yaml:"port,omitempty"`
	Username       string `yaml:"username,omitempty"`
	Password       string `yaml:"password,omitempty"`
	Database       string `yaml:"database,omitempty"`
	SSLMode        string `yaml:"sslmode,omitempty"`
	MaxIdleConns   int    `yaml:"maxIdleConns,omitempty"`
	MaxOpenConns   int    `yaml:"maxOpenConns,omitempty"`
	ExistingSecret string `yaml:"existingSecret,omitempty"`
}

type NewValkey struct {
	Enabled      bool                   `yaml:"enabled"`
	Architecture string                 `yaml:"architecture,omitempty"`
	Auth         NewValkeyAuth          `yaml:"auth,omitempty"`
	Master       map[string]interface{} `yaml:"master,omitempty"`
}

type NewValkeyAuth struct {
	Enabled  bool   `yaml:"enabled"`
	Password string `yaml:"password,omitempty"`
}

type NewExternalRedis struct {
	Host              string `yaml:"host,omitempty"`
	Port              int    `yaml:"port,omitempty"`
	Password          string `yaml:"password,omitempty"`
	SentinelMasterSet string `yaml:"sentinelMasterSet,omitempty"`
	ExistingSecret    string `yaml:"existingSecret,omitempty"`
}

type NewTrivy struct {
	Enabled bool `yaml:"enabled"`
}

type NewTLS struct {
	CertManager   NewCertManager   `yaml:"certManager,omitempty"`
	CustomSecrets NewCustomSecrets `yaml:"customSecrets,omitempty"`
}

type NewCertManager struct {
	Enabled     bool                   `yaml:"enabled"`
	IssuerRef   map[string]interface{} `yaml:"issuerRef,omitempty"`
	Duration    string                 `yaml:"duration,omitempty"`
	RenewBefore string                 `yaml:"renewBefore,omitempty"`
}

type NewCustomSecrets struct {
	Core     string `yaml:"core,omitempty"`
	Registry string `yaml:"registry,omitempty"`
}

type NewMetrics struct {
	ServiceMonitor NewServiceMonitor `yaml:"serviceMonitor,omitempty"`
}

type NewServiceMonitor struct {
	Enabled       bool              `yaml:"enabled"`
	Namespace     string            `yaml:"namespace,omitempty"`
	Labels        map[string]string `yaml:"labels,omitempty"`
	Interval      string            `yaml:"interval,omitempty"`
	ScrapeTimeout string            `yaml:"scrapeTimeout,omitempty"`
	HonorLabels   bool              `yaml:"honorLabels,omitempty"`
}

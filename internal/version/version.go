package version

// Service é o identificador estável do serviço (health check e observabilidade).
const Service = "retech-core-api"

// Version é injetada no build: -ldflags '-X github.com/theretech/retech-core-api/internal/version.Version=1.2.3'
var Version = "dev"

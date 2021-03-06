package cloud

// Cloud related constants
const (
	defaultCloudResource     = "cloud"
	DefaultCloudResourcePath = "/" + defaultCloudResource

	DefaultWorkRoot           = "/var/tmp/cloud"
	DefaultTemplateRoot       = "./pkg/cloud/configs"
	DefaultGenInventoryScript = "transform/generate_inventories.py"
	DefaultAWSPlanTF          = "aws.tf.json"
	DefaultAzurePlanTF        = "azure.tf.json"
	DefaultGCPPlanTF          = "google.tf.json"
	DefaultTFStateFile        = "terraform.tfstate"
	DefaultTopologyFile       = "topology.yml"
	DefaultSecretFile         = "secret.yml"
	DefaultMultiCloudDir      = "/usr/share/contrail/"
	DefaultMultiCloudRepo     = "contrail-multi-cloud"

	MultiCloudDockerImage     = "contrail-multicloud-deployer"
	MultiCloudContainerPrefix = "multicloud-deployer"

	StatusField          = "provisioning_state"
	StatusNoState        = "NOSTATE"
	StatusCreated        = "CREATED"
	StatusCreateProgress = "CREATE_IN_PROGRESS"
	StatusCreateFailed   = "CREATE_FAILED"

	StatusUpdated        = "UPDATED"
	StatusUpdateProgress = "UPDATE_IN_PROGRESS"
	StatusUpdateFailed   = "UPDATE_FAILED"

	CreateAction      = "create"
	UpdateAction      = "update"
	DeleteCloudAction = "DELETE_CLOUD"

	aws    = "aws"
	azure  = "azure"
	gcp    = "gcp"
	google = "google"
	onPrem = "private"

	defaultRWOnlyPerm = 0600
	defaultSSHKeyRepo = "keypair"

	pubSSHKey     = "PUBLIC_SSH_KEY"
	privateSSHKey = "PRIVATE_SSH_KEY"
)

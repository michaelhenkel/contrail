package cluster

import (
	"github.com/Juniper/asf/pkg/logutil"
	"github.com/Juniper/contrail/pkg/ansible"
	"github.com/Juniper/contrail/pkg/client"
	"github.com/Juniper/contrail/pkg/deploy/base"
	"github.com/Juniper/contrail/pkg/deploy/rhospd/overcloud"
	"github.com/Juniper/contrail/pkg/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Config represents Command configuration.
type Config struct {
	// http client of api server
	APIServer *client.HTTP
	// UUID of resource to be managed.
	ClusterID string
	// Action to the performed with the cluster (values: create, update, delete).
	Action string
	// Logging level
	LogLevel string
	// Logging  file
	LogFile string
	// Template root directory
	TemplateRoot string
	// Work root directory
	WorkRoot string
	// Service user name for keystone
	ServiceUserID string
	// Service user password for keystone
	ServiceUserPassword string

	// Optional ansible sudo password
	AnsibleSudoPass string
	// Optional ansible deployer cherry pick url
	AnsibleFetchURL string
	// Optional ansible deployer cherry pick revison(commit id)
	AnsibleCherryPickRevision string
	// Optional ansible deployer revision(commit id)
	AnsibleRevision string
	// Optional Test var to run command in test mode
	Test bool
}

// Cluster represents contrail cluster manager
type Cluster struct {
	config            *Config
	APIServer         *client.HTTP
	log               *logrus.Entry
	containerExecutor ContainerExecutor
}

// NewCluster creates Cluster with given configuration.
func NewCluster(c *Config, containerExecutor ContainerExecutor) (*Cluster, error) {
	return &Cluster{
		config:            c,
		APIServer:         c.APIServer,
		log:               logutil.NewFileLogger("cluster", c.LogFile),
		containerExecutor: containerExecutor,
	}, nil
}

// GetDeployer creates new deployer based on the type
// TODO(Daniel): this should not be Cluster's method
func (c *Cluster) GetDeployer() (base.Deployer, error) {
	cData, err := data(c)
	if err != nil {
		return nil, err
	}

	switch deployerType(cData, c.config.Action) {
	case "rhospd":
		return newOvercloudDeployer(c)
	case "ansible", "tripleo", "juju":
		return newAnsibleDeployer(c, cData, c.containerExecutor)
	case MCProvisioner:
		return newMCProvisioner(c, cData, c.containerExecutor)
	}
	return nil, errors.New("unsupported deployer type")
}

func data(c *Cluster) (*base.Data, error) {
	if c.config.Action == "delete" {
		return &base.Data{Reader: c.APIServer}, nil
	}
	return base.NewResourceManager(c.APIServer, c.config.LogFile).GetClusterDetails(c.config.ClusterID)
}

func deployerType(cData *base.Data, action string) string {
	if cData.ClusterInfo != nil && cData.ClusterInfo.ProvisionerType != "" {
		return cData.ClusterInfo.ProvisionerType
	}

	if action != DeleteAction && isMCProvisioner(cData) {
		return MCProvisioner
	}

	return DefaultDeployer
}

func isMCProvisioner(cData *base.Data) bool {
	if hasCloudRefs(cData) && hasMCGWNodes(cData.ClusterInfo) {
		switch cData.ClusterInfo.ProvisioningAction {
		case addCloud, updateCloud, deleteCloud:
			return true
		}
	}
	return false
}

func hasCloudRefs(d *base.Data) bool {
	return d.CloudInfo != nil
}

func hasMCGWNodes(cc *models.ContrailCluster) bool {
	return cc.ContrailMulticloudGWNodes != nil
}

func newOvercloudDeployer(c *Cluster) (base.Deployer, error) {
	o, err := overcloud.NewOverCloud(&overcloud.Config{
		APIServer:    c.APIServer,
		ResourceID:   c.config.ClusterID,
		Action:       c.config.Action,
		TemplateRoot: c.config.TemplateRoot,
		LogLevel:     c.config.LogLevel,
		LogFile:      c.config.LogFile,
	})
	if err != nil {
		return nil, err
	}

	return o.GetDeployer()
}

func newAnsibleDeployer(
	c *Cluster, cData *base.Data, containerExecutor ContainerExecutor,
) (*ContrailAnsibleDeployer, error) {
	d := newDeployCluster(c, cData, "contrail-ansible-deployer")

	return &ContrailAnsibleDeployer{
		deployCluster: *d,
		ansibleClient: ansible.NewCLIClient(
			d.Reporter,
			c.config.LogFile,
			d.getWorkingDir(),
			c.config.Test,
		),
		containerExecutor: containerExecutor,
	}, nil
}

func newMCProvisioner(
	c *Cluster, cData *base.Data, containerExecutor ContainerExecutor,
) (*multiCloudProvisioner, error) {
	d := newDeployCluster(c, cData, "multi-cloud-provisioner")

	return &multiCloudProvisioner{
		ContrailAnsibleDeployer: ContrailAnsibleDeployer{
			deployCluster: *d,
			ansibleClient: ansible.NewCLIClient(
				d.Reporter,
				c.config.LogFile,
				d.getWorkingDir(),
				c.config.Test,
			),
			containerExecutor: containerExecutor,
		},
		workDir: "",
	}, nil
}

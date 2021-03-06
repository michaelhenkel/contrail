package deploy

import (
	"fmt"

	"github.com/Juniper/asf/pkg/logutil"
	"github.com/Juniper/asf/pkg/logutil/report"
	"github.com/Juniper/contrail/pkg/ansible"
	"github.com/Juniper/contrail/pkg/deploy/base"
	"github.com/Juniper/contrail/pkg/deploy/cluster"
	"github.com/Juniper/contrail/pkg/deploy/rhospd/undercloud"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// manager inerface to manage resources
type manager interface {
	manage() error
}

type oneShotManager struct {
	deploy *Deploy
	log    *logrus.Entry
}

func (o *oneShotManager) manage() error {
	deployer, err := newDeployer(o.deploy)
	if err != nil {
		return err
	}
	return deployer.Deploy()
}

func newOneShotManager(deploy *Deploy) (*oneShotManager, error) {
	return &oneShotManager{
		deploy: deploy,
		log:    logutil.NewLogger("oneshot-manager"),
	}, nil
}

func newManager(deploy *Deploy) (manager, error) {
	switch deploy.managerType {
	case "oneshot":
		return newOneShotManager(deploy)
	}
	//TODO(ijohnson) Support daemon manager with etcd
	return nil, errors.New("unsupported manager type")
}

func newDeployer(deploy *Deploy) (base.Deployer, error) {
	switch deploy.config.ResourceType {
	case "contrail_cluster":
		player, err := ansible.NewContainerPlayer(report.NewReporter(
			deploy.APIServer,
			fmt.Sprintf("%s/%s", cluster.DefaultResourcePath, deploy.config.ResourceID),
			logutil.NewFileLogger("reporter", deploy.config.LogFile).WithField("cloudID", deploy.config.ResourceID),
		), deploy.config.LogFile)
		if err != nil {
			return nil, err
		}
		c, err := cluster.NewCluster(&cluster.Config{
			APIServer:                 deploy.APIServer,
			ClusterID:                 deploy.config.ResourceID,
			Action:                    deploy.config.Action,
			TemplateRoot:              deploy.config.TemplateRoot,
			LogLevel:                  deploy.config.LogLevel,
			LogFile:                   deploy.config.LogFile,
			AnsibleSudoPass:           deploy.config.AnsibleSudoPass,
			AnsibleFetchURL:           deploy.config.AnsibleFetchURL,
			AnsibleCherryPickRevision: deploy.config.AnsibleCherryPickRevision,
			AnsibleRevision:           deploy.config.AnsibleRevision,
			ServiceUserID:             deploy.config.ServiceUserID,
			ServiceUserPassword:       deploy.config.ServiceUserPassword,
		}, player)
		if err != nil {
			return nil, err
		}
		return c.GetDeployer()
	case "rhospd_cloud_manager":
		u, err := undercloud.NewUnderCloud(&undercloud.Config{
			APIServer:    deploy.APIServer,
			ResourceID:   deploy.config.ResourceID,
			Action:       deploy.config.Action,
			TemplateRoot: deploy.config.TemplateRoot,
			LogLevel:     deploy.config.LogLevel,
			LogFile:      deploy.config.LogFile,
		})
		if err != nil {
			return nil, err
		}
		return u.GetDeployer()
	}
	return nil, errors.New("unsupported resource type")
}

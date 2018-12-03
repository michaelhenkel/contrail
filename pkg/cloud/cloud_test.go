package cloud

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/flosch/pongo2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/Juniper/contrail/pkg/apisrv"
)

const (
	allInOneCloudTemplatePath       = "./test_data/test_all_in_one_azure_cloud.tmpl"
	allInOneCloudDeleteTemplatePath = "./test_data/test_all_in_one_azure_cloud_delete.tmpl"
	allInOneCloudUpdateTemplatePath = "./test_data/test_all_in_one_azure_cloud_update.tmpl"
	expectedAZCmdForCreateUpdate    = "./test_data/expected_azure_cmd_for_create_update.yaml"
	expectedAZTopology              = "./test_data/expected_azure_cloud_topology.yaml"
	expectedAZSecret                = "./test_data/expected_azure_cloud_secret.yaml"
	cloudID                         = "cloud_uuid"
)

func TestMain(m *testing.M) {
	apisrv.SetupAndRunTest(m)
}

func TestAzureCloud(t *testing.T) {
	runAllInOneCloudTest(t, azure)
}

func runAllInOneCloudTest(t *testing.T, cloudType string) {

	context := pongo2.Context{
		"CLOUD_TYPE": cloudType,
	}

	switch cloudType {
	case azure:
		runCloudTest(t, expectedAZTopology, expectedAZSecret,
			expectedAZCmdForCreateUpdate, context)
	}
}

func runCloudTest(t *testing.T, expectedTopology, expectedSecret string,
	expectedCmdForCreateUpdate string, context map[string]interface{}) {

	// mock keystone to let access server after cluster create
	keystoneAuthURL := viper.GetString("keystone.authurl")
	ksPublic := apisrv.MockServerWithKeystone("127.0.0.1:35357", keystoneAuthURL)
	defer ksPublic.Close()
	ksPrivate := apisrv.MockServerWithKeystone("127.0.0.1:5000", keystoneAuthURL)
	defer ksPrivate.Close()

	// create cloud related objects
	var cloudTestScenario apisrv.TestScenario
	err := apisrv.LoadTestScenario(&cloudTestScenario, allInOneCloudTemplatePath, context)
	assert.NoErrorf(t, err, "failed to load cloud test data from file: %s", allInOneCloudTemplatePath)
	cleanup := apisrv.RunDirtyTestScenario(t, &cloudTestScenario)

	defer cleanup()

	// creating cloud config
	config := &Config{
		ID:           "alice",
		Password:     "alice_password",
		ProjectID:    "admin",
		AuthURL:      apisrv.TestServer.URL + "/keystone/v3",
		Endpoint:     apisrv.TestServer.URL,
		InSecure:     true,
		CloudID:      cloudID,
		Action:       createAction,
		LogLevel:     "debug",
		TemplateRoot: "configs/",
		Test:         true,
	}

	// delete previously created
	if _, err = os.Stat(executedCommandsPath()); err == nil {
		// cleanup old executed command file
		err = os.Remove(executedCommandsPath())
		if err != nil {
			assert.NoError(t, err, "failed to delete executed commands yaml")
		}
	}
	cloud, err := NewCloud(config)
	assert.NoError(t, err, "failed to create cloud struct")

	err = cloud.Manage()
	assert.NoError(t, err, "failed to manage cloud, while creating cloud")

	assert.True(t, compareGeneratedTopology(t, expectedTopology),
		"topology file created during cloud create is not as expected")
	assert.True(t, compareGeneratedSecret(t, expectedSecret),
		"secret file created during cloud create is not as expected")
	assert.True(t, verifyCommandsExecuted(t, expectedCmdForCreateUpdate),
		"Expected list of create commands are not executed")

	// Wait for the in-memory endpoint cache to get updated
	apisrv.APIServer.ForceProxyUpdate()

	//update cloud
	config.Action = updateAction

	var cloudUpdateTestScenario apisrv.TestScenario
	err = apisrv.LoadTestScenario(&cloudUpdateTestScenario, allInOneCloudUpdateTemplatePath, context)
	assert.NoErrorf(t, err, "failed to load cloud test data from file: %s", allInOneCloudUpdateTemplatePath)
	_ = apisrv.RunDirtyTestScenario(t, &cloudUpdateTestScenario)

	// delete previously created files

	// Remove topology file and secret file
	err = os.Remove(generatedTopoPath())
	if err != nil {
		assert.NoError(t, err, "failed to delete topology.yml file, during update")
	}
	err = os.Remove(generatedSecretPath())
	if err != nil {
		assert.NoError(t, err, "failed to delete secret.yml file, during update")
	}

	if _, err = os.Stat(executedCommandsPath()); err == nil {
		// cleanup old executed command file
		err = os.Remove(executedCommandsPath())
		if err != nil {
			assert.NoError(t, err, "failed to delete executed cmd yml, during update")
		}
	}

	cloud, err = NewCloud(config)
	assert.NoError(t, err, "failed to create cloud struct for update action")

	err = cloud.Manage()
	assert.NoError(t, err, "failed to manage cloud, while updating cloud")

	assert.True(t, compareGeneratedTopology(t, expectedTopology),
		"topology file created during cloud update is not as expected")
	assert.True(t, compareGeneratedSecret(t, expectedSecret),
		"secret file created during cloud update is not as expected")
	assert.True(t, verifyCommandsExecuted(t, expectedCmdForCreateUpdate),
		"Expected list of update commands are not executed")

	// delete cloud
	var cloudDeleteTestScenario apisrv.TestScenario
	err = apisrv.LoadTestScenario(&cloudDeleteTestScenario, allInOneCloudDeleteTemplatePath, context)
	assert.NoErrorf(t, err, "failed to load cloud test data from file: %s", allInOneCloudDeleteTemplatePath)
	_ = apisrv.RunDirtyTestScenario(t, &cloudDeleteTestScenario)

	// delete previously created
	if _, err = os.Stat(executedCommandsPath()); err == nil {
		// cleanup old executed command file
		err = os.Remove(executedCommandsPath())
		if err != nil {
			assert.NoError(t, err, "failed to delete executed cmd yml, during delete")
		}
	}

	cloud, err = NewCloud(config)
	assert.NoError(t, err, "failed to create cloud struct for delete action")

	err = cloud.Manage()
	assert.NoError(t, err, "failed to manage cloud, while deleting cloud")

	// make sure cloud is removed
	assert.True(t, verifyCloudDeleted(), "Topo/Secret file is not deleted during cloud delete")

}

func compareFiles(t *testing.T, expectedFile, generatedFile string) bool {
	generatedData, err := ioutil.ReadFile(generatedFile)
	assert.NoErrorf(t, err, "unable to read generated: %s", generatedFile)
	expectedData, err := ioutil.ReadFile(expectedFile)
	assert.NoErrorf(t, err, "unable to read expected: %s", expectedFile)
	return bytes.Equal(generatedData, expectedData)

}

func compareGeneratedTopology(t *testing.T, expectedTopology string) bool {
	return compareFiles(t, expectedTopology, generatedTopoPath())
}

func compareGeneratedSecret(t *testing.T, expectedSecretFile string) bool {
	return compareFiles(t, expectedSecretFile, generatedSecretPath())
}

func verifyCommandsExecuted(t *testing.T, expectedCmdForCreateUpdate string) bool {
	return compareFiles(t, expectedCmdForCreateUpdate, executedCommandsPath())
}

func generatedTopoPath() string {
	return defaultWorkRoot + "/" + cloudID + "/topology.yml"
}

//new-comment

func generatedSecretPath() string {
	return defaultWorkRoot + "/" + cloudID + "/secret.yml"
}

func executedCommandsPath() string {
	return defaultWorkRoot + "/" + cloudID + "/" + executedCmdTestFile
}

func verifyCloudDeleted() bool {

	if _, err := os.Stat(defaultWorkRoot + "/" + cloudID); err == nil {
		// working dir not deleted
		return false
	}
	return true

}
package deployer

import (
	"context"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/grid3-go/mocks"
	client "github.com/threefoldtech/grid3-go/node"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/substrate-client"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

func constructTestK8s(t *testing.T, mock bool) (
	K8sDeployer,
	*mocks.RMBMockClient,
	*mocks.MockSubstrateExt,
	*mocks.MockNodeClientGetter,
	*mocks.MockDeployerInterface,
	*mocks.MockClient,
) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tfPluginClient, err := setup()
	assert.NoError(t, err)

	cl := mocks.NewRMBMockClient(ctrl)
	sub := mocks.NewMockSubstrateExt(ctrl)
	ncPool := mocks.NewMockNodeClientGetter(ctrl)
	deployer := mocks.NewMockDeployerInterface(ctrl)
	gridProxyCl := mocks.NewMockClient(ctrl)

	if mock {
		tfPluginClient.twinID = twinID

		tfPluginClient.SubstrateConn = sub
		tfPluginClient.NcPool = ncPool
		tfPluginClient.RMB = cl
		tfPluginClient.GridProxyClient = gridProxyCl

		tfPluginClient.StateLoader.ncPool = ncPool
		tfPluginClient.StateLoader.substrate = sub

		tfPluginClient.K8sDeployer.deployer = deployer
		tfPluginClient.K8sDeployer.tfPluginClient = &tfPluginClient
	}
	net := constructTestNetwork()
	tfPluginClient.StateLoader.networks = networkState{net.Name: network{
		subnets:               map[uint32]string{nodeID: net.IPRange.String()},
		nodeDeploymentHostIDs: map[uint32]deploymentHostIDs{nodeID: map[uint64][]byte{contractID: {}}},
	}}

	return tfPluginClient.K8sDeployer, cl, sub, ncPool, deployer, gridProxyCl
}

func k8sMockValidation(identity substrate.Identity, cl *mocks.RMBMockClient, sub *mocks.MockSubstrateExt, ncPool *mocks.MockNodeClientGetter, proxyCl *mocks.MockClient, d K8sDeployer) {
	sub.EXPECT().
		GetBalance(d.tfPluginClient.identity).
		Return(substrate.Balance{
			Free: types.U128{
				Int: big.NewInt(100000),
			},
		}, nil)

	cl.EXPECT().
		Call(
			gomock.Any(),
			nodeID,
			"zos.system.version",
			nil,
			gomock.Any(),
		).Return(nil).AnyTimes()

	ncPool.EXPECT().
		GetNodeClient(
			gomock.Any(),
			nodeID,
		).Return(client.NewNodeClient(nodeID, cl), nil)

}

func constructK8sCluster() (workloads.K8sCluster, error) {
	flist := "https://hub.grid.tf/tf-official-apps/threefoldtech-k3s-latest.flist"
	flistCheckSum, err := workloads.GetFlistChecksum(flist)
	if err != nil {
		return workloads.K8sCluster{}, err
	}

	master := workloads.K8sNode{
		Name:          "K8sforTesting",
		Node:          nodeID,
		DiskSize:      5,
		PublicIP:      true,
		PublicIP6:     true,
		Planetary:     true,
		Flist:         flist,
		FlistChecksum: flistCheckSum,
		ComputedIP:    "5.5.5.5/24",
		ComputedIP6:   "::7/64",
		YggIP:         "::8/64",
		IP:            "10.1.0.2",
		CPU:           2,
		Memory:        1024,
	}

	worker := workloads.K8sNode{
		Name:          "worker1",
		Node:          nodeID,
		DiskSize:      5,
		PublicIP:      true,
		PublicIP6:     true,
		Planetary:     true,
		Flist:         flist,
		FlistChecksum: flistCheckSum,
		ComputedIP:    "",
		ComputedIP6:   "",
		YggIP:         "",
		IP:            "",
		CPU:           2,
		Memory:        1024,
	}
	workers := []workloads.K8sNode{worker}
	Cluster := workloads.K8sCluster{
		Master:       &master,
		Workers:      workers[:],
		Token:        "tokens",
		SSHKey:       "",
		NetworkName:  "network",
		NodesIPRange: make(map[uint32]gridtypes.IPNet),
		// NodeDeploymentID: map[uint32]uint64{nodeID: contractID},
		ContractID: 0,
	}
	return Cluster, nil
}

func TestValidateMasterReachable(t *testing.T) {
	d, cl, sub, ncPool, _, proxyCl := constructTestK8s(t, true)
	k8sMockValidation(d.tfPluginClient.identity, cl, sub, ncPool, proxyCl, d)

	k8s, err := constructK8sCluster()
	assert.NoError(t, err)

	err = d.assignNodeIPRange(&k8s)
	assert.NoError(t, err)

	err = d.Validate(context.Background(), &k8s)
	assert.NoError(t, err)
}

func TestGenerateK8sDeployment(t *testing.T) {
	d, _, _, _, _, _ := constructTestK8s(t, true)
	k8s, err := constructK8sCluster()
	assert.NoError(t, err)

	err = d.assignNodeIPRange(&k8s)
	assert.NoError(t, err)

	dls, err := d.GenerateVersionlessDeployments(context.Background(), &k8s)
	assert.NoError(t, err)

	nodeWorkloads := make(map[uint32][]gridtypes.Workload)
	masterWorkloads := k8s.Master.MasterZosWorkload(&k8s)
	nodeWorkloads[k8s.Master.Node] = append(nodeWorkloads[k8s.Master.Node], masterWorkloads...)
	for _, w := range k8s.Workers {
		workerWorkloads := w.WorkerZosWorkload(&k8s)
		nodeWorkloads[w.Node] = append(nodeWorkloads[w.Node], workerWorkloads...)
	}

	wl := nodeWorkloads[nodeID]
	assert.Equal(t, dls, map[uint32]gridtypes.Deployment{
		nodeID: workloads.NewGridDeployment(d.tfPluginClient.twinID, wl),
	})
}

func TestDeploy(t *testing.T) {
	d, cl, sub, ncPool, deployer, proxyCl := constructTestK8s(t, true)

	k8sCluster, err := constructK8sCluster()
	assert.NoError(t, err)

	err = d.assignNodeIPRange(&k8sCluster)
	assert.NoError(t, err)

	dls, err := d.GenerateVersionlessDeployments(context.Background(), &k8sCluster)
	assert.NoError(t, err)

	k8sMockValidation(d.tfPluginClient.identity, cl, sub, ncPool, proxyCl, d)

	deploymentData := workloads.DeploymentData{
		Name:        k8sCluster.Master.Name,
		Type:        "K8s",
		ProjectName: "",
	}
	newDeploymentsData := make(map[uint32]workloads.DeploymentData)
	newDeploymentsSolutionProvider := make(map[uint32]*uint64)

	newDeploymentsData[k8sCluster.Master.Node] = deploymentData
	newDeploymentsSolutionProvider[k8sCluster.Master.Node] = nil

	deployer.EXPECT().Deploy(
		gomock.Any(),
		map[uint32]uint64{},
		dls,
		newDeploymentsData,
		newDeploymentsSolutionProvider,
	).Return(map[uint32]uint64{nodeID: contractID}, nil)

	err = d.Deploy(context.Background(), &k8sCluster)
	assert.NoError(t, err)

	assert.NotEqual(t, k8sCluster.ContractID, 0)
	assert.Equal(t, k8sCluster.NodeDeploymentID, map[uint32]uint64{nodeID: contractID})
}

func TestUpdateK8s(t *testing.T) {
	d, cl, sub, ncPool, deployer, proxyCl := constructTestK8s(t, true)

	k8sCluster, err := constructK8sCluster()
	assert.NoError(t, err)

	d.tfPluginClient.StateLoader.currentNodeDeployment[nodeID] = contractID

	err = d.assignNodeIPRange(&k8sCluster)
	assert.NoError(t, err)

	dls, err := d.GenerateVersionlessDeployments(context.Background(), &k8sCluster)
	assert.NoError(t, err)

	k8sMockValidation(d.tfPluginClient.identity, cl, sub, ncPool, proxyCl, d)

	deployer.EXPECT().Deploy(
		gomock.Any(),
		map[uint32]uint64{nodeID: contractID},
		dls,
		gomock.Any(),
		gomock.Any(),
	).Return(map[uint32]uint64{nodeID: contractID}, nil)

	err = d.Deploy(context.Background(), &k8sCluster)
	assert.NoError(t, err)
	assert.Equal(t, k8sCluster.NodeDeploymentID, map[uint32]uint64{nodeID: contractID})

}

func TestUpdateK8sFailed(t *testing.T) {
	d, cl, sub, ncPool, deployer, proxyCl := constructTestK8s(t, true)

	k8sCluster, err := constructK8sCluster()
	assert.NoError(t, err)

	d.tfPluginClient.StateLoader.currentNodeDeployment[nodeID] = contractID

	err = d.assignNodeIPRange(&k8sCluster)
	assert.NoError(t, err)

	dls, err := d.GenerateVersionlessDeployments(context.Background(), &k8sCluster)
	assert.NoError(t, err)

	k8sMockValidation(d.tfPluginClient.identity, cl, sub, ncPool, proxyCl, d)

	deployer.EXPECT().Deploy(
		gomock.Any(),
		map[uint32]uint64{nodeID: contractID},
		dls,
		gomock.Any(),
		gomock.Any(),
	).Return(map[uint32]uint64{nodeID: contractID}, errors.New("error"))

	err = d.Deploy(context.Background(), &k8sCluster)
	assert.Error(t, err)
	assert.Equal(t, k8sCluster.NodeDeploymentID, map[uint32]uint64{nodeID: contractID})
}

func TestCancelK8s(t *testing.T) {
	d, cl, sub, ncPool, deployer, proxyCl := constructTestK8s(t, true)

	d.tfPluginClient.StateLoader.currentNodeDeployment[nodeID] = contractID

	k8sCluster, err := constructK8sCluster()
	assert.NoError(t, err)
	k8sCluster.NodesIPRange = map[uint32]gridtypes.IPNet{uint32(10): {}}

	k8sMockValidation(d.tfPluginClient.identity, cl, sub, ncPool, proxyCl, d)

	deployer.EXPECT().Cancel(
		gomock.Any(), contractID,
	).Return(nil).AnyTimes()

	err = d.Cancel(context.Background(), &k8sCluster)
	assert.NoError(t, err)
	assert.Empty(t, k8sCluster.NodeDeploymentID)
}

func TestCancelK8sFailed(t *testing.T) {
	d, cl, sub, ncPool, deployer, proxyCl := constructTestK8s(t, true)

	d.tfPluginClient.StateLoader.currentNodeDeployment[nodeID] = contractID

	k8sCluster, err := constructK8sCluster()
	k8sCluster.NodeDeploymentID = map[uint32]uint64{nodeID: contractID}
	assert.NoError(t, err)
	k8sCluster.NodesIPRange = map[uint32]gridtypes.IPNet{uint32(10): {}}

	k8sMockValidation(d.tfPluginClient.identity, cl, sub, ncPool, proxyCl, d)

	deployer.EXPECT().Cancel(
		gomock.Any(), contractID,
	).Return(errors.New("error"))

	err = d.Cancel(context.Background(), &k8sCluster)
	assert.Error(t, err)
	assert.Equal(t, k8sCluster.NodeDeploymentID, map[uint32]uint64{nodeID: contractID})
}

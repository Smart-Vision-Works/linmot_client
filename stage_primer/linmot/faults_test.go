package linmot

import (
	"context"
	"errors"
	"testing"
	"time"

	linmot_client "github.com/Smart-Vision-Works/staged_robot/client"
	client_command_tables "github.com/Smart-Vision-Works/staged_robot/client/rtc/command_tables"
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type faultTestClient struct {
	checkCalls int
	closeCalls int
	ackCalls   int
	ackErr     error
	lastVAI    struct {
		pos float64
		vel float64
		acc float64
		dec float64
	}
}

func (c *faultTestClient) Close() error {
	c.closeCalls++
	return nil
}

func (c *faultTestClient) GetPosition(ctx context.Context) (float64, error) {
	return 0, nil
}

func (c *faultTestClient) GetStatus(ctx context.Context) (*protocol_common.Status, error) {
	return &protocol_common.Status{}, nil
}

func (c *faultTestClient) CheckDriveFault(ctx context.Context) error {
	c.checkCalls++
	return nil
}

func (c *faultTestClient) EnableDrive(ctx context.Context) (*protocol_common.Status, error) {
	return &protocol_common.Status{}, nil
}

func (c *faultTestClient) AcknowledgeError(ctx context.Context) (*protocol_common.Status, error) {
	c.ackCalls++
	if c.ackErr != nil {
		return nil, c.ackErr
	}
	return &protocol_common.Status{}, nil
}

func (c *faultTestClient) SendControlWord(ctx context.Context, word uint16) (*protocol_common.Status, error) {
	return &protocol_common.Status{}, nil
}

func (c *faultTestClient) VAIGoToPosition(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	c.lastVAI.pos = positionMM
	c.lastVAI.vel = velocityMS
	c.lastVAI.acc = accelMS2
	c.lastVAI.dec = decelMS2
	return &protocol_common.Status{}, nil
}

func (c *faultTestClient) VAIGoToPositionFromActual(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	c.lastVAI.pos = positionMM
	c.lastVAI.vel = velocityMS
	c.lastVAI.acc = accelMS2
	c.lastVAI.dec = decelMS2
	return &protocol_common.Status{}, nil
}

func (c *faultTestClient) SetPosition1(ctx context.Context, position float64, storage protocol_common.ParameterStorageType) error {
	return nil
}

func (c *faultTestClient) SetPosition2(ctx context.Context, position float64, storage protocol_common.ParameterStorageType) error {
	return nil
}

func (c *faultTestClient) SetRunMode(ctx context.Context, mode protocol_common.RunMode, storageType protocol_common.ParameterStorageType) error {
	return nil
}

func (c *faultTestClient) SetCommandTable(ctx context.Context, table *client_command_tables.CommandTable) error {
	return nil
}

func (c *faultTestClient) WriteOutputsWithMask(ctx context.Context, bitMask, bitValue uint16) (*protocol_common.Status, error) {
	return nil, nil
}

func (c *faultTestClient) SetEasyStepsAutoStart(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error {
	return nil
}

func (c *faultTestClient) SetEasyStepsAutoHome(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error {
	return nil
}

func (c *faultTestClient) SetEasyStepsInputRisingEdgeFunction(ctx context.Context, inputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error {
	return nil
}

func (c *faultTestClient) SetEasyStepsInputCurveCmdID(ctx context.Context, inputNumber protocol_common.IOPinNumber, curveCmdID int32, storageType protocol_common.ParameterStorageType) error {
	return nil
}

func (c *faultTestClient) SetIODefOutputFunction(ctx context.Context, outputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error {
	return nil
}

func (c *faultTestClient) SetIODefInputFunction(ctx context.Context, inputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error {
	return nil
}

func (c *faultTestClient) SetTriggerMode(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error {
	return nil
}
func (c *faultTestClient) Home(ctx context.Context) (*protocol_common.Status, error) {
	return &protocol_common.Status{}, nil
}
func (c *faultTestClient) SetCommandTableWithOptions(context.Context, *client_command_tables.CommandTable, client_command_tables.SetCommandTableOptions) error {
	return nil
}
func (c *faultTestClient) GetCommandTable(context.Context) (*client_command_tables.CommandTable, error) {
	return &client_command_tables.CommandTable{}, nil
}
func (c *faultTestClient) StopMotionController(context.Context) error    { return nil }
func (c *faultTestClient) StartMotionController(context.Context) error   { return nil }
func (c *faultTestClient) SaveCommandTableToFlash(context.Context) error { return nil }
func (c *faultTestClient) ReadRAM(context.Context, uint16) (int32, error) {
	return 0, nil
}
func (c *faultTestClient) WriteRAMAndROM(context.Context, uint16, int32) error {
	return nil
}

type faultTestFactory struct {
	client LinMotClient
}

func (f *faultTestFactory) CreateClient(linmotIP string) (LinMotClient, error) {
	return f.client, nil
}

func (f *faultTestFactory) Close() {}

func resetFaultHandlingState(ip string) {
	ResetFaultBudget(ip)
	resetFaultLifecycleStateForTests()
}

func TestCheckDriveFaultOnce_DoesNotCloseFactoryManagedClient(t *testing.T) {
	mockClient := &faultTestClient{}
	SetClientFactory(&faultTestFactory{client: mockClient})
	defer ResetClientFactory()

	err := checkDriveFaultOnce(context.Background(), "192.168.1.100")
	require.NoError(t, err)

	assert.Equal(t, 1, mockClient.checkCalls)
	assert.Equal(t, 0, mockClient.closeCalls)
}

func TestAddFaultListener_UnregistersListener(t *testing.T) {
	calls := 0
	remove := AddFaultListener(func(ip string, err error) {
		calls++
	})

	broadcastFault("192.168.1.10", errors.New("fault-1"))
	assert.Equal(t, 1, calls)

	remove()
	broadcastFault("192.168.1.10", errors.New("fault-2"))
	assert.Equal(t, 1, calls)
}

func TestHandleDetectedFault_NonFatalAutoRecoverSuccess(t *testing.T) {
	ip := "10.0.0.7"
	resetFaultHandlingState(ip)

	mockClient := &faultTestClient{}
	SetClientFactory(&faultTestFactory{client: mockClient})
	defer ResetClientFactory()

	faultCalls := 0
	escalationCalls := 0
	removeFault := AddFaultListener(func(_ string, _ error) {
		faultCalls++
	})
	defer removeFault()
	removeEscalation := AddEscalationListener(func(_ string, _ *linmot_client.DriveFaultError, _ bool) {
		escalationCalls++
	})
	t.Cleanup(removeEscalation)

	fault := &linmot_client.DriveFaultError{
		StatusWord: 0x0000,
		ErrorCode:  0x0042,
	}
	handleDetectedFault(context.Background(), ip, fault)

	assert.Equal(t, 1, mockClient.ackCalls, "auto-recovery should acknowledge once")
	assert.Equal(t, 0, faultCalls, "successful auto-recovery should suppress fault broadcast")
	assert.Equal(t, 0, escalationCalls, "successful auto-recovery should suppress escalation")
	snapshot := getFaultLifecycleSnapshot(ip)
	assert.Equal(t, FaultLifecycleStateHealthy, snapshot.State)
	assert.NotZero(t, snapshot.LastRecoverAttemptTime.UnixNano())
	assert.Zero(t, snapshot.ConsecutiveRecoverFailures)
}

func TestHandleDetectedFault_NonFatalAutoRecoverFailureStaysRecoveringWithinBudget(t *testing.T) {
	ip := "10.0.0.8"
	resetFaultHandlingState(ip)

	mockClient := &faultTestClient{ackErr: errors.New("ack failed")}
	SetClientFactory(&faultTestFactory{client: mockClient})
	defer ResetClientFactory()

	faultCalls := 0
	var gotEscalationFatal bool
	escalationCalls := 0
	removeFault := AddFaultListener(func(_ string, _ error) {
		faultCalls++
	})
	defer removeFault()
	removeEscalation := AddEscalationListener(func(_ string, _ *linmot_client.DriveFaultError, isFatal bool) {
		escalationCalls++
		gotEscalationFatal = isFatal
	})
	t.Cleanup(removeEscalation)

	fault := &linmot_client.DriveFaultError{
		StatusWord: 0x0000,
		ErrorCode:  0x0042,
	}
	handleDetectedFault(context.Background(), ip, fault)

	assert.Equal(t, 1, mockClient.ackCalls)
	assert.Equal(t, 1, faultCalls, "first failed auto-recovery should broadcast the fault once")
	assert.Equal(t, 0, escalationCalls, "failed auto-recovery should not escalate until the retry budget is exhausted")
	assert.False(t, gotEscalationFatal)
	snapshot := getFaultLifecycleSnapshot(ip)
	assert.Equal(t, FaultLifecycleStateRecovering, snapshot.State)
	assert.Equal(t, 1, snapshot.ConsecutiveRecoverFailures)
	assert.True(t, snapshot.LastEscalationTime.IsZero())
}

func TestHandleDetectedFault_BudgetExceededEscalatesWithoutAcknowledge(t *testing.T) {
	ip := "10.0.0.9"
	resetFaultHandlingState(ip)
	for i := 0; i < autoRecoveryMaxFaults; i++ {
		assert.True(t, globalFaultBudget.tryConsume(ip))
	}

	mockClient := &faultTestClient{}
	SetClientFactory(&faultTestFactory{client: mockClient})
	defer ResetClientFactory()

	faultCalls := 0
	escalationCalls := 0
	removeFault := AddFaultListener(func(_ string, _ error) {
		faultCalls++
	})
	defer removeFault()
	removeEscalation := AddEscalationListener(func(_ string, _ *linmot_client.DriveFaultError, _ bool) {
		escalationCalls++
	})
	t.Cleanup(removeEscalation)

	fault := &linmot_client.DriveFaultError{
		StatusWord: 0x0000,
		ErrorCode:  0x0042,
	}
	handleDetectedFault(context.Background(), ip, fault)

	assert.Equal(t, 0, mockClient.ackCalls, "budget exceeded should skip auto-recovery acknowledge")
	assert.Equal(t, 1, faultCalls)
	assert.Equal(t, 1, escalationCalls)
	snapshot := getFaultLifecycleSnapshot(ip)
	assert.Equal(t, FaultLifecycleStateEscalated, snapshot.State)
	assert.Equal(t, 0, snapshot.ConsecutiveRecoverFailures)
}

func TestHandleDetectedFault_PersistentNonFatalFaultRetriesThroughBudgetThenEscalates(t *testing.T) {
	ip := "10.0.0.16"
	resetFaultHandlingState(ip)

	mockClient := &faultTestClient{ackErr: errors.New("ack failed")}
	SetClientFactory(&faultTestFactory{client: mockClient})
	defer ResetClientFactory()

	faultCalls := 0
	escalationCalls := 0
	removeFault := AddFaultListener(func(_ string, _ error) { faultCalls++ })
	defer removeFault()
	removeEscalation := AddEscalationListener(func(_ string, _ *linmot_client.DriveFaultError, _ bool) {
		escalationCalls++
	})
	t.Cleanup(removeEscalation)

	fault := &linmot_client.DriveFaultError{
		StatusWord: 0x0000,
		ErrorCode:  0x0042,
	}

	for attempt := 1; attempt <= autoRecoveryMaxFaults; attempt++ {
		handleDetectedFault(context.Background(), ip, fault)

		snapshot := getFaultLifecycleSnapshot(ip)
		assert.Equal(t, FaultLifecycleStateRecovering, snapshot.State)
		assert.Equal(t, attempt, snapshot.ConsecutiveRecoverFailures)
		assert.Equal(t, attempt, mockClient.ackCalls)
		assert.Equal(t, 1, faultCalls, "persistent identical fault should only broadcast once while recovering")
		assert.Equal(t, 0, escalationCalls, "escalation should wait until budget is exhausted")
	}

	handleDetectedFault(context.Background(), ip, fault)

	snapshot := getFaultLifecycleSnapshot(ip)
	assert.Equal(t, autoRecoveryMaxFaults, mockClient.ackCalls, "budget exhaustion should not trigger an additional acknowledge")
	assert.Equal(t, 1, faultCalls, "fault notification should remain deduped on final escalation")
	assert.Equal(t, 1, escalationCalls, "final budget exhaustion should escalate once")
	assert.Equal(t, FaultLifecycleStateEscalated, snapshot.State)
	assert.Equal(t, autoRecoveryMaxFaults, snapshot.ConsecutiveRecoverFailures)
	assert.NotZero(t, snapshot.LastEscalationTime.UnixNano())

	handleDetectedFault(context.Background(), ip, fault)
	assert.Equal(t, 1, faultCalls, "persistent escalated fault should keep fault notification deduped")
	assert.Equal(t, 1, escalationCalls, "persistent escalated fault should keep escalation deduped")
}

func TestHandleDetectedFault_FatalEscalatesImmediately(t *testing.T) {
	ip := "10.0.0.10"
	resetFaultHandlingState(ip)

	mockClient := &faultTestClient{}
	SetClientFactory(&faultTestFactory{client: mockClient})
	defer ResetClientFactory()

	faultCalls := 0
	var gotEscalationFatal bool
	escalationCalls := 0
	removeFault := AddFaultListener(func(_ string, _ error) {
		faultCalls++
	})
	defer removeFault()
	removeEscalation := AddEscalationListener(func(_ string, _ *linmot_client.DriveFaultError, isFatal bool) {
		escalationCalls++
		gotEscalationFatal = isFatal
	})
	t.Cleanup(removeEscalation)

	fault := &linmot_client.DriveFaultError{
		StatusWord: 0x1000,
		ErrorCode:  0x0001,
	}
	handleDetectedFault(context.Background(), ip, fault)

	assert.Equal(t, 0, mockClient.ackCalls, "fatal faults should not attempt auto-recovery acknowledge")
	assert.Equal(t, 1, faultCalls)
	assert.Equal(t, 1, escalationCalls)
	assert.True(t, gotEscalationFatal)
	snapshot := getFaultLifecycleSnapshot(ip)
	assert.Equal(t, FaultLifecycleStateEscalated, snapshot.State)
	assert.Equal(t, FaultLevelFatal, snapshot.LastFaultLevel)
}

func TestHandleDetectedFault_WarningDoesNotEscalate(t *testing.T) {
	ip := "10.0.0.11"
	resetFaultHandlingState(ip)

	mockClient := &faultTestClient{}
	SetClientFactory(&faultTestFactory{client: mockClient})
	defer ResetClientFactory()

	faultCalls := 0
	escalationCalls := 0
	removeFault := AddFaultListener(func(_ string, _ error) { faultCalls++ })
	defer removeFault()
	removeEscalation := AddEscalationListener(func(_ string, _ *linmot_client.DriveFaultError, _ bool) {
		escalationCalls++
	})
	t.Cleanup(removeEscalation)

	warning := &linmot_client.DriveFaultError{
		StatusWord:  0x0000,
		ErrorCode:   0x0000,
		WarningWord: 0x0001,
	}
	handleDetectedFault(context.Background(), ip, warning)

	assert.Equal(t, 0, mockClient.ackCalls)
	assert.Equal(t, 0, faultCalls)
	assert.Equal(t, 0, escalationCalls)
	assert.Equal(t, FaultLifecycleStateHealthy, getFaultLifecycleSnapshot(ip).State)
}

func TestHandleDetectedFault_NetworkErrorBroadcastsDirectly(t *testing.T) {
	ip := "10.0.0.12"
	resetFaultHandlingState(ip)

	mockClient := &faultTestClient{}
	SetClientFactory(&faultTestFactory{client: mockClient})
	defer ResetClientFactory()

	var gotErr error
	faultCalls := 0
	removeFault := AddFaultListener(func(_ string, err error) {
		faultCalls++
		gotErr = err
	})
	defer removeFault()

	handleDetectedFault(context.Background(), ip, errors.New("dial timeout"))
	assert.Equal(t, 1, faultCalls)
	require.EqualError(t, gotErr, "dial timeout")
	assert.Equal(t, 0, mockClient.ackCalls)
}

func TestHandleDetectedFault_DedupesRepeatedEscalationForSamePersistentFault(t *testing.T) {
	ip := "10.0.0.13"
	resetFaultHandlingState(ip)

	mockClient := &faultTestClient{}
	SetClientFactory(&faultTestFactory{client: mockClient})
	defer ResetClientFactory()

	faultCalls := 0
	escalationCalls := 0
	removeFault := AddFaultListener(func(_ string, _ error) { faultCalls++ })
	defer removeFault()
	removeEscalation := AddEscalationListener(func(_ string, _ *linmot_client.DriveFaultError, _ bool) {
		escalationCalls++
	})
	t.Cleanup(removeEscalation)

	fatalFault := &linmot_client.DriveFaultError{
		StatusWord: 0x1000,
		ErrorCode:  0x0001,
	}

	handleDetectedFault(context.Background(), ip, fatalFault)
	firstSnapshot := getFaultLifecycleSnapshot(ip)
	time.Sleep(5 * time.Millisecond)
	handleDetectedFault(context.Background(), ip, fatalFault)
	secondSnapshot := getFaultLifecycleSnapshot(ip)

	assert.Equal(t, 1, faultCalls, "identical persistent fault should only notify once")
	assert.Equal(t, 1, escalationCalls, "identical persistent fault should only escalate once")
	assert.Equal(t, firstSnapshot.LastEscalationTime.UnixNano(), secondSnapshot.LastEscalationTime.UnixNano())
}

func TestFaultLifecycleStateTransitionsAreDeterministic(t *testing.T) {
	ipSuccess := "10.0.0.14"
	resetFaultHandlingState(ipSuccess)

	successClient := &faultTestClient{}
	SetClientFactory(&faultTestFactory{client: successClient})

	nonFatalFault := &linmot_client.DriveFaultError{
		StatusWord: 0x0000,
		ErrorCode:  0x0042,
	}
	handleDetectedFault(context.Background(), ipSuccess, nonFatalFault)
	successSnapshot := getFaultLifecycleSnapshot(ipSuccess)
	assert.Equal(t, FaultLifecycleStateHealthy, successSnapshot.State)
	assert.NotZero(t, successSnapshot.LastRecoverAttemptTime.UnixNano())
	assert.Zero(t, successSnapshot.ConsecutiveRecoverFailures)

	ipFailure := "10.0.0.15"
	resetFaultHandlingState(ipFailure)
	failureClient := &faultTestClient{ackErr: errors.New("ack failed")}
	SetClientFactory(&faultTestFactory{client: failureClient})
	defer ResetClientFactory()

	handleDetectedFault(context.Background(), ipFailure, nonFatalFault)
	failureSnapshot := getFaultLifecycleSnapshot(ipFailure)
	assert.Equal(t, FaultLifecycleStateRecovering, failureSnapshot.State)
	assert.NotZero(t, failureSnapshot.LastRecoverAttemptTime.UnixNano())
	assert.True(t, failureSnapshot.LastEscalationTime.IsZero())
	assert.Equal(t, 1, failureSnapshot.ConsecutiveRecoverFailures)
}

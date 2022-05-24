package node_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/stretchr/testify/require"

	"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite-hq/cli/integration"
)

func TestNodeTxBankSend(t *testing.T) {
	var (
		env           = envtest.New(t)
		path          = env.Scaffold("github.com/test/blog", "--address-prefix", testPrefix)
		servers       = env.RandomizeServerPorts(path, "")
		rndWorkdir    = t.TempDir() // To make sure we can run these commands from anywhere
		accKeyringDir = t.TempDir()
	)

	env.SetKeyringBackend(path, "", keyring.BackendTest)
	env.SetConfigMnemonic(path, "", "alice", aliceMnemonic)
	env.SetConfigMnemonic(path, "", "bob", bobMnemonic)

	var (
		ctx, cancel = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
	)

	go func() {
		defer cancel()
		isBackendAliveErr := env.IsAppServed(ctx, servers)
		require.NoError(t, isBackendAliveErr, "app cannot get online in time")

		// error "account doesn't have any balances" occurs if a sleep is not included
		time.Sleep(time.Second * 1)

		env.Must(env.Exec("import alice",
			step.NewSteps(step.New(
				step.Exec(envtest.IgniteApp, "account", "import", "alice", "--keyring-dir", accKeyringDir, "--non-interactive", "--secret", aliceMnemonic),
			)),
		))
		env.Must(env.Exec("import bob",
			step.NewSteps(step.New(
				step.Exec(envtest.IgniteApp, "account", "import", "bob", "--keyring-dir", accKeyringDir, "--non-interactive", "--secret", bobMnemonic),
			)),
		))

		env.Must(env.Exec("send 100token from alice to bob",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"tx",
					"bank",
					"send",
					"alice",
					"bob",
					"100token",
					"--from",
					"alice",
					"--rpc",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(rndWorkdir),
			)),
		))

		env.Must(env.Exec("send 2stake from bob to alice using addresses",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"tx",
					"bank",
					"send",
					bobAddress,
					aliceAddress,
					"2stake",
					"--from",
					"bob",
					"--rpc",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(rndWorkdir),
			)),
		))

		env.Must(env.Exec("send 5token from alice to bob using a combination of address and account",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"tx",
					"bank",
					"send",
					"alice",
					bobAddress,
					"5token",
					"--from",
					"alice",
					"--rpc",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(rndWorkdir),
			)),
		))

		time.Sleep(time.Second * 1)

		var aliceBalanceCheckBuffer = &bytes.Buffer{}
		env.Must(env.Exec("query bank balances for alice",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"alice",
					"--rpc",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(rndWorkdir),
			)),
			envtest.ExecStdout(aliceBalanceCheckBuffer),
		))
		require.True(t, strings.Contains(aliceBalanceCheckBuffer.String(), `Amount 		Denom 	
100000002 	stake 	
19895 		token`))

		var bobBalanceCheckBuffer = &bytes.Buffer{}
		env.Must(env.Exec("query bank balances for bob",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"bob",
					"--rpc",
					"http://"+servers.RPC,
					"--keyring-dir",
					accKeyringDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(rndWorkdir),
			)),
			envtest.ExecStdout(bobBalanceCheckBuffer),
		))
		require.True(t, strings.Contains(bobBalanceCheckBuffer.String(), `Amount 		Denom 	
99999998 	stake 	
10105 		token`))

	}()
	env.Must(env.Serve("should serve with Stargate version", path, "", "", envtest.ExecCtx(ctx)))
}

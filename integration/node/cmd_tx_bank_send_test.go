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
	"github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
	envtest "github.com/ignite-hq/cli/integration"
)

func TestNodeTxBankSend(t *testing.T) {
	var (
		env     = envtest.New(t)
		path    = env.Scaffold("github.com/test/blog", "--address-prefix", testPrefix)
		servers = env.RandomizeServerPorts(path, "")
		homeDir = env.SetRandomHomeConfig(path, "")
	)

	env.SetKeyringBackend(keyring.BackendTest, path, "")

	var (
		ctx, cancel = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
	)

	go func() {
		defer cancel()
		isBackendAliveErr := env.IsAppServed(ctx, servers)
		require.NoError(t, isBackendAliveErr, "app cannot get online in time")

		// error "account doesn't have any balances" occurs if a sleep is not included
		time.Sleep(time.Second * 1)

		accounts := env.AccountsInKeyring(homeDir, cosmosaccount.KeyringTest, testPrefix)
		var alice envtest.Account
		var bob envtest.Account
		for _, acc := range accounts {
			if acc.Name == "alice" {
				alice = acc
			}
			if acc.Name == "bob" {
				bob = acc
				break
			}
		}

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
					"--home",
					homeDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(path),
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
					bob.Address,
					alice.Address,
					"2stake",
					"--from",
					"bob",
					"--rpc",
					"http://"+servers.RPC,
					"--home",
					homeDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(path),
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
					bob.Address,
					"5token",
					"--from",
					"alice",
					"--rpc",
					"http://"+servers.RPC,
					"--home",
					homeDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(path),
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
					"--home",
					homeDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(path),
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
					"--home",
					homeDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(path),
			)),
			envtest.ExecStdout(bobBalanceCheckBuffer),
		))
		require.True(t, strings.Contains(bobBalanceCheckBuffer.String(), `Amount 		Denom 	
99999998 	stake 	
10105 		token`))

	}()
	env.Must(env.Serve("should serve with Stargate version", path, "", "", envtest.ExecCtx(ctx)))
}

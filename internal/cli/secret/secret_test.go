package secret_test

import (
	"testing"
	"time"

	"github.com/holos-run/holos/internal/cli"
	"github.com/holos-run/holos/internal/holos"
	"github.com/rogpeppe/go-internal/testscript"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

const clientsetKey = "clientset"

var secret = v1.Secret{
	TypeMeta: metav1.TypeMeta{
		Kind:       "Secret",
		APIVersion: "v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "k2-talos",
		Namespace: "secrets",
		Labels: map[string]string{
			"holos.run/owner.name":  "jeff",
			"holos.run/secret.name": "k2-talos",
		},
		CreationTimestamp: metav1.Time{
			Time: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	},
	Data: map[string][]byte{
		"secrets.yaml": []byte("content: secret\n"),
	},
	Type: "Opaque",
}

var loginSecret = v1.Secret{
	TypeMeta: metav1.TypeMeta{
		Kind:       "Secret",
		APIVersion: "v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "zitadel-admin-d7fgbgbfbt",
		Namespace: "secrets",
		Labels: map[string]string{
			"holos.run/owner.name":  "jeff",
			"holos.run/secret.name": "zitadel-admin",
		},
		CreationTimestamp: metav1.Time{
			Time: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	},
	Data: map[string][]byte{
		"url":      []byte("https://login.example.com"),
		"username": []byte("zitadel-admin@zitadel.login.example.com"),
		"password": []byte("Password1!"),
	},
	Type: "Opaque",
}

// cmdHolos executes the holos root command with a kubernetes.Interface that
// persists for the duration of the testscript. holos is NOT executed in a
// subprocess, the current working directory is not and should not be changed.
// Take care to read and write to $WORK in the test scripts using flags.
func cmdHolos(ts *testscript.TestScript, neg bool, args []string) {
	clientset, ok := ts.Value(clientsetKey).(kubernetes.Interface)
	if clientset == nil || !ok {
		ts.Fatalf("missing kubernetes.Interface")
	}

	cfg := holos.New(
		holos.ProvisionerClientset(clientset),
		holos.Stdout(ts.Stdout()),
		holos.Stderr(ts.Stderr()),
	)

	cmd := cli.New(cfg)
	cmd.SetArgs(args)
	err := cmd.Execute()

	if neg {
		if err == nil {
			ts.Fatalf("\nwant: error\nhave: %v", err)
		} else {
			cli.HandleError(cmd.Context(), err, cfg)
		}
	} else {
		ts.Check(err)
	}
}

func TestSecrets(t *testing.T) {
	// Add TestWork: true to the Params to keep the $WORK directory around.
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Setup: func(env *testscript.Env) error {
			env.Values[clientsetKey] = fake.NewSimpleClientset(&secret, &loginSecret)
			return nil
		},
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"holos": cmdHolos,
		},
	})
}

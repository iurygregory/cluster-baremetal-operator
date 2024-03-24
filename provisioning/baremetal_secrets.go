package provisioning

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	coreclientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
)

const (
	ironicUsernameKey        = "username"
	ironicPasswordKey        = "password"
	ironicHtpasswdKey        = "htpasswd"
	ironicConfigKey          = "auth-config"
	ironicSecretName         = "metal3-ironic-password"
	ironicUsername           = "ironic-user"
	inspectorSecretName      = "metal3-ironic-inspector-password"
	inspectorUsername        = "inspector-user"
	tlsSecretName            = "metal3-ironic-tls" // #nosec
	openshiftConfigSecretKey = ".dockerconfigjson" // #nosec
	// NOTE(dtantsur): this is kept here to be able to remove the old
	// secret when a Provisioning is removed.
	ironicrpcSecretName = "metal3-ironic-rpc-password" // #nosec
	baremetalSecretName = "metal3-mariadb-password"    // #nosec

	// OpenshiftConfigNamespace holds the name of the openshift-config namespace.
	OpenshiftConfigNamespace = "openshift-config"
	// PullSecretName holds the name of the pull-secret in openshift-config and openshift-machine-config.
	PullSecretName = "pull-secret"
)

type shouldUpdateDataFn func(existing *corev1.Secret) (bool, error)

func doNotUpdateData(existing *corev1.Secret) (bool, error) {
	return false, nil
}

// applySecret merges objectmeta, applies data if the secret does not exist or shouldUpdateDataFn returns false.
func applySecret(client coreclientv1.SecretsGetter, recorder events.Recorder, requiredInput *corev1.Secret, shouldUpdateData shouldUpdateDataFn) error {
	needsApply := false
	existing, err := client.Secrets(requiredInput.Namespace).Get(context.TODO(), requiredInput.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		err = nil
		needsApply = true
	} else if err != nil {
		return err
	} else {
		// Allow the caller to decide whether update.
		needsApply, err = shouldUpdateData(existing)
		if err != nil {
			return err
		}
	}

	if needsApply {
		_, _, err = resourceapply.ApplySecret(context.TODO(), client, recorder, requiredInput)
	}

	return err
}

func createIronicSecret(info *ProvisioningInfo, name string, username string, configSection string) error {
	password, err := generateRandomPassword()
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 5) // Use same cost as htpasswd default
	if err != nil {
		return err
	}
	// Change hash version from $2a$ to $2y$, as generated by htpasswd.
	// These are equivalent for our purposes.
	// Some background information about this : https://en.wikipedia.org/wiki/Bcrypt#Versioning_history
	// There was a bug 9 years ago in PHP's implementation of 2a, so they decided to call the fixed version 2y.
	// httpd decided to adopt this (if it sees 2a it uses elaborate heuristic workarounds to mitigate against the bug,
	// but 2y is assumed to not need them), but everyone else (including go) was just decided to not implement the bug in 2a.
	// The bug only affects passwords containing characters with the high bit set, i.e. not ASCII passwords generated here.

	// Anyway, Ironic implemented their own basic auth verification and originally hard-coded 2y because that's what
	// htpasswd produces (see https://review.opendev.org/738718). It is better to keep this as one day we may move the auth
	// to httpd and this would prevent triggering the workarounds.
	hash[2] = 'y'

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: info.Namespace,
		},
		StringData: map[string]string{
			ironicUsernameKey: username,
			ironicPasswordKey: password,
			ironicHtpasswdKey: fmt.Sprintf("%s:%s", username, hash),
			ironicConfigKey: fmt.Sprintf(`[%s]
auth_type = http_basic
username = %s
password = %s
`,
				configSection, username, password),
		},
	}

	if err := controllerutil.SetControllerReference(info.ProvConfig, secret, info.Scheme); err != nil {
		return err
	}

	return applySecret(info.Client.CoreV1(), info.EventRecorder, secret, doNotUpdateData)
}

// createRegistryPullSecret creates a copy of the pull-secret in the
// openshift-config namespace for use with LocalObjectReference
func createRegistryPullSecret(info *ProvisioningInfo) error {
	client := info.Client.CoreV1()
	openshiftConfigSecret, err := client.Secrets(OpenshiftConfigNamespace).Get(context.TODO(), PullSecretName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("could not get secret %s/%s, err: %w", OpenshiftConfigNamespace, PullSecretName, err)
	}
	openshiftConfigSecretKeyData, ok := openshiftConfigSecret.Data[openshiftConfigSecretKey]
	if !ok {
		return fmt.Errorf("could not find key %q in secret %s/%s", openshiftConfigSecretKey, OpenshiftConfigNamespace, PullSecretName)
	}

	machineAPINamespace := info.Namespace
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PullSecretName,
			Namespace: machineAPINamespace,
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			// The openshift-machine-api/pull-secret .dockerconfigjson field should be double encoded due to PR
			// https://github.com/openshift/cluster-baremetal-operator/pull/184
			openshiftConfigSecretKey: base64.StdEncoding.EncodeToString(openshiftConfigSecretKeyData),
		},
	}

	if err := controllerutil.SetControllerReference(info.ProvConfig, secret, info.Scheme); err != nil {
		return err
	}
	_, changed, err := resourceapply.ApplySecret(context.TODO(), client, info.EventRecorder, secret)
	if changed {
		reportRegistryPullSecretReconcile()
	}
	return err
}

// reportRegistryPullSecretReconcile is used for unit testing, to report that the reconciler was triggered.
var reportRegistryPullSecretReconcile = func() {}

func EnsureAllSecrets(info *ProvisioningInfo) (bool, error) {
	// Create a Secret for the Ironic Password
	if err := createIronicSecret(info, ironicSecretName, ironicUsername, "ironic"); err != nil {
		return false, errors.Wrap(err, "failed to create Ironic password")
	}
	// Create a Secret for the Ironic Inspector Password
	if err := createIronicSecret(info, inspectorSecretName, inspectorUsername, "inspector"); err != nil {
		return false, errors.Wrap(err, "failed to create Inspector password")
	}
	// Generate/update TLS certificate
	if err := createOrUpdateTlsSecret(info); err != nil {
		return false, errors.Wrap(err, "failed to create TLS certificate")
	}
	// Create a Secret for the Registry Pull Secret
	if err := createRegistryPullSecret(info); err != nil {
		return false, errors.Wrap(err, "failed to create Registry pull secret")
	}
	return false, nil // ApplySecret does not use Generation, so just return false for updated
}

func DeleteAllSecrets(info *ProvisioningInfo) error {
	var secretErrors []error
	for _, sn := range []string{baremetalSecretName, ironicSecretName, inspectorSecretName, ironicrpcSecretName, tlsSecretName, PullSecretName} {
		if err := client.IgnoreNotFound(info.Client.CoreV1().Secrets(info.Namespace).Delete(context.Background(), sn, metav1.DeleteOptions{})); err != nil {
			secretErrors = append(secretErrors, err)
		}
	}
	return utilerrors.NewAggregate(secretErrors)
}

// createOrUpdateTlsSecret creates a Secret for the Ironic and Inspector TLS.
// It updates the secret if the existing certificate is close to expiration.
func createOrUpdateTlsSecret(info *ProvisioningInfo) error {
	cert, err := generateTlsCertificate(info.ProvConfig.Spec.ProvisioningIP)
	if err != nil {
		return err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tlsSecretName,
			Namespace: info.Namespace,
		},
		Data: map[string][]byte{
			corev1.TLSCertKey:       cert.certificate,
			corev1.TLSPrivateKeyKey: cert.privateKey,
		},
	}

	if err := controllerutil.SetControllerReference(info.ProvConfig, secret, info.Scheme); err != nil {
		return err
	}

	return applySecret(info.Client.CoreV1(), info.EventRecorder, secret, func(existing *corev1.Secret) (bool, error) {
		expired, err := isTlsCertificateExpired(existing.Data[corev1.TLSCertKey])
		if err != nil {
			return false, err
		}
		return expired, nil
	})
}

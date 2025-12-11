# Deployment Helpers (Go Reference)

Go-based reference implementations for automated deployment tasks.

## Storage & Key Preparation

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// PrepareStorageAndKeys sets up hostPath and Ed25519 signing keys
func PrepareStorageAndKeys(basePath string) error {
	// Create directories
	storageDir := filepath.Join(basePath, "n-audit-data")
	signingDir := filepath.Join(storageDir, "signing")
	
	if err := os.MkdirAll(signingDir, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	
	// Generate Ed25519 key (non-interactive)
	keyPath := filepath.Join(signingDir, "id_ed25519")
	cmd := exec.Command("ssh-keygen",
		"-t", "ed25519",
		"-N", "",
		"-f", keyPath,
		"-C", "n-audit-sentinel@localhost",
	)
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ssh-keygen: %w", err)
	}
	
	// Set secure permissions
	if err := os.Chmod(keyPath, 0600); err != nil {
		return fmt.Errorf("chmod private key: %w", err)
	}
	if err := os.Chmod(keyPath+".pub", 0644); err != nil {
		return fmt.Errorf("chmod public key: %w", err)
	}
	
	fmt.Printf("✓ Storage prepared at: %s\n", storageDir)
	fmt.Printf("✓ Keys generated in: %s\n", signingDir)
	return nil
}
```

## RBAC Creation (Go Client)

```go
package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ApplyRBAC creates ServiceAccount and ClusterRole bindings
func ApplyRBAC(kubeClient kubernetes.Interface) error {
	ctx := context.Background()
	
	// Create ServiceAccount
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "n-audit-sentinel",
			Namespace: "default",
		},
	}
	
	if _, err := kubeClient.CoreV1().ServiceAccounts("default").Create(ctx, sa, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("create ServiceAccount: %w", err)
		}
	}
	
	// Create ClusterRole for Cilium policies
	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{Name: "n-audit-cilium-policy"},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"cilium.io"},
				Resources: []string{"ciliumnetworkpolicies"},
				Verbs:     []string{"get", "list", "create", "delete", "update", "patch"},
			},
		},
	}
	
	if _, err := kubeClient.RbacV1().ClusterRoles().Create(ctx, role, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("create ClusterRole: %w", err)
		}
	}
	
	// Create ClusterRoleBinding
	binding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{Name: "n-audit-cilium-policy"},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "n-audit-cilium-policy",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "n-audit-sentinel",
				Namespace: "default",
			},
		},
	}
	
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(ctx, binding, metav1.CreateOptions{}); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("create ClusterRoleBinding: %w", err)
		}
	}
	
	fmt.Println("✓ ServiceAccount and RBAC configured")
	return nil
}
```

## Pod Deployment (Go Client)

```go
package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DeployPod creates the N-Audit Sentinel pod
func DeployPod(kubeClient kubernetes.Interface, imageName, imageTag string) error {
	ctx := context.Background()
	
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "n-audit-sentinel",
			Namespace: "default",
			Labels: map[string]string{
				"app": "n-audit-sentinel",
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "n-audit-sentinel",
			SecurityContext: &corev1.PodSecurityContext{
				RunAsUser: func(i int64) *int64 { return &i }(0),
			},
			Containers: []corev1.Container{
				{
					Name:            "sentinel",
					Image:           imageName + ":" + imageTag,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Stdin:           true,
					TTY:             true,
					Env: []corev1.EnvVar{
						{
							Name:  "SSH_SIGN_KEY_PATH",
							Value: "/var/lib/n-audit/signing/id_ed25519",
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "data",
							MountPath: "/var/lib/n-audit",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "data",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/mnt/n-audit-data",
							Type: func(s corev1.HostPathType) *corev1.HostPathType { return &s }(corev1.HostPathDirectoryOrCreate),
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyAlways,
		},
	}
	
	if _, err := kubeClient.CoreV1().Pods("default").Create(ctx, pod, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create pod: %w", err)
	}
	
	fmt.Println("✓ Pod deployed successfully")
	return nil
}
```

## Health Check (Go Reference)

```go
package main

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// WaitForPod waits for pod to be running
func WaitForPod(kubeClient kubernetes.Interface, name, namespace string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	for {
		pod, err := kubeClient.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("get pod: %w", err)
		}
		
		if pod.Status.Phase == "Running" {
			fmt.Printf("✓ Pod %s is running\n", name)
			return nil
		}
		
		if pod.Status.Phase == "Failed" || pod.Status.Phase == "Unknown" {
			return fmt.Errorf("pod failed: %s", pod.Status.Phase)
		}
		
		fmt.Printf("Waiting for pod... (phase: %s)\n", pod.Status.Phase)
		time.Sleep(2 * time.Second)
	}
}
```

## Image Build & Load (Go Reference)

```go
package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// BuildAndLoadImage builds Docker image and loads into K3s
func BuildAndLoadImage(imageName, imageTag, dockerfilePath string) error {
	fullImage := imageName + ":" + imageTag
	
	// Build image
	fmt.Printf("Building image: %s\n", fullImage)
	buildCmd := exec.Command("docker", "build",
		"-t", fullImage,
		"-f", dockerfilePath,
		".",
	)
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("build: %w", err)
	}
	
	// Save and import into K3s
	fmt.Println("Importing into K3s...")
	saveCmd := exec.Command("docker", "save", fullImage)
	importCmd := exec.Command("sudo", "k3s", "ctr", "images", "import", "-")
	
	pipe, err := saveCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("pipe: %w", err)
	}
	
	importCmd.Stdin = pipe
	
	if err := saveCmd.Start(); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	
	if err := importCmd.Run(); err != nil {
		return fmt.Errorf("import: %w", err)
	}
	
	if err := saveCmd.Wait(); err != nil {
		return fmt.Errorf("save wait: %w", err)
	}
	
	fmt.Printf("✓ Image loaded: %s\n", fullImage)
	return nil
}
```

## Verification Workflow (Go Reference)

```go
package main

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// VerifyDeployment checks all deployment components
func VerifyDeployment(kubeClient kubernetes.Interface) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	checks := []struct {
		name string
		fn   func() error
	}{
		{
			name: "ServiceAccount exists",
			fn: func() error {
				_, err := kubeClient.CoreV1().ServiceAccounts("default").Get(ctx, "n-audit-sentinel", metav1.GetOptions{})
				return err
			},
		},
		{
			name: "Pod is running",
			fn: func() error {
				pod, err := kubeClient.CoreV1().Pods("default").Get(ctx, "n-audit-sentinel", metav1.GetOptions{})
				if err != nil {
					return err
				}
				if pod.Status.Phase != "Running" {
					return fmt.Errorf("pod not running: %s", pod.Status.Phase)
				}
				return nil
			},
		},
		{
			name: "RBAC ClusterRole exists",
			fn: func() error {
				_, err := kubeClient.RbacV1().ClusterRoles().Get(ctx, "n-audit-cilium-policy", metav1.GetOptions{})
				return err
			},
		},
	}
	
	for _, check := range checks {
		if err := check.fn(); err != nil {
			fmt.Printf("✗ %s: %v\n", check.name, err)
		} else {
			fmt.Printf("✓ %s\n", check.name)
		}
	}
	
	return nil
}
```

## Integration Example

```go
// Complete deployment workflow
func DeploymentWorkflow(kubeClient kubernetes.Interface, imageName, imageTag string) error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Prepare storage", func() error { return PrepareStorageAndKeys("/mnt") }},
		{"Apply RBAC", func() error { return ApplyRBAC(kubeClient) }},
		{"Build image", func() error { return BuildAndLoadImage(imageName, imageTag, "Dockerfile") }},
		{"Deploy pod", func() error { return DeployPod(kubeClient, imageName, imageTag) }},
		{"Wait for ready", func() error { return WaitForPod(kubeClient, "n-audit-sentinel", "default", 30*time.Second) }},
		{"Verify", func() error { return VerifyDeployment(kubeClient) }},
	}
	
	for _, step := range steps {
		fmt.Printf("\n[%s]\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("%s failed: %w", step.name, err)
		}
	}
	
	fmt.Println("\n✅ Deployment complete!")
	return nil
}
```

## See Also

- [DEPLOYMENT.md](../DEPLOYMENT.md) — Full deployment guide
- [TESTING_AND_VERIFICATION.md](TESTING_AND_VERIFICATION.md) — Testing procedures
- [README.md](../README.md) — Architecture overview

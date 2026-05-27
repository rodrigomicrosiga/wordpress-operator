package factory

import (
	"fmt"

	"github.com/rodrigomicrosiga/wordpress-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// commonLabels gera as labels padrão para amarrar todos os recursos ao WordpressSite
func commonLabels(site *v1alpha1.WordpressSite, component string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "wordpress",
		"app.kubernetes.io/instance":   site.Name,
		"app.kubernetes.io/component":  component,
		"app.kubernetes.io/managed-by": "wordpress-operator",
	}
}

// --- BANCO DE DADOS (MYSQL) ---

func BuildMySQLSecret(site *v1alpha1.WordpressSite, password string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-mysql-secret", site.Name),
			Namespace: site.Namespace,
			Labels:    commonLabels(site, "database"),
		},
		StringData: map[string]string{
			"mysql-password":      password,
			"mysql-root-password": password, // Simplificação: mesma senha para root
		},
	}
}

func BuildMySQLStatefulSet(site *v1alpha1.WordpressSite) *appsv1.StatefulSet {
	labels := commonLabels(site, "database")
	replicas := int32(1)

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-mysql", site.Name),
			Namespace: site.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "mysql",
						Image: site.Spec.Database.Image,
						Env: []corev1.EnvVar{
							{Name: "MYSQL_DATABASE", Value: site.Spec.Database.Name},
							{Name: "MYSQL_USER", Value: site.Spec.Database.User},
							{
								Name: "MYSQL_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: fmt.Sprintf("%s-mysql-secret", site.Name),
										},
										Key: "mysql-password",
									},
								},
							},
							{
								Name: "MYSQL_ROOT_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: fmt.Sprintf("%s-mysql-secret", site.Name),
										},
										Key: "mysql-root-password",
									},
								},
							},
						},
						Ports: []corev1.ContainerPort{{ContainerPort: 3306}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "mysql-data",
							MountPath: "/var/lib/mysql",
						}},
					}},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: metav1.ObjectMeta{Name: "mysql-data"},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					Resources: corev1.VolumeResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: resource.MustParse(site.Spec.Database.StorageSize),
						},
					},
				},
			}},
		},
	}
}

func BuildMySQLService(site *v1alpha1.WordpressSite) *corev1.Service {
	labels := commonLabels(site, "database")
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-mysql", site.Name),
			Namespace: site.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector:  labels,
			Ports:     []corev1.ServicePort{{Port: 3306, TargetPort: intstr.FromInt(3306)}},
			ClusterIP: "None", // Headless service para StatefulSet
		},
	}
}

// --- APLICAÇÃO (WORDPRESS) ---

func BuildWordpressConfigMap(site *v1alpha1.WordpressSite) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-wp-config", site.Name),
			Namespace: site.Namespace,
			Labels:    commonLabels(site, "application"),
		},
		Data: map[string]string{
			"WORDPRESS_DB_HOST": fmt.Sprintf("%s-mysql", site.Name),
			"WORDPRESS_DB_NAME": site.Spec.Database.Name,
			"WORDPRESS_DB_USER": site.Spec.Database.User,
		},
	}
}

func BuildWordpressPVC(site *v1alpha1.WordpressSite) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-wp-content", site.Name),
			Namespace: site.Namespace,
			Labels:    commonLabels(site, "application"),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(site.Spec.Wordpress.StorageSize),
				},
			},
		},
	}
}

func BuildWordpressDeployment(site *v1alpha1.WordpressSite) *appsv1.Deployment {
	labels := commonLabels(site, "application")
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      site.Name,
			Namespace: site.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &site.Spec.Wordpress.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "wordpress",
						Image: site.Spec.Wordpress.Image,
						EnvFrom: []corev1.EnvFromSource{
							{ConfigMapRef: &corev1.ConfigMapEnvSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: fmt.Sprintf("%s-wp-config", site.Name),
								},
							}},
						},
						Env: []corev1.EnvVar{
							{
								Name: "WORDPRESS_DB_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: fmt.Sprintf("%s-mysql-secret", site.Name),
										},
										Key: "mysql-password",
									},
								},
							},
						},
						Ports: []corev1.ContainerPort{{ContainerPort: 80}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "wp-content",
							MountPath: "/var/www/html/wp-content",
						}},
					}},
					Volumes: []corev1.Volume{{
						Name: "wp-content",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: fmt.Sprintf("%s-wp-content", site.Name),
							},
						},
					}},
				},
			},
		},
	}
}

func BuildWordpressService(site *v1alpha1.WordpressSite) *corev1.Service {
	labels := commonLabels(site, "application")
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      site.Name,
			Namespace: site.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports:    []corev1.ServicePort{{Port: 80, TargetPort: intstr.FromInt(80)}},
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
}

func BuildWordpressIngress(site *v1alpha1.WordpressSite) *networkingv1.Ingress {
	labels := commonLabels(site, "application")
	pathType := networkingv1.PathTypePrefix

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      site.Name,
			Namespace: site.Namespace,
			Labels:    labels,
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{{
				Host: site.Spec.Domain,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{{
							Path:     "/",
							PathType: &pathType,
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: site.Name,
									Port: networkingv1.ServiceBackendPort{Number: 80},
								},
							},
						}},
					},
				},
			}},
		},
	}

	if site.Spec.IngressClassName != "" {
		ingress.Spec.IngressClassName = &site.Spec.IngressClassName
	}

	return ingress
}

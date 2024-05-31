package client

import (
	"context"
	"fmt"
	"log"
	"sort"

	"gopkg.in/yaml.v2"
)

type ObjectMeta struct {
	Name      string `yaml:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Namespace string `yaml:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
}

type SecretStoreRef struct {
	Name string `yaml:"name"`
	Kind string `yaml:"kind,omitempty"`
}

type ExternalSecretTargetTemplate struct {
	Type string `yaml:"type,omitempty"`
}

type ExternalSecretTarget struct {
	Name     string                       `yaml:"name,omitempty"`
	Template ExternalSecretTargetTemplate `yaml:"template,omitempty"`
}

type ExternalSecretData struct {
	SecretKey string                      `yaml:"secretKey"`
	RemoteRef ExternalSecretDataRemoteRef `yaml:"remoteRef"`
}

type ExternalSecretDataRemoteRef struct {
	Key      string `yaml:"key"`
	Property string `yaml:"property,omitempty"`
}

type ExternalSecretSpec struct {
	SecretStoreRef  SecretStoreRef       `yaml:"secretStoreRef,omitempty"`
	Target          ExternalSecretTarget `yaml:"target,omitempty"`
	RefreshInterval string               `yaml:"refreshInterval,omitempty"`
	Data            []ExternalSecretData `yaml:"data,omitempty"`
}

type ExternalSecret struct {
	Kind       string     `yaml:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	APIVersion string     `yaml:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
	Metadata   ObjectMeta `yaml:"metadata,omitempty"`

	Spec ExternalSecretSpec `yaml:"spec,omitempty"`
}

type SecretStoreSpec struct {
	Provider        *SecretStoreProvider `yaml:"provider"`
	RefreshInterval int                  `yaml:"refreshInterval,omitempty"`
}

type SecretKeySelector struct {
	Name string `yaml:"name,omitempty"`
	Key  string `yaml:"key,omitempty"`
}

type VaultAuth struct {
	TokenSecretRef *SecretKeySelector `yaml:"tokenSecretRef,omitempty"`
}

type VaultProvider struct {
	Auth    VaultAuth `yaml:"auth"`
	Server  string    `yaml:"server"`
	Path    string    `yaml:"path,omitempty"`
	Version string    `yaml:"version"`
}

type SecretStoreProvider struct {
	Vault *VaultProvider `yaml:"vault,omitempty"`
}

type SecretStore struct {
	Kind       string     `yaml:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	APIVersion string     `yaml:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
	Metadata   ObjectMeta `yaml:"metadata,omitempty"`

	Spec SecretStoreSpec `yaml:"spec,omitempty"`
}

func (k *K8SClient) CreateOrUpdateSecretStore(ctx context.Context, secretStoreName, namespace, vaultAddr,
	tokenSecretName, tokenSecretKey string) (err error) {
	secretStore := SecretStore{
		APIVersion: "external-secrets.io/v1beta1",
		Kind:       "SecretStore",
		Metadata: ObjectMeta{
			Name:      secretStoreName,
			Namespace: namespace,
		},
		Spec: SecretStoreSpec{
			RefreshInterval: 10,
			Provider: &SecretStoreProvider{
				Vault: &VaultProvider{
					Server:  vaultAddr,
					Path:    "secret",
					Version: "v2",
					Auth: VaultAuth{
						TokenSecretRef: &SecretKeySelector{
							Key:  tokenSecretKey,
							Name: tokenSecretName,
						},
					},
				},
			},
		},
	}

	secretStoreData, err := yaml.Marshal(&secretStore)
	if err != nil {
		return
	}
	_, _, err = k.DynamicClient.CreateResource(ctx, []byte(secretStoreData))
	if err != nil {
		err = fmt.Errorf("failed to create cluter vault token secret %s/%s, %v", namespace, secretStoreName, err)
		return
	}
	return
}

// func (k *K8SClient) CreateOrUpdateExternalSecret(ctx context.Context, externalSecretName, namespace,
// 	secretStoreRefName, secretName, secretType string, vaultKeyPathdata map[string]string) (err error) {
// 	secretKeysData := []ExternalSecretData{}
// 	for key, path := range vaultKeyPathdata {
// 		secretKeyData := ExternalSecretData{
// 			SecretKey: key,
// 			RemoteRef: ExternalSecretDataRemoteRef{
// 				Key:      path,
// 				Property: key,
// 			},
// 		}
// 		secretKeysData = append(secretKeysData, secretKeyData)
// 	}
// 	externalSecret := ExternalSecret{
// 		APIVersion: "external-secrets.io/v1beta1",
// 		Kind:       "ExternalSecret",
// 		Metadata: ObjectMeta{
// 			Name:      externalSecretName,
// 			Namespace: namespace,
// 		},
// 		Spec: ExternalSecretSpec{
// 			RefreshInterval: "10s",
// 			Target: ExternalSecretTarget{
// 				Name:     secretName,
// 				Template: ExternalSecretTargetTemplate{Type: secretType}},
// 			SecretStoreRef: SecretStoreRef{
// 				Name: secretStoreRefName,
// 				Kind: "SecretStore",
// 			},
// 			Data: secretKeysData,
// 		},
// 	}

// 	externalSecretData, err := yaml.Marshal(&externalSecret)
// 	if err != nil {
// 		return
// 	}

// 	_, _, err = k.DynamicClient.CreateResource(ctx, []byte(externalSecretData))
// 	if err != nil {
// 		err = fmt.Errorf("failed to create vault external secret %s/%s, %v", namespace, externalSecretName, err)
// 		return
// 	}
// 	return
// }

// func (k *K8SClient) CreateOrUpdateExternalSecret(ctx context.Context, externalSecretName, namespace,
// 	secretStoreRefName, secretName, secretType string, vaultKeyPathdata, properties map[string]string) (err error) {
// 	secretKeysData := []ExternalSecretData{}
// 	log.Println("Vault Path data", vaultKeyPathdata)
// 	log.Println("properties", properties)
// 	for key, path := range vaultKeyPathdata {
// 		property := properties[key]
// 		if property == "" {
// 			property = key // Default to key if property is not specified
// 		}
// 		log.Println("Property", property)
// 		secretKeyData := ExternalSecretData{
// 			SecretKey: key,
// 			RemoteRef: ExternalSecretDataRemoteRef{
// 				Key:      path,
// 				Property: property,
// 			},
// 		}
// 		secretKeysData = append(secretKeysData, secretKeyData)
// 		log.Println("secret keys data", secretKeysData)
// 	}
// 	externalSecret := ExternalSecret{
// 		APIVersion: "external-secrets.io/v1beta1",
// 		Kind:       "ExternalSecret",
// 		Metadata: ObjectMeta{
// 			Name:      externalSecretName,
// 			Namespace: namespace,
// 		},
// 		Spec: ExternalSecretSpec{
// 			RefreshInterval: "10s",
// 			Target: ExternalSecretTarget{
// 				Name:     secretName,
// 				Template: ExternalSecretTargetTemplate{Type: secretType}},
// 			SecretStoreRef: SecretStoreRef{
// 				Name: secretStoreRefName,
// 				Kind: "SecretStore",
// 			},
// 			Data: secretKeysData,
// 		},
// 	}

// 	externalSecretData, err := yaml.Marshal(&externalSecret)
// 	if err != nil {
// 		return
// 	}

// 	_, _, err = k.DynamicClient.CreateResource(ctx, []byte(externalSecretData))
// 	if err != nil {
// 		err = fmt.Errorf("failed to create vault external secret %s/%s, %v", namespace, externalSecretName, err)
// 		return
// 	}
// 	return
// }

// func (k *K8SClient) CreateOrUpdateExternalSecret(ctx context.Context, externalSecretName, namespace,
// 	secretStoreRefName, secretName, secretType string, vaultKeyPathdata, secretProperties map[string]string) (err error) {
// 	secretKeysData := []ExternalSecretData{}
// 	for key, path := range vaultKeyPathdata {
// 		property := secretProperties[key]
// 		secretKeyData := ExternalSecretData{
// 			SecretKey: key,
// 			RemoteRef: ExternalSecretDataRemoteRef{
// 				Key:      path,
// 				Property: property,
// 			},
// 		}
// 		secretKeysData = append(secretKeysData, secretKeyData)
// 		log.Println("Secret keys data", secretKeysData)
// 		log.Println("property", property)
// 	}
// 	externalSecret := ExternalSecret{
// 		APIVersion: "external-secrets.io/v1beta1",
// 		Kind:       "ExternalSecret",
// 		Metadata: ObjectMeta{
// 			Name:      externalSecretName,
// 			Namespace: namespace,
// 		},
// 		Spec: ExternalSecretSpec{
// 			RefreshInterval: "10s",
// 			Target: ExternalSecretTarget{
// 				Name:     secretName,
// 				Template: ExternalSecretTargetTemplate{Type: secretType}},
// 			SecretStoreRef: SecretStoreRef{
// 				Name: secretStoreRefName,
// 				Kind: "SecretStore",
// 			},
// 			Data: secretKeysData,
// 		},
// 	}

// 	externalSecretData, err := yaml.Marshal(&externalSecret)
// 	if err != nil {
// 		return
// 	}

// 	_, _, err = k.DynamicClient.CreateResource(ctx, []byte(externalSecretData))
// 	if err != nil {
// 		err = fmt.Errorf("failed to create vault external secret %s/%s, %v", namespace, externalSecretName, err)
// 		return
// 	}
// 	return
// }

// func (k *K8SClient) CreateOrUpdateExternalSecret(ctx context.Context, externalSecretName, namespace,
// 	secretStoreRefName, secretName, secretType string, vaultKeyPathdata, secretProperties map[string][]string) (err error) {
// 	secretKeysData := []ExternalSecretData{}
// 	for key, paths := range vaultKeyPathdata {
// 		// Ensure that the length of paths and properties are the same
// 		if len(paths) != len(secretProperties[key]) {
// 			err = fmt.Errorf("length of paths and properties must be the same for key: %s", key)
// 			return
// 		}
// 		for i, path := range paths {
// 			property := secretProperties[key][i]
// 			secretKeyData := ExternalSecretData{
// 				SecretKey: key,
// 				RemoteRef: ExternalSecretDataRemoteRef{
// 					Key:      path,
// 					Property: property,
// 				},
// 			}
// 			secretKeysData = append(secretKeysData, secretKeyData)

// 			//	log.Println("property", property)
// 		}
// 	}
// 	externalSecret := ExternalSecret{
// 		APIVersion: "external-secrets.io/v1beta1",
// 		Kind:       "ExternalSecret",
// 		Metadata: ObjectMeta{
// 			Name:      externalSecretName,
// 			Namespace: namespace,
// 		},
// 		Spec: ExternalSecretSpec{
// 			RefreshInterval: "10s",
// 			Target: ExternalSecretTarget{
// 				Name:     secretName,
// 				Template: ExternalSecretTargetTemplate{Type: secretType}},
// 			SecretStoreRef: SecretStoreRef{
// 				Name: secretStoreRefName,
// 				Kind: "SecretStore",
// 			},
// 			Data: secretKeysData,
// 		},
// 	}
// 	log.Println("Secret keys data", secretKeysData)
// 	externalSecretData, err := yaml.Marshal(&externalSecret)
// 	if err != nil {
// 		return
// 	}

// 	_, _, err = k.DynamicClient.CreateResource(ctx, []byte(externalSecretData))
// 	if err != nil {
// 		err = fmt.Errorf("failed to create vault external secret %s/%s, %v", namespace, externalSecretName, err)
// 		return
// 	}
// 	return
// }

// func (k *K8SClient) CreateOrUpdateExternalSecret(ctx context.Context, externalSecretName, namespace,
// 	secretStoreRefName, secretName, secretType string, vaultKeyPathdata, secretProperties map[string][]string) (err error) {
// 	secretKeysData := []ExternalSecretData{}

// 	// Extract and sort the keys
// 	keys := make([]string, 0, len(vaultKeyPathdata))
// 	for key := range vaultKeyPathdata {
// 		keys = append(keys, key)
// 	}
// 	sort.Strings(keys)

// 	for _, key := range keys {
// 		paths := vaultKeyPathdata[key]
// 		properties := secretProperties[key]

// 		// Ensure that the length of paths and properties are the same
// 		if len(paths) != len(properties) {
// 			err = fmt.Errorf("length of paths and properties must be the same for key: %s", key)
// 			return
// 		}

// 		for i, path := range paths {
// 			property := properties[i]
// 			secretKeyData := ExternalSecretData{
// 				SecretKey: key,
// 				RemoteRef: ExternalSecretDataRemoteRef{
// 					Key:      path,
// 					Property: property,
// 				},
// 			}
// 			secretKeysData = append(secretKeysData, secretKeyData)
// 		}
// 	}

// 	externalSecret := ExternalSecret{
// 		APIVersion: "external-secrets.io/v1beta1",
// 		Kind:       "ExternalSecret",
// 		Metadata: ObjectMeta{
// 			Name:      externalSecretName,
// 			Namespace: namespace,
// 		},
// 		Spec: ExternalSecretSpec{
// 			RefreshInterval: "10s",
// 			Target: ExternalSecretTarget{
// 				Name:     secretName,
// 				Template: ExternalSecretTargetTemplate{Type: secretType}},
// 			SecretStoreRef: SecretStoreRef{
// 				Name: secretStoreRefName,
// 				Kind: "SecretStore",
// 			},
// 			Data: secretKeysData,
// 		},
// 	}
// 	log.Println("Secret keys data", secretKeysData)
// 	externalSecretData, err := yaml.Marshal(&externalSecret)
// 	if err != nil {
// 		return
// 	}

// 	_, _, err = k.DynamicClient.CreateResource(ctx, []byte(externalSecretData))
// 	if err != nil {
// 		err = fmt.Errorf("failed to create vault external secret %s/%s, %v", namespace, externalSecretName, err)
// 		return
// 	}
// 	return
// }

// func (k *K8SClient) CreateOrUpdateExternalSecret(ctx context.Context, externalSecretName, namespace,
// 	secretStoreRefName, secretName, secretType string, vaultKeyPathdata, secretProperties map[string][]string) (err error) {
// 	secretKeysData := []ExternalSecretData{}

// 	// Extract and sort the keys
// 	keys := make([]string, 0, len(vaultKeyPathdata))
// 	for key := range vaultKeyPathdata {
// 		keys = append(keys, key)
// 	}
// 	sort.Strings(keys)

// 	for _, key := range keys {
// 		paths := vaultKeyPathdata[key]
// 		properties := secretProperties[key]

// 		// Ensure that the length of paths and properties are the same
// 		if len(paths) != len(properties) {
// 			err = fmt.Errorf("length of paths and properties must be the same for key: %s", key)
// 			return
// 		}

// 		for i, path := range paths {
// 			property := properties[i]
// 			secretKeyData := ExternalSecretData{
// 				SecretKey: key,
// 				RemoteRef: ExternalSecretDataRemoteRef{
// 					Key:      path,
// 					Property: property,
// 				},
// 			}
// 			secretKeysData = append(secretKeysData, secretKeyData)
// 		}
// 	}

// 	// Sort the secretKeysData to ensure consistent order
// 	sort.Slice(secretKeysData, func(i, j int) bool {
// 		if secretKeysData[i].SecretKey != secretKeysData[j].SecretKey {
// 			return secretKeysData[i].SecretKey < secretKeysData[j].SecretKey
// 		}
// 		return secretKeysData[i].RemoteRef.Property < secretKeysData[j].RemoteRef.Property
// 	})

// 	externalSecret := ExternalSecret{
// 		APIVersion: "external-secrets.io/v1beta1",
// 		Kind:       "ExternalSecret",
// 		Metadata: ObjectMeta{
// 			Name:      externalSecretName,
// 			Namespace: namespace,
// 		},
// 		Spec: ExternalSecretSpec{
// 			RefreshInterval: "10s",
// 			Target: ExternalSecretTarget{
// 				Name:     secretName,
// 				Template: ExternalSecretTargetTemplate{Type: secretType}},
// 			SecretStoreRef: SecretStoreRef{
// 				Name: secretStoreRefName,
// 				Kind: "SecretStore",
// 			},
// 			Data: secretKeysData,
// 		},
// 	}
// 	log.Println("Secret keys data", secretKeysData)
// 	externalSecretData, err := yaml.Marshal(&externalSecret)
// 	if err != nil {
// 		return
// 	}

// 	_, _, err = k.DynamicClient.CreateResource(ctx, []byte(externalSecretData))
// 	if err != nil {
// 		err = fmt.Errorf("failed to create vault external secret %s/%s, %v", namespace, externalSecretName, err)
// 		return
// 	}
// 	return
// }

// func (k *K8SClient) CreateOrUpdateExternalSecret(ctx context.Context, externalSecretName, namespace,
// 	secretStoreRefName, secretName, secretType string, vaultKeyPathdata, secretProperties map[string][]string) (err error) {
// 	secretKeysData := []ExternalSecretData{}

// 	// Extract and sort the keys
// 	keys := make([]string, 0, len(vaultKeyPathdata))
// 	for key := range vaultKeyPathdata {
// 		keys = append(keys, key)
// 	}
// 	sort.Strings(keys)

// 	// Iterate over the sorted keys
// 	for _, key := range keys {
// 		paths := vaultKeyPathdata[key]
// 		properties := secretProperties[key]

// 		// Ensure that the length of paths and properties are the same
// 		if len(paths) != len(properties) {
// 			err = fmt.Errorf("length of paths and properties must be the same for key: %s", key)
// 			return
// 		}

// 		// Sort paths and properties to ensure consistent order
// 		sort.Strings(paths)
// 		sort.Strings(properties)

// 		// Iterate over the sorted paths and properties
// 		for i := range paths {
// 			path := paths[i]
// 			property := properties[i]
// 			secretKeyData := ExternalSecretData{
// 				SecretKey: key,
// 				RemoteRef: ExternalSecretDataRemoteRef{
// 					Key:      path,
// 					Property: property,
// 				},
// 			}
// 			secretKeysData = append(secretKeysData, secretKeyData)
// 		}
// 	}

// 	// Sort the secretKeysData to ensure consistent order
// 	sort.Slice(secretKeysData, func(i, j int) bool {
// 		if secretKeysData[i].SecretKey != secretKeysData[j].SecretKey {
// 			return secretKeysData[i].SecretKey < secretKeysData[j].SecretKey
// 		}
// 		return secretKeysData[i].RemoteRef.Property < secretKeysData[j].RemoteRef.Property
// 	})

// 	externalSecret := ExternalSecret{
// 		APIVersion: "external-secrets.io/v1beta1",
// 		Kind:       "ExternalSecret",
// 		Metadata: ObjectMeta{
// 			Name:      externalSecretName,
// 			Namespace: namespace,
// 		},
// 		Spec: ExternalSecretSpec{
// 			RefreshInterval: "10s",
// 			Target: ExternalSecretTarget{
// 				Name:     secretName,
// 				Template: ExternalSecretTargetTemplate{Type: secretType}},
// 			SecretStoreRef: SecretStoreRef{
// 				Name: secretStoreRefName,
// 				Kind: "SecretStore",
// 			},
// 			Data: secretKeysData,
// 		},
// 	}
// 	log.Println("Secret keys data", secretKeysData)
// 	externalSecretData, err := yaml.Marshal(&externalSecret)
// 	if err != nil {
// 		return
// 	}

//		_, _, err = k.DynamicClient.CreateResource(ctx, []byte(externalSecretData))
//		if err != nil {
//			err = fmt.Errorf("failed to create vault external secret %s/%s, %v", namespace, externalSecretName, err)
//			return
//		}
//		return
//	}
func (k *K8SClient) CreateOrUpdateExternalSecret(ctx context.Context, externalSecretName, namespace,
	secretStoreRefName, secretName, secretType string, vaultKeyPathdata, secretProperties map[string][]string) (err error) {
	secretKeysData := []ExternalSecretData{}

	// Extract and sort the keys from vaultKeyPathdata
	keys := make([]string, 0, len(vaultKeyPathdata))
	for key := range vaultKeyPathdata {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		paths := vaultKeyPathdata[key]

		// Ensure that secretProperties contains the key
		properties, exists := secretProperties[key]
		if !exists {
			err = fmt.Errorf("secretProperties does not contain key: %s", key)
			return
		}

		// Ensure that the length of paths and properties are the same
		if len(paths) != len(properties) {
			err = fmt.Errorf("length of paths and properties must be the same for key: %s", key)
			return
		}

		// Sort paths and properties to ensure consistent order
		sortedData := make([]struct {
			Path, Property string
		}, len(paths))

		for i := range paths {
			sortedData[i] = struct {
				Path, Property string
			}{
				Path:     paths[i],
				Property: properties[i],
			}
		}

		// Sort the combined data
		sort.Slice(sortedData, func(i, j int) bool {
			if sortedData[i].Path != sortedData[j].Path {
				return sortedData[i].Path < sortedData[j].Path
			}
			return sortedData[i].Property < sortedData[j].Property
		})

		// Append sorted data to secretKeysData
		for _, data := range sortedData {
			secretKeyData := ExternalSecretData{
				SecretKey: key,
				RemoteRef: ExternalSecretDataRemoteRef{
					Key:      data.Path,
					Property: data.Property,
				},
			}
			secretKeysData = append(secretKeysData, secretKeyData)
		}
	}

	// Sort the secretKeysData to ensure consistent order
	sort.Slice(secretKeysData, func(i, j int) bool {
		if secretKeysData[i].SecretKey != secretKeysData[j].SecretKey {
			return secretKeysData[i].SecretKey < secretKeysData[j].SecretKey
		}
		return secretKeysData[i].RemoteRef.Property < secretKeysData[j].RemoteRef.Property
	})

	externalSecret := ExternalSecret{
		APIVersion: "external-secrets.io/v1beta1",
		Kind:       "ExternalSecret",
		Metadata: ObjectMeta{
			Name:      externalSecretName,
			Namespace: namespace,
		},
		Spec: ExternalSecretSpec{
			RefreshInterval: "10s",
			Target: ExternalSecretTarget{
				Name:     secretName,
				Template: ExternalSecretTargetTemplate{Type: secretType}},
			SecretStoreRef: SecretStoreRef{
				Name: secretStoreRefName,
				Kind: "SecretStore",
			},
			Data: secretKeysData,
		},
	}
	log.Println("Secret keys data", secretKeysData)
	externalSecretData, err := yaml.Marshal(&externalSecret)
	if err != nil {
		return
	}

	_, _, err = k.DynamicClient.CreateResource(ctx, []byte(externalSecretData))
	if err != nil {
		err = fmt.Errorf("failed to create vault external secret %s/%s, %v", namespace, externalSecretName, err)
		return
	}
	return
}

package registry_metadata

import (
	"context"
	"encoding/json"
	"github.com/google/go-containerregistry/pkg/v1/types"
	ociSpec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/sirupsen/logrus"
	orasContent "oras.land/oras-go/v2/content"
	orasRegistry "oras.land/oras-go/v2/registry"
	orasRemote "oras.land/oras-go/v2/registry/remote"
)

type ImageRequestQuery struct {
	Reference        string `json:"reference" binding:"required"`
	ociSpec.Platform `json:",inline"`
}

type ImageResponseData struct {
	Manifest ociSpec.Manifest `json:"manifest"`
	Config   ociSpec.Image    `json:"config"`
}

func ConstructRegistryMetadata(logEntry *logrus.Entry,
	query ImageRequestQuery) (*ImageResponseData, error) {

	parseReference, err := orasRegistry.ParseReference(query.Reference)
	if err != nil {
		logEntry.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to parse image reference")
	}

	r, err := orasRemote.NewRepository(parseReference.Registry + "/" + parseReference.Repository)

	if err != nil {
		logEntry.WithFields(logrus.Fields{
			"error": err,
		}).Error("error constructing remote repository")
		return nil, err
	}

	ctx := context.Background()
	registryMetadata := &ImageResponseData{}

	registryMetadata.Manifest, err = getFullImageManifest(ctx, logEntry, r, parseReference, query.Platform)
	if err != nil {
		logEntry.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to get full image manifest")
		return nil, err
	}

	registryMetadata.Config, err = getConfig(ctx, logEntry, r, registryMetadata.Manifest.Config)
	if err != nil {
		logEntry.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to get image config")
		return nil, err
	}

	return registryMetadata, nil
}

func getFullImageManifest(ctx context.Context, logEntry *logrus.Entry,
	repository *orasRemote.Repository, reference orasRegistry.Reference,
	platform ociSpec.Platform) (ociSpec.Manifest, error) {

	descriptor, rc, err := repository.FetchReference(ctx, reference.String())
	if err != nil {
		logEntry.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to fetch image manifest")
		return ociSpec.Manifest{}, err
	}

	imageContent, err := orasContent.ReadAll(rc, descriptor)
	if err != nil {
		logEntry.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to read image manifest")
		return ociSpec.Manifest{}, err
	}

	if descriptor.MediaType == ociSpec.MediaTypeImageIndex || descriptor.MediaType == string(types.DockerManifestList) {
		var index ociSpec.Index
		if err := json.Unmarshal(imageContent, &index); err != nil {
			logEntry.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Failed to unmarshal image manifest")
			return ociSpec.Manifest{}, err
		}

		var manifestForPlatform ociSpec.Descriptor
		for _, m := range index.Manifests {
			if m.Platform.Architecture == platform.Architecture && m.Platform.OS == platform.OS {
				manifestForPlatform = m
			}
		}
		reference.Reference = manifestForPlatform.Digest.String()
		return getFullImageManifest(ctx, logEntry, repository, reference, ociSpec.Platform{})
	}

	var manifest ociSpec.Manifest
	if err := json.Unmarshal(imageContent, &manifest); err != nil {
		logEntry.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to unmarshal image manifest")
		return ociSpec.Manifest{}, err
	}
	manifest.Subject = &descriptor

	return manifest, nil
}

func getConfig(ctx context.Context, logEntry *logrus.Entry,
	repository orasContent.Fetcher, descriptor ociSpec.Descriptor) (ociSpec.Image, error) {

	configContent, err := orasContent.FetchAll(ctx, repository, descriptor)
	if err != nil {
		logEntry.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to fetch image config")
		return ociSpec.Image{}, err
	}

	var config ociSpec.Image
	if err := json.Unmarshal(configContent, &config); err != nil {
		logEntry.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to unmarshal image config")
		return ociSpec.Image{}, err
	}

	return config, nil
}

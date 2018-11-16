package engine

import (
	"bytes"
	"encoding/json"
	"github.com/seizadi/cmdb/pkg/pb"
	"os/exec"
	"strings"
)
import "github.com/seizadi/template-engine/util"

// TODO - Check errors from file operations

func CreateValues(providers *pb.ListCloudProvidersResponse) (*pb.Application, error) {
	// TODO - We pick the first one but we should go through the list
	provider := providers.Results[0]
	region := provider.Regions[0]
	// TODO - We pick the first one but we should go through the list
	env := region.Environments[0]
	err := util.PutMap(Environment{Environment: env},"tmp/environment.yaml")
	if err != nil {
		return nil,err
	}
	
	// TODO - We pick the first one but we should grow through the list
	application := env.Applications[0]

	versionTag := application.VersionTag

	err = util.PutText(versionTag.Name, "tmp/release_name.txt")
	if err != nil {
		return nil,err
	}
	
	err = util.PutMap(Application{Application: application},"tmp/application.yaml")
	if err != nil {
		return nil,err
	}
	
	containers := application.Containers
	err = util.PutMap(Containers{Containers: containers},"tmp/containers.yaml")
	if err != nil {
		return nil,err
	}
	
	manifest := application.Manifest
	var manifestValues ManifestValues
	json.Unmarshal([]byte(manifest.Values.Value), &manifestValues)
	var manifestServices []ManifestService
	json.Unmarshal([]byte(manifest.Services.Value), &manifestServices)
	var manifestIngress ManifestIngress
	json.Unmarshal([]byte(manifest.Ingress.Value), &manifestIngress)
	err = util.PutMap(Manifest{
		Manifest: manifest,
		ManifestValues: manifestValues.Values,
		Services: manifestServices,
		Ingress: manifestIngress,
	},"tmp/manifest.yaml")
	if err != nil {
		return nil,err
	}

	vault := manifest.Vault
	err = util.PutMap(Vault{Vault: vault},"tmp/vault.yaml")
	if err != nil {
		return nil,err
	}
	
	secrets := vault.Secrets
	err = util.PutMap(Secrets{Secrets: secrets},"tmp/secrets.yaml")
	if err != nil {
		return nil,err
	}
	
	namespace := strings.ToLower(application.Name + "-" + env.Name + "-" + (env.Code).String())
	namespace = strings.Replace(namespace, " ", "-", -1)
	err = util.PutMap(NameSpace{NameSpace: namespace},"tmp/name_space.yaml")
	if err != nil {
		return nil,err
	}
	
	return application,nil
}

func GetArchiveFile(app *pb.Application) string {
	return app.VersionTag.Repo + "/" + app.VersionTag.Version + ".zip"
}

func ResolveManifest(path string) error {
	var out bytes.Buffer
	helmCmd := "#!/bin/bash\n" +
		"helm template " +
	"-f ./tmp/name_space.yaml " +
	"-f ./tmp/containers.yaml " +
	"-f ./tmp/manifest.yaml " +
	"-f ./tmp/vault.yaml " +
	"-f ./tmp/secrets.yaml " +
	"-f ./tmp/cmdb-0.0.4/repo/cmdb/cmdb.yaml " +
	"-n $(cat ./tmp/release_name.txt) " +
	"./"+ path
	cmd := exec.Command("echo", helmCmd)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	util.CopyBufferContents(out.Bytes(), "./tmp/run.sh")
	
	out.Reset()
	cmd = exec.Command("bash", "./tmp/run.sh")
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return err
	}
	util.CopyBufferContents(out.Bytes(), "./tmp/resolved.yaml")
	
	return nil
}
package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/seizadi/cmdb/pkg/pb"
	"os/exec"
	"strings"
)
import "github.com/seizadi/template-engine/util"

// TODO - Check errors from file operations

func CreateTemplate(regions *pb.ListRegionsResponse) {
	region := regions.Results[0]
	// TODO - We pick the first one but we should go through the list
	env := region.Environments[0]
	util.PutMap(Environment{Environment: env},"tmp/environment.yaml")
	
	// TODO - We pick the first one but we should grow through the list
	application := env.Applications[0]

	versionTag := application.VersionTag

	util.PutText(versionTag.Name, "tmp/release_name.txt")

	util.PutMap(Application{Application: application},"tmp/application.yaml")

	containers := application.Containers
	util.PutMap(Containers{Containers: containers},"tmp/containers.yaml")

	manifest := application.Manifest
	var manifestValues ManifestValues
	json.Unmarshal([]byte(manifest.Values.Value), &manifestValues)
	var manifestServices []ManifestService
	json.Unmarshal([]byte(manifest.Services.Value), &manifestServices)
	var manifestIngress ManifestIngress
	json.Unmarshal([]byte(manifest.Ingress.Value), &manifestIngress)
	util.PutMap(Manifest{
		Manifest: manifest,
		ManifestValues: manifestValues.Values,
		Services: manifestServices,
		Ingress: manifestIngress,
	},"tmp/manifest.yaml")


	vault := manifest.Vault
	util.PutMap(Vault{Vault: vault},"tmp/vault.yaml")

	secrets := vault.Secrets
	util.PutMap(Secrets{Secrets: secrets},"tmp/secrets.yaml")

	namespace := strings.ToLower(application.Name + "-" + env.Name + "-" + (env.Code).String())
	namespace = strings.Replace(namespace, " ", "-", -1)
	util.PutMap(NameSpace{NameSpace: namespace},"tmp/name_space.yaml")
	
	//cmd := exec.Command("/usr/local/bin/helm",
	//	"template",
	//	"-f ./tmp/name_space.yaml",
	//	"-f ./tmp/containers.yaml",
	//	"-f ./tmp/manifest.yaml",
	//	"-f ./tmp/vault.yaml",
	//	"-f ./tmp/secrets.yaml",
	//	"-f ./tmp/cmdb-0.0.4/repo/cmdb/cmdb.yaml",
	//	"-n " + versionTag.Name,
	//	"./tmp/cmdb-0.0.4/repo/cmdb")
	cmd := exec.Command("bash", "run.sh")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return
	}
	
	util.CopyBufferContents(out.Bytes(), "./tmp/resolve.yaml")
}

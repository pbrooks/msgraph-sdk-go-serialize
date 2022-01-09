package main

import (
	"fmt"
	"os"
	"testing"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	azure "github.com/microsoft/kiota/authentication/go/azure"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models/microsoft/graph"
	"github.com/microsoftgraph/msgraph-sdk-go/schemaextensions"
	"github.com/stretchr/testify/assert"
)

// errors are ignored, as they don't catch the problem

func authenticate() *msgraphsdk.GraphRequestAdapter{
    // Fake a basic authentication
    // Interest in having msgraphsdk construct json serializers

	os.Setenv("AZURE_TENANT_ID", uuid.New().String())
	os.Setenv("AZURE_CLIENT_ID", uuid.New().String())
	os.Setenv("AZURE_USERNAME", "_")
	os.Setenv("AZURE_PASSWORD", "_")
	credential, _ := azidentity.NewEnvironmentCredential(nil)

	auth, _ := azure.NewAzureIdentityAuthenticationProviderWithScopes(credential, []string{})
	requestAdapter, _ := msgraphsdk.NewGraphRequestAdapter(auth)

    msgraphsdk.NewGraphServiceClient(requestAdapter)

    return requestAdapter
}

func TestSerialization(t *testing.T) {

    requestAdapter := authenticate()

	extension_id := "test_id"
	extension_desc := "test_desc"
	extension := graph.NewSchemaExtension()
	extension.SetId(&extension_id)
	extension.SetDescription(&extension_desc)
	extension.SetTargetTypes([]string{"Device"})

	options := schemaextensions.SchemaExtensionsRequestBuilderPostOptions{
		Body: extension,
	}
	extensionRequestBuilder := schemaextensions.NewSchemaExtensionsRequestBuilder("", requestAdapter)
	requestInfo, _ := extensionRequestBuilder.CreatePostRequestInformation(&options)

    additionalData := extension.GetAdditionalData()
    assert.NotNil(t, additionalData, "AdditionalData isn't nil")
    for key, value := range additionalData {
        fmt.Println("AdditionalData item ", key, value)
    }

	assert.Equal(t, "{\"id\":\"test_id\",\"description\":\"test_desc\",\"targetTypes\":[\"Device\"]}", 
        string(requestInfo.Content), "Expected JSON matches output")

	fmt.Println("Content ", string(requestInfo.Content))
}

func TestSerialization_Mitigation(t *testing.T) {

    requestAdapter := authenticate()

	extension_id := "test_id"
	extension_desc := "test_desc"
	extension := graph.NewSchemaExtension()
	extension.SetId(&extension_id)
	extension.SetDescription(&extension_desc)
	extension.SetTargetTypes([]string{"Device"})
    extension.SetAdditionalData(nil)

	options := schemaextensions.SchemaExtensionsRequestBuilderPostOptions{
		Body: extension,
	}
	extensionRequestBuilder := schemaextensions.NewSchemaExtensionsRequestBuilder("", requestAdapter)
	requestInfo, _ := extensionRequestBuilder.CreatePostRequestInformation(&options)

    additionalData := extension.GetAdditionalData()
    assert.Nil(t, additionalData, "AdditionalData is nil")
    for key, value := range additionalData {
        fmt.Println("AdditionalData item ", key, value)
    }

	assert.Equal(t, "{\"id\":\"test_id\",\"description\":\"test_desc\",\"targetTypes\":[\"Device\"]}", 
        string(requestInfo.Content), "Expected JSON matches output")

	fmt.Println("Content ", string(requestInfo.Content))
}

// az deployment group create --resource-group test-aca --template-file main.bicep --parameters @main.parameters.json
// az account set --subscription "00f9e387-4d2a-4d2d-8113-785d5827878f"

param location string = resourceGroup().location

// ---------------------------------------------------------------------------------------------
//  Optional Parameters
//    Each resource has an optional override for the default azd resource naming conventions.
//    Update the main.parameters.json file to specify them:
//
//    "webServiceName": {
//      "value": "my-web-service"
//    }
// ---------------------------------------------------------------------------------------------

param azureContainerAppsName string = ''
param azureContainerRegistryName string = ''
param azureContainerRegistryUserName string = ''
@secure()
param azureContainerRegistryPassword string = ''
param azureEventHubConnectionString string = ''


// create the azure container registry
// resource acr 'Microsoft.ContainerRegistry/registries@2021-09-01' = {
//   name: toLower('${resourceGroup().name}acr')
//   location: location
//   sku: {
//     name: 'Basic'
//   }
//   properties: {
//     adminUserEnabled: true
//   }
// }

// create the aca environment
module env 'environment.bicep' = {
  name: 'containerAppEnvironment'
  params: {
    acaName: azureContainerAppsName
    location: location
  }
}

// create the various config pairs
var shared_config = [
  {
    name: 'ASPNETCORE_ENVIRONMENT'
    value: 'Development'
  }
  {
    name: 'APPINSIGHTS_INSTRUMENTATIONKEY'
    value: env.outputs.appInsightsInstrumentationKey
  }
  {
    name: 'APPLICATIONINSIGHTS_CONNECTION_STRING'
    value: env.outputs.appInsightsConnectionString
  } 
  {
    name: 'load-test-requests'
    value: '0'
  }
  {
    name: 'load-test-url'
    value: 'https://apim-mckq3zpiubjie.azure-api.net/mock/get'
  }
  {
    name: 'load-test-eventHub'
    value: azureEventHubConnectionString
  }
]

// create the service container app
module loadTest 'container_app.bicep' = {
  name: 'resty-load-test'
  params: {
    name: 'resty-load-test'
    location: location
    registryPassword: azureContainerRegistryPassword
    registryUsername: azureContainerRegistryUserName
    containerAppEnvironmentId: env.outputs.id
    registry: azureContainerRegistryName
    envVars: shared_config
  }
}

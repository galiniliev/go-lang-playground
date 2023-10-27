param name string
param location string = resourceGroup().location
param containerAppEnvironmentId string
param repositoryImage string = 'apimloadtest.azurecr.io/loadtestresty:latest'
param envVars array = []
param registry string
param minReplicas int = 1
param maxReplicas int = 1
param registryUsername string
@secure()
param registryPassword string
@secure()
param azureEventHubConnectionString string

resource containerApp 'Microsoft.App/containerApps@2023-05-01' ={
  name: name
  location: location
  properties:{
    managedEnvironmentId: containerAppEnvironmentId
    configuration: {
      activeRevisionsMode: 'multiple'
      secrets: [
        {
          name: 'container-registry-password'
          value: registryPassword
        }
        {
          name: 'load-test-eventhub'
          value: azureEventHubConnectionString
        }
      ]      
      registries: [
        {
          server: registry
          username: registryUsername
          passwordSecretRef: 'container-registry-password'
        }
      ]
      ingress: {
        external: true
        targetPort: 80
        transport: 'http'
        allowInsecure: true
      }
    }
    template: {
      containers: [
        {
          image: repositoryImage
          name: name
          env: envVars
        }
      ]
      scale: {
        minReplicas: minReplicas
        maxReplicas: maxReplicas
      }
    }
  }
}

output fqdn string = containerApp.properties.configuration.ingress.fqdn

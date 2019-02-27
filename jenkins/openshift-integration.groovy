def testRunner = evaluate(readTrusted('jenkins/vars/testRunner.groovy'))

testRunner {
    testScript = "hack/integration.sh"
    testType   = "integration"
    kubeconfigCredentialsID = 'openshift-metering-ci-kubeconfig'
    deployPlatform = "openshift"
    meteringHttpsAPI = true
}

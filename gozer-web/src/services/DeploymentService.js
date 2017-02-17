import _ from 'underscore';
import Vue from 'vue';
import VueResource from 'vue-resource';

Vue.use(VueResource);

const deploymentsResource = Vue.resource('/deployments{/name}');

function convertToKeyValueObjects(kvp) {
  return {
    key: kvp[0],
    value: kvp[1]
  }
}

function transformDeploymentToViewObj(deployment) {

  return {
    $$originalObject: deployment,
    name: deployment.metadata.name,
    uid: deployment.metadata.uid,
    labels: _.pairs(deployment.metadata.labels)
             .map(convertToKeyValueObjects),
    containers: deployment.spec.template.spec.containers
                          .map((container) => {
                            return {
                              name: container.name,
                              image: container.image,
                              ports: container.ports
                            }
                          }),
    maxUnresponsiveTimeBeforeTermination: deployment.spec.template.spec.terminationGracePeriodSeconds,
    status: _.pairs(_.omit(deployment.status, 'conditions'))
             .map(convertToKeyValueObjects)
  };
}

function list() {
  return deploymentsResource.get().then(function (response) {
    return response.data.deployments.map(transformDeploymentToViewObj)
  });
}

function get(name) {
  return deploymentsResource.get({name: name}).then(function (response) {
    return transformDeploymentToViewObj(response.data.deployment);
  })
}

export default {list, get}


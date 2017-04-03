import _ from 'underscore';
import Vue from 'vue';
import Promise from 'bluebird';
import VueResource from 'vue-resource';
import imagesService from './ImagesService';
import moment from 'moment';

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
  return  deploymentsResource.get({name: name}).then(function (response) {

    return transformDeploymentToViewObj(response.data);
  })
}

function getDeploymentWithImageData(name) {
  return get(name)
    .then(function (deployment) {
      return Promise.reduce(
        deployment.containers,
        function (image, container) {
          if (image) {
            return image;
          }
          else {
            let containerImageParts = container.image.split(':');
            let containerName = containerImageParts[0];
            let containerVersion = containerImageParts[1];

            return imagesService
              .getByVersion(containerName, containerVersion)
              .then(function (imageData) {
                imageData.containerName = containerName;
                return imageData;
              })
          }
        }, null).then(function (image) {
        if (image) {
          return imagesService
            .get(image.containerName)
            .then(function (list) {
              list.forEach(function (image) {
                let imageParts = image.name.split(':');
                let imageName = imageParts[0];
                let imageVersion = imageParts[1];
                image.age = moment(image.metadata["created-at"]).fromNow();
                image.deploymentImageName = [imageName, imageVersion].join(':');
              });
              list.sort(function (leftImage, rightImage) {return moment(leftImage.metadata['created-at']).isBefore(rightImage.metadata['created-at'])});
              return {
                deployment: deployment,
                deploymentAppImage: image,
                deploymentAppImageName: image.name,
                deploymentImages: list
              }
            });
        }
        else {
          Promise.reject('No image data found');
        }
      });
    });
}

function set(deploymentName, container_name, image_name) {
  return deploymentsResource
    .update({name: deploymentName},
            {
              container_name: container_name,
              image: image_name
            })
}

export default {list, get, set, getDeploymentWithImageData}


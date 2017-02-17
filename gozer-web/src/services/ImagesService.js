import _ from 'underscore';
import Vue from 'vue';

import VueResource from 'vue-resource';

Vue.use(VueResource);

const storedContainersResource = Vue.resource('/images{/name}{/version}');

function get(name) {
  return storedContainersResource
    .get({name: name})
    .then(function (res) {
      return res.data;
    });
}

function getByVersion(name, version) {
  return storedContainersResource
    .get({name:name, version:version})
    .then(function (res) {
      return res.data;
    });
    // .then(function (images) {
    //   return _.find(images, function (item) {
    //     let nameParts = item.name.split(':');
    //     return nameParts[1] === version;
    //   })
    // })
}

export default {get, getByVersion}

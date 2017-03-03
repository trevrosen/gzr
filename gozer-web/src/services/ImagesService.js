import _ from 'underscore';
import Vue from 'vue';

import VueResource from 'vue-resource';

Vue.use(VueResource);

const storedContainersResource = Vue.resource('/images{/name}{/version}');

function get(name) {
  return storedContainersResource
    .get({name: name})
    .then(function (res) {
      return res.data.images;
    });
}

function getByVersion(name, version) {
  return storedContainersResource
    .get({name:name, version:version})
    .then(function (res) {
      return res.data;
    });
}

export default {get, getByVersion}

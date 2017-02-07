import _ from 'underscore';
import Vue from 'vue';

import VueResource from 'vue-resource';

Vue.use(VueResource);

const storedContainersResource = Vue.resource('/static/images.json');

function get(name){
    return storedContainersResource.get({name:name}).then(function (res) {
      if(name === "bypass/admin")
        return res.data.images;
      else return []
    });
}

function getByVersion(name, version) {
  return get(name)
    .then(function(images){
       return  _.find(images, function(item){
         let nameParts = item.name.split(':');
         return nameParts[1] === version;
       })
    })
}


export default {get, getByVersion}

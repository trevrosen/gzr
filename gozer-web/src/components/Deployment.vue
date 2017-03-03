<template>
  <div>
    <ol class="breadcrumb">
      <li>
        <router-link to="/">&lt; Back</router-link>
      </li>
      <li v-if="deployment">
        <span class="active">{{deployment.name}}</span>
      </li>
    </ol>
    <div v-if="loading">
      <spinner v-model="loading" size="lg" fixed text="Loading Deployments"></spinner>
    </div>
    <div v-else-if="!error" class="panel panel-default">
      <div class="panel-heading">
        <h2>{{deployment.name}}</h2>
      </div>
      <div class="panel-body">
        <div class="list-group">
          <div class="list-group-item" :class="{active: image.name === deploymentAppImageName}" v-for="image in deploymentImages">
            <h4 class="list-group-item-heading">{{image.name}}<span class="label label-default pull-right">{{image.age}}</span></h4>
            <div class="list-group-item-text">
              <div class="row">
                <div class="col-xs-10">
                  <dl class="dl-horizontal">
                    <dt>Commit</dt>
                    <dd>{{image.metadata['git-commit']}} (Github button)</dd>

                    <dt>Annotations</dt>
                    <dd><span class="label label-info" v-for="annotation in image.metadata['git-annotation']">{{annotation}}</span></dd>

                    <dt>Tags</dt>
                    <dd><span class="label label-info" v-for="tag in image.metadata['git-tag']">{{tag}}</span></dd>
                  </dl>
                </div>
                <div class="col-xs-2">
                  <button class="btn btn-primary">Deploy</button>
                </div>
              </div>

            </div>
          </div>
        </div>


        <!--Image: <v-select v-model="deploymentAppImageName" :options="deploymentImages" search options-value="name" options-label="name"></v-select>-->
        <span></span>
      </div>
    </div>
    <div v-else>
      <div class="row">
        <div class="col-md-4 col-md-offset-4 col-xs-10 col-xs-offset-1">
          <div class="media">
            <div class="media-left media-middle">
              <img src="../img/sad_stay_puft.png" class="media-object">
            </div>
            <div class="media-body media-middle">
              <h3 class="media-heading">Gozer Error</h3>
              Sorry there was an error loading {{name}}.
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
  import { accordion, panel, spinner } from 'vue-strap'
  import Multiselect from 'vue-multiselect'
  import moment from 'moment';
  import imagesService from '../services/ImagesService'
  import deploymentService from '../services/DeploymentService'
  export default {
    props: {name: {type: String}},
    data() {
      return {
        loading: true,
        error: undefined,
        deploymentAppImage: {},
        deploymentImages: [],
        deploymentAppImageName: '',
        deployment: {}
      };
    },
    created: function () {
      const vm = this;
      deploymentService
        .get(vm.name)
        .then(function (deployment) {
          vm.deployment = deployment;
          let promises = vm.deployment
                           .containers
                           .map(function (container) {
                             let containerImageParts = container.image.split(':');
                             let containerName = containerImageParts[0];
                             let containerVersion = containerImageParts[1];
                             return imagesService
                               .getByVersion(containerName, containerVersion)
                               .then(function (image) {
                                 if (image) {
                                   vm.deploymentAppImage = image;
                                   vm.deploymentAppImageName = image.name;
                                   return imagesService
                                     .get(containerName)
                                     .then(function (list) {
                                       list.forEach(function (image) {
                                         image.age = moment(image.metadata["created-at"]).fromNow();
                                       });
                                       vm.deploymentImages = list;
                                     });
                                 }
                               })
                           });
          return Promise.all(promises);
        })
        .finally(function () {
          vm.loading = false;
        })
      ;

    },
    components: {
      accordion,
      panel,
      spinner,
      Multiselect
    }
  };
</script>

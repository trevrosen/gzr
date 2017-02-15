<template>
  <div>
    <div v-if="loading">
      <spinner v-model="loading" size="lg" fixed text="Loading Deployments"></spinner>
    </div>
    <div v-else-if="!error">
      <div class="gzr-cards">
        <div class="flex-row">
          <div v-for="deployment in list">
            <deployment-dash-item :deployment="deployment"></deployment-dash-item>
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <div class="row">
        <div class="col-md-4 col-md-offset-4">
          <div class="media">
            <div class="media-left media-middle">
              <img src="../img/sad_stay_puft.png" class="media-object">
            </div>
            <div class="media-body media-middle">
              <h3 class="media-heading">Gozer Error</h3>
              Sorry there was an error loading deployments.
            </div>
          </div>
        </div>

      </div>


    </div>
  </div>
</template>
<script>
  import {spinner} from 'vue-strap'
  import DeploymentService from '../services/DeploymentService'
  import DeploymentDashItem from './DeploymentDashItem'
  export default {
    name      : 'index',
    data ()  {
      return {
        list   : [],
        loading: true,
        error  : null
      }
    },
    created() {
      const vm = this;
      DeploymentService
        .list()
        .then(function (list) {
          vm.list = list;
        })
        .catch(function (err) {
          vm.error = err;
        })
        .finally(function () {
          vm.loading = false;
        })
    },
    components: {
      DeploymentDashItem,
      spinner
    }
  };
</script>
<style scoped>
  .gzr-cards {
    display: flex;
    flex-direction: column;
  }

  .flex-row {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: space-around;

  }

  .flex-row > div {
    margin-bottom: 15px;
  }

</style>

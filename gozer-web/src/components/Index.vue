<template>
  <div>
    <div class="gzr-cards">
      <div class="flex-row">
        <div v-for="deployment in list">
          <deployment-dash-item :deployment="deployment"></deployment-dash-item>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
  import DeploymentService from '../services/DeploymentService'
  import DeploymentDashItem from './DeploymentDashItem'
  export default {
    name      : 'index',
    data ()  {
      return {
        list: []
      }
    },
    created() {
      const vm = this;
      DeploymentService
        .list()
        .then(function (list) {
          vm.list = list;
        })
    },
    components: {
      DeploymentDashItem
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

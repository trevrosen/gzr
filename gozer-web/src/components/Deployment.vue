<template>
  <div class="panel panel-default">
    <div class="panel-heading">
      <h3 class="panel-title">{{deploymentInternal.name}}</h3>
    </div>
    <div class="panel-body">
      <accordion :one-at-atime="false">
        <panel is-open class="panel-default github-panel">
          <div slot="header">
            <span class="github-header"><img src="../img/GitHub-Mark-32px.png" title="Github"/> </span><span class="link-text">Github</span>
          </div>
          <ul class="list-unstyled" v-for="container in deploymentInternal.containers">
            <li><strong>Container Name:</strong> {{container.name}}</li>
            <li><strong>Container Label:</strong> {{container.image}}</li>
          </ul>
        </panel>
        <panel class="panel-default" header="Containers" is-open>
          <ul class="list-unstyled" v-for="container in deploymentInternal.containers">
            <li><strong>Container Name:</strong> {{container.name}}</li>
            <li><strong>Container Label:</strong> {{container.image}}</li>
          </ul>
        </panel><!-- put in tabs for this other shit-->
        <panel class="panel-default" header="Labels" >
          <ul class="list-unstyled">
            <li v-for="label in deploymentInternal.labels"><strong> {{label.key}}:</strong> {{label.value}}</li>
          </ul>
        </panel>
        <panel class="panel-default" header="Status" >
          <ul class="list-unstyled">
            <li v-for="status in deploymentInternal.status"><strong> {{status.key}}:</strong> {{status.value}}</li>
          </ul>
        </panel>
      </accordion>
    </div>
  </div>
</template>

<script>
  import {accordion, panel} from 'vue-strap'
  export default {
    name: 'deployment',
    props:{deployment:{type:Object}},
    data() {
      return {
        deploymentInternal: JSON.parse(JSON.stringify(this.deployment)),
      };
    },
    components:{
      accordion,
      panel
    }
  };
</script>

<style>

  .btn-github {
    padding: 2px 6px;
  }

  .btn-github img {
    height: 28px;
    width: 28px;
  }

  .github-panel .panel-title a:hover,
  .github-panel .panel-title a:focus{
    text-decoration: none;
  }

  .github-panel .panel-title a:hover .link-text,
  .github-panel .panel-title a:focus .link-text{
    text-decoration: underline;
  }

  .github-header img {
    height: 28px;
    width: 28px;
  }
</style>

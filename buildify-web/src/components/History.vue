<script lang="ts">
import Build from './Build.vue'

export default {
  data() {
    return {
      builds: [],
    }
  },

  props: [
    "link",
    "maxAmount"
  ],

  created: function() {
    var _this = this;
    fetch("http://" + _this.link + "/builds")
      .then(res => res.json())
      .then(out => _this.builds = out.reverse());
  },
  components: {
    Build
  }
}
</script>

<template>
  <div class="builds">
    <template v-for="build, i in builds">
      <Build v-if="maxAmount > i" 
              :baseLink="link"
              :is-first="i == 0"
              :id="build.Id"
              :time="build.Time" 
              :hash="build.Hash" 
              :message="build.Message" 
              :download-link="build.DownloadLink"/>
    </template>
  </div>
</template>
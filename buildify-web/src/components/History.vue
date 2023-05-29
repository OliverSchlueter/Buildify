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
  },

  methods: {
      totalDownloads: function(){
        return this.builds
                .map(build => build.Downloads)
                .reduce((count, val) => count + val)
      },
  }
}
</script>

<template>
  <div class="builds">
    <template v-if="builds.length == 0">
      <p class="error">Could not load data.</p>
    </template>
    <template v-else>
      <p class="center-text">Total downloads: {{ totalDownloads() }}</p>
      <template v-for="build, i in builds">
        <Build v-if="maxAmount > i" 
                :baseLink="link"
                :is-latest="i == 0"
                :id="build.Id"
                :time="build.Time" 
                :hash="build.Hash" 
                :message="build.Message"
                :downloads="build.Downloads"
                :file-name="build.FileName"
                :download-link="build.DownloadLink"/>
      </template>
    </template>

  </div>
</template>
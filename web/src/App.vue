<template>
  <n-config-provider :theme="darkTheme" :hljs="hljs">
    <n-global-style/>
    <n-layout>
      <n-layout-header style="display: none;">
        <div id="nav">
          <router-link to="/">Home</router-link>
          |
          <router-link to="/">About</router-link>
        </div>
      </n-layout-header>
      <n-layout-content>
        <router-view/>
      </n-layout-content>
    </n-layout>
  </n-config-provider>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { NConfigProvider, NGlobalStyle, NLayout, NLayoutHeader, NLayoutContent } from 'naive-ui';
import { darkTheme } from 'naive-ui';

import hljs from 'highlight.js/lib/core';
import json from 'highlight.js/lib/languages/json';

hljs.registerLanguage('json', json);

export default defineComponent({
  components: {
    NConfigProvider,
    NGlobalStyle,
    NLayout,
    NLayoutHeader,
    NLayoutContent,
  },
  setup() {
    return {
      darkTheme,
      hljs,
    };
  },
});
</script>

<style>
</style>

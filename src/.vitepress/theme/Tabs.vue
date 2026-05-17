<template>
  <div class="custom-tabs">
    <div class="tab-headers">
      <button
        v-for="(tab, i) in tabs"
        :key="i"
        :class="['tab-header', { active: activeTab === i }]"
        @click="activeTab = i"
      >
        {{ tab }}
      </button>
    </div>
    <div class="tab-panels">
      <div v-for="(tab, i) in tabs" :key="i" v-show="activeTab === i" class="tab-panel">
        <slot :name="'tab-' + i" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  tabs: {
    type: Array,
    required: true
  }
})

const activeTab = ref(0)
</script>

<style scoped>
.custom-tabs {
  margin: 16px 0;
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  overflow: hidden;
}
.tab-headers {
  display: flex;
  border-bottom: 1px solid var(--vp-c-divider);
  background: var(--vp-c-bg-soft);
}
.tab-header {
  padding: 8px 16px;
  border: none;
  background: none;
  cursor: pointer;
  font-size: 14px;
  color: var(--vp-c-text-2);
  border-bottom: 2px solid transparent;
}
.tab-header.active {
  color: var(--vp-c-brand-1);
  border-bottom-color: var(--vp-c-brand-1);
}
.tab-panel {
  padding: 16px;
}
.tab-panel :deep(img) {
  max-width: 100%;
}
.tab-panel :deep(div[class*="language-"]) {
  margin: 0;
}
</style>

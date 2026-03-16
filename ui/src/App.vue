<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { fetchGraph, type GraphData, type GraphNode } from './services/api'
import SidebarLeft from './components/SidebarLeft.vue'
import GraphView from './components/GraphView.vue'
import PropsPanel from './components/PropsPanel.vue'

const graphData = ref<GraphData | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const selectedNode = ref<GraphNode | null>(null)
const focusedNodeId = ref<string | null>(null)
const graphView = ref<InstanceType<typeof GraphView> | null>(null)
const isDark = ref(false)

function toggleDark() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
}

onMounted(async () => {
  try {
    graphData.value = await fetchGraph()
  } catch (e: any) {
    error.value = e.message || 'Failed to load graph data'
  } finally {
    loading.value = false
  }
})

function onNodeSelect(node: GraphNode | null) {
  if (node) {
    selectedNode.value = node
    focusedNodeId.value = node.id
  }
}

function onSidebarSelect(nodeId: string) {
  focusedNodeId.value = nodeId
  const node = graphData.value?.nodes.find(n => n.id === nodeId) || null
  selectedNode.value = node
}

function clearFocus() {
  focusedNodeId.value = null
  selectedNode.value = null
}
</script>

<template>
  <div class="app-container">
    <SidebarLeft
      :files="graphData?.files ?? []"
      :nodes="graphData?.nodes ?? []"
      :selectedId="focusedNodeId"
      @select="onSidebarSelect"
    />

    <main class="main-content">
      <div class="top-bar">
        <div class="stats" v-if="graphData && focusedNodeId">
          Tracing: <strong>{{ focusedNodeId }}</strong>
          <button class="btn-clear" @click="clearFocus" title="Show all">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"/></svg>
            Clear
          </button>
        </div>
        <div class="stats" v-else-if="graphData">
          <strong>{{ graphData.nodes.length }}</strong> resources found · Select one to trace
        </div>
        <div class="stats" v-else>Terraform Tracer</div>
        <div class="actions">
          <button class="theme-toggle" @click="toggleDark" :title="isDark ? 'Light mode' : 'Dark mode'">
            <span class="toggle-knob">
              <svg v-if="!isDark" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="5"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
            </span>
          </button>
          <button class="btn btn-icon" @click="graphView?.fit()" title="Fit view">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M15 3h6v6M9 21H3v-6M21 3l-7 7M3 21l7-7"/></svg>
          </button>
        </div>
      </div>

      <div v-if="loading" class="loading-overlay">
        <div class="spinner"></div>
        <p>Analyzing dependencies...</p>
      </div>

      <div v-else-if="error" class="error-overlay">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 8v4M12 16h.01"/></svg>
        <p>{{ error }}</p>
      </div>

      <GraphView
        v-else-if="graphData"
        ref="graphView"
        :graph="graphData"
        :focusNodeId="focusedNodeId"
        @node-select="onNodeSelect"
      />
    </main>

    <PropsPanel
      :node="selectedNode"
      :environments="graphData?.environments ?? []"
    />
  </div>
</template>

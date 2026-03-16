<script setup lang="ts">
import { computed, ref } from 'vue'
import type { GraphNode } from '../services/api'

const props = defineProps<{
  files: string[]
  nodes: GraphNode[]
  selectedId: string | null
}>()

const emit = defineEmits<{ (e: 'select', nodeId: string): void }>()

const search = ref('')

// Group nodes by file
const fileTree = computed(() => {
  const tree: Record<string, GraphNode[]> = {}
  for (const node of props.nodes) {
    const file = node.file || 'unknown'
    if (!tree[file]) tree[file] = []
    tree[file].push(node)
  }
  return tree
})

const filteredFiles = computed(() => {
  const q = search.value.toLowerCase()
  if (!q) return Object.keys(fileTree.value)
  return Object.keys(fileTree.value).filter(f =>
    f.toLowerCase().includes(q) ||
    fileTree.value[f].some(n => n.label.toLowerCase().includes(q))
  )
})

function filteredNodes(file: string): GraphNode[] {
  const q = search.value.toLowerCase()
  const nodes = fileTree.value[file] || []
  if (!q) return nodes
  return nodes.filter(n =>
    n.label.toLowerCase().includes(q) ||
    n.id.toLowerCase().includes(q) ||
    n.provider.toLowerCase().includes(q)
  )
}

function selectNode(nodeId: string) {
  emit('select', nodeId)
}
</script>

<template>
  <aside class="sidebar-left">
    <div class="sidebar-header">
      <div class="logo">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 2L2 7l10 5 10-5-10-5z"/>
          <path d="M2 17l10 5 10-5"/>
          <path d="M2 12l10 5 10-5"/>
        </svg>
        <span>Tf Tracer</span>
      </div>
    </div>

    <div class="search-box">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
      </svg>
      <input v-model="search" type="text" placeholder="Filter resources..." />
    </div>

    <div class="file-tree">
      <div v-for="file in filteredFiles" :key="file" class="tree-item">
        <div class="tree-label tree-file-label">
          <svg class="tree-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
          </svg>
          <span>{{ file }}</span>
        </div>
        <div class="tree-children">
          <div
            v-for="node in filteredNodes(file)"
            :key="node.id"
            class="tree-label"
            :class="{ active: selectedId === node.id }"
            @click="selectNode(node.id)"
          >
            <svg class="tree-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="3" width="18" height="18" rx="2"/>
              <path d="M3 9h18"/>
            </svg>
            <span>{{ node.label }}</span>
          </div>
        </div>
      </div>

      <div v-if="filteredFiles.length === 0" class="tree-empty">
        No resources match "{{ search }}"
      </div>
    </div>
  </aside>
</template>

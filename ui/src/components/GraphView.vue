<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch, computed } from 'vue'
import { Network, DataSet } from 'vis-network/standalone'
import type { GraphData, GraphNode } from '../services/api'

const props = defineProps<{
  graph: GraphData
  focusNodeId: string | null
}>()
const emit = defineEmits<{ (e: 'node-select', node: GraphNode | null): void }>()

const container = ref<HTMLDivElement>()
let network: Network | null = null
let isRebuilding = false

// Compute the neighborhood: only the focused node + its direct connections
const visibleData = computed(() => {
  const { nodes, edges } = props.graph
  const focusId = props.focusNodeId

  if (!focusId) {
    return { nodes: [], edges: [] }
  }

  const connectedEdges = edges.filter(e => e.from === focusId || e.to === focusId)
  const connectedNodeIds = new Set<string>([focusId])
  for (const e of connectedEdges) {
    connectedNodeIds.add(e.from)
    connectedNodeIds.add(e.to)
  }

  const visibleNodes = nodes.filter(n => connectedNodeIds.has(n.id))
  return { nodes: visibleNodes, edges: connectedEdges }
})

function buildNetwork() {
  if (!container.value) return

  isRebuilding = true

  const { nodes: visNodes, edges: visEdges } = visibleData.value

  const nodes = new DataSet(
    visNodes.map(n => ({
      id: n.id,
      label: n.label,
      group: n.type,
      title: `${n.type}: ${n.label}\nFile: ${n.file}`,
      borderWidth: n.id === props.focusNodeId ? 3 : 2,
      color: n.id === props.focusNodeId ? {
        border: '#a30eff',
        background: '#f3e8ff',
        highlight: { border: '#a30eff', background: '#ede9fe' },
        hover: { border: '#a30eff', background: '#f5f0ff' },
      } : undefined,
      font: n.id === props.focusNodeId ? {
        face: 'Inter, system-ui, sans-serif',
        size: 14,
        color: '#7c3aed',
        bold: true as unknown as string,
      } : undefined,
    }))
  )

  const edges = new DataSet(
    visEdges.map((e, i) => ({
      id: `e-${i}`,
      from: e.from,
      to: e.to,
      arrows: 'to',
      dashes: e.edgeType === 'implicit',
      label: e.label || undefined,
      font: { size: 10, color: '#94a3b8', align: 'middle' as const },
    }))
  )

  const options = {
    nodes: {
      shape: 'box',
      fixed: false,
      margin: { top: 10, bottom: 10, left: 14, right: 14 },
      borderWidth: 2,
      borderWidthSelected: 3,
      color: {
        border: '#e2e8f0',
        background: '#ffffff',
        highlight: { border: '#a30eff', background: '#f5f0ff' },
        hover: { border: '#c4b5fd', background: '#faf5ff' },
      },
      font: {
        face: 'Inter, system-ui, sans-serif',
        size: 13,
        color: '#1f2937',
      },
      shadow: {
        enabled: true,
        color: 'rgba(0,0,0,0.04)',
        size: 8,
        x: 0,
        y: 3,
      },
    },
    edges: {
      width: 1.5,
      color: { color: '#cbd5e1', highlight: '#a30eff', hover: '#94a3b8' },
      smooth: { enabled: true, type: 'cubicBezier', forceDirection: 'horizontal', roundness: 0.5 },
    },
    groups: {
      module: {
        shape: 'hexagon',
        color: { border: '#3b82f6', background: '#eff6ff' },
        font: { color: '#1d4ed8' },
        borderWidth: 2,
        shadow: true,
      },
      resource: {
        shape: 'box',
        color: { border: '#e2e8f0', background: '#ffffff' },
      },
      data: {
        shape: 'ellipse',
        color: { border: '#f59e0b', background: '#fffbeb' },
        font: { color: '#92400e' },
      },
    },
    layout: {
      hierarchical: {
        enabled: visNodes.length > 1,
        direction: 'LR',
        sortMethod: 'directed',
        nodeSpacing: 120,
        levelSeparation: 250,
      },
    },
    physics: false,
    interaction: { hover: true, tooltipDelay: 150, zoomView: true, dragNodes: true, dragView: true },
  }

  // Destroy previous network (suppress deselectNode during rebuild)
  if (network) {
    network.destroy()
    network = null
  }

  network = new Network(container.value, { nodes, edges }, options)

  network.on('selectNode', (params: any) => {
    const nodeId = params.nodes[0]
    const node = props.graph.nodes.find(n => n.id === nodeId) || null
    emit('node-select', node)
  })

  network.on('deselectNode', () => {
    // Don't clear selection during rebuild or if we have a focused node
    if (!isRebuilding) {
      emit('node-select', null)
    }
  })

  isRebuilding = false

  // Auto-fit after rendering
  setTimeout(() => {
    network?.fit({ animation: { duration: 400, easingFunction: 'easeInOutQuad' } })
  }, 100)
}

function fit() {
  network?.fit({ animation: { duration: 400, easingFunction: 'easeInOutQuad' } })
}

defineExpose({ fit })

onMounted(buildNetwork)

onUnmounted(() => {
  if (network) {
    network.destroy()
    network = null
  }
})

watch(() => props.focusNodeId, buildNetwork)
watch(() => props.graph, buildNetwork)
</script>

<template>
  <div class="graph-wrapper" style="position: relative; width: 100%; height: 100%;">
    <div id="network-container" ref="container" style="width: 100%; height: 100%; outline: none;"></div>
    
    <div v-if="!focusNodeId" class="graph-prompt">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.3">
        <path d="M12 2L2 7l10 5 10-5-10-5z"/>
        <path d="M2 17l10 5 10-5"/>
        <path d="M2 12l10 5 10-5"/>
      </svg>
      <p>Select a resource from the sidebar to trace its dependencies</p>
    </div>
  </div>
</template>

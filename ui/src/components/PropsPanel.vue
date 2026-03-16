<script setup lang="ts">
import { ref, computed } from 'vue'
import type { GraphNode, Environment } from '../services/api'

const props = defineProps<{
  node: GraphNode | null
  environments: Environment[]
}>()

const activeTab = ref<'envs' | 'props'>('envs')

function envBadgeClass(name: string): string {
  const lower = name.toLowerCase()
  if (lower.includes('prod')) return 'prod'
  if (lower.includes('stg') || lower.includes('staging')) return 'stg'
  if (lower.includes('dev')) return 'dev'
  return 'default'
}

// Only show environment variables that are actually used in the selected resource
const filteredEnvironments = computed(() => {
  if (!props.node || !props.environments) return []
  
  // Se o backend extraiu dependências explícitas (novo comportamento)
  const usedVars = props.node.variables || []
  
  // Fallback pra varrer as props brutas caso o node não traga variables listadas
  const nodeProps = props.node.properties || {}
  const propValues = Object.values(nodeProps).join(' ')

  return props.environments.map(env => {
    const filteredValues: Record<string, any> = {}
    let hasValues = false

    for (const [key, val] of Object.entries(env.values)) {
      // Verifica se está explicitamente na lista de variáveis OU num fallback de texto brute
      if (usedVars.includes(key) || propValues.includes(`var.${key}`) || propValues.includes(key)) {
        filteredValues[key] = val
        hasValues = true
      }
    }

    return {
      ...env,
      values: filteredValues,
      hasValues
    }
  }).filter(env => env.hasValues)
})
</script>

<template>
  <aside class="sidebar-right">
    <!-- Resource selected -->
    <div v-if="node" class="panel-inner">
      <div class="props-header">
        <h3>{{ node.label }}</h3>
      </div>

      <div class="props-content">
        <div class="info-card">
          <div class="info-row">
            <span class="label">Type</span>
            <span class="value">{{ node.type }}</span>
          </div>
          <div class="info-row">
            <span class="label">Provider</span>
            <span class="value">{{ node.provider }}</span>
          </div>
          <div class="info-row">
            <span class="label">File</span>
            <span class="value">{{ node.file }}</span>
          </div>
          <div class="info-row">
            <span class="label">Line</span>
            <span class="value">{{ node.lineStart }}</span>
          </div>
        </div>

        <div class="tabs">
          <button class="tab" :class="{ active: activeTab === 'envs' }" @click="activeTab = 'envs'">
            Environments
          </button>
          <button class="tab" :class="{ active: activeTab === 'props' }" @click="activeTab = 'props'">
            Properties
          </button>
        </div>

        <!-- Environments -->
        <div v-if="activeTab === 'envs'" class="env-list">
          <div v-if="filteredEnvironments.length === 0" style="font-size: 0.8125rem; color: var(--text-muted);">
            No environment variables are referenced by this resource.
          </div>
          <div v-for="env in filteredEnvironments" :key="env.name" class="env-item">
            <div class="env-header">
              <span class="env-badge" :class="envBadgeClass(env.name)">{{ env.name }}</span>
              <span class="env-file">{{ env.filePath }}</span>
            </div>
            <div class="env-vars">
              <div v-for="(val, key) in env.values" :key="key" class="var-item">
                <span>{{ key }}</span>
                <code>{{ val }}</code>
              </div>
            </div>
          </div>
        </div>

        <!-- Properties -->
        <div v-if="activeTab === 'props'">
          <div v-if="!node.properties || Object.keys(node.properties).length === 0"
               style="font-size: 0.8125rem; color: var(--text-muted);">
            No extracted properties.
          </div>
          <div v-else class="code-block">
            <div v-for="(val, key) in node.properties" :key="key">
              <span style="color: #93c5fd;">{{ key }}</span> = <span style="color: #86efac;">"{{ val }}"</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- No resource selected -->
    <div v-else class="panel-empty">
      <div class="props-header" style="width: 100%; border-bottom: 1px solid var(--border-color); padding: 16px; position: absolute; top: 0; left: 0;">
        <h3 style="font-size: 1rem; font-weight: 600; color: var(--text-primary);">Details</h3>
      </div>
      <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.3">
        <path d="M15 3h6v6M9 21H3v-6M21 3l-7 7M3 21l7-7"/>
      </svg>
      <p>Select a resource to see its details and environment variables</p>
    </div>
  </aside>
</template>

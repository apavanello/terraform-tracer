export interface GraphNode {
  id: string
  type: string
  provider: string
  name: string
  label: string
  file: string
  lineStart: number
  properties: Record<string, string>
  variables?: string[]
}

export interface GraphEdge {
  from: string
  to: string
  edgeType: string
  label: string
}

export interface Environment {
  name: string
  filePath: string
  values: Record<string, string>
}

export interface Variable {
  name: string
  type: string
  default: string
}

export interface GraphData {
  nodes: GraphNode[]
  edges: GraphEdge[]
  environments: Environment[]
  variables: Variable[]
  files: string[]
}

const API_BASE = import.meta.env.DEV ? 'http://localhost:8080' : ''

export async function fetchGraph(): Promise<GraphData> {
  const res = await fetch(`${API_BASE}/api/v1/graph`)
  if (!res.ok) {
    throw new Error(`API error: ${res.status} ${res.statusText}`)
  }
  return res.json()
}

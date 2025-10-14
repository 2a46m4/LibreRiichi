export interface AgentInfo {
    name: string
    id: string
    order: number
}

export function sort_agents(agents: AgentInfo[]) {
    agents.sort((a, b) => a.order - b.order)
}

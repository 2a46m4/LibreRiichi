export interface AgentInfo {
    name: string
    id: string
    // The agent's numeric ID
    order: number
}

export function sort_agents(agents: AgentInfo[]) {
    return agents.sort((a, b) => a.order - b.order)
}

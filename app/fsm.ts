export type CallbackArgs = any[]

// Branded types for compile-time type safety
export type StateName<T extends string> = T & { readonly __stateBrand: unique symbol }
export type EventName<T extends string> = T & { readonly __eventBrand: unique symbol } & { readonly toString(): T }

// Helper functions to create branded names
export function createStateName<T extends string>(name: T): StateName<T> {
    return name as StateName<T>
}

export function createEventName<T extends string>(name: T): EventName<T> {
    return name as EventName<T>
}

// Helper to create typed event constants with their arguments
export function createTypedEvent<T extends string, Args extends CallbackArgs>(
    name: T,
    args: Args
): { name: EventName<T>; args: Args } {
    return {
        name: createEventName(name),
        args
    }
}

export type EventCallback<
    FromData extends any,
    ToData extends any,
    Args extends CallbackArgs
> = (from: FromData, to: ToData, ...args: Args) => void | Promise<void>

export type GuardCallback<
    FromData extends any,
    ToData extends any,
    Args extends CallbackArgs
> = (from: FromData, to: ToData, ...args: Args) => boolean

export type TransitionType = 'immediate' | 'deferred'

/**
 * Transition Type Explanation:
 *
 * 'immediate': The transition happens immediately, and any events that occur during
 *              the transition are queued and processed after the transition completes.
 *
 * 'deferred': The transition itself is queued and happens later (useful for animation-driven
 *             transitions or when you want to batch multiple operations). You must manually
 *             trigger the transition by calling process_deferred_transitions().
 */

// Enhanced type-safe state interface
export interface State<Data = {}> {
    name: StateName<string>
    data: Data
    on_enter?: (from_state: StateName<any>, from_data: any) => void
    on_exit?: (to_state: StateName<any>, to_data: any) => void
}

// Type-safe event interface with explicit state constraints
export interface Event<
    FromState extends State<any>,
    ToState extends State<any>,
    Args extends CallbackArgs = []
> {
    name: EventName<string>
    from: FromState
    to: ToState
    callback?: EventCallback<FromState['data'], ToState['data'], Args>
    guard?: GuardCallback<FromState['data'], ToState['data'], Args>
    transition_type?: TransitionType
}

// Type-safe configuration that enforces event-state relationships
export interface FSMConfig<
    TStates extends readonly State<any>[],
    TEvents extends readonly Event<any, any, any>[]
> {
    states: TStates
    events: TEvents
    initial_state: TStates[number]
}

// Type-safe event map for better autocomplete
export type EventMap<TEvents extends readonly Event<any, any, any>[]> = {
    [K in TEvents[number]['name']]: Extract<TEvents[number], { name: K }>
}

// Type-safe event trigger map that links event names to their argument types
export type EventTriggerMap<TEvents extends readonly Event<any, any, any>[]> = {
    [K in TEvents[number]['name']]: Extract<TEvents[number], { name: K }> extends Event<any, any, infer Args>
        ? Args
        : never
}

// Helper to get the argument types for a specific event name
export type EventArguments<TEvents extends readonly Event<any, any, any>[], EventName extends string> =
    Extract<TEvents[number], { name: EventName }> extends Event<any, any, infer Args>
        ? Args
        : never

// Helper to get argument types by branded event name
export type EventArgumentsByBrandedName<TEvents extends readonly Event<any, any, any>[], T extends string> =
    Extract<TEvents[number], { name: EventName<T> }> extends Event<any, any, infer Args>
        ? Args
        : never

export class FSMError extends Error {
    constructor(message: string, public readonly from_state?: string, public readonly event_name?: string) {
        super(message)
        this.name = 'FSMError'
    }
}

// Type-safe FSM class with enhanced constraints
export class FSM<
    TStates extends readonly State<any>[],
    TEvents extends readonly Event<any, any, any>[]
> {
    private current_state: TStates[number]
    private deferred_queue: Array<{
        event: TEvents[number];
        args: any[];
        reason: 'during_transition' | 'explicit_deferred'
    }> = []
    private transition_in_progress = false

    constructor(private config: FSMConfig<TStates, TEvents>) {
        this.validate_configuration()
        this.current_state = config.initial_state
    }

    /**
     * Get the current state
     */
    get_current_state(): TStates[number] {
        return this.current_state
    }

    /**
     * Get the current state name
     */
    get_current_state_name(): string {
        return this.current_state.name
    }

    /**
     * Get the current state data
     */
    get_current_state_data(): TStates[number]['data'] {
        return this.current_state.data
    }

    /**
     * Check if the FSM is in a specific state
     */
    is_state<T extends TStates[number]>(state_name: T['name']): boolean {
        return this.current_state.name === state_name
    }

    /**
     * Type-safe event triggering with compile-time type checking using branded types
     * This provides true compile-time type safety for event arguments
     */
    async trigger_event<T extends string>(
        event_name: EventName<T>,
        ...args: EventArgumentsByBrandedName<TEvents, T>
    ): Promise<void> {
        // Find the event first to check if it's deferred
        const event = this.find_event_by_branded_name(event_name)
        if (!event) {
            throw new FSMError(`Event '${event_name}' not found from current state '${this.current_state.name}'`, this.current_state.name, event_name as string)
        }

        // Handle deferred transitions first (before checking transition_in_progress)
        if (event.transition_type === 'deferred') {
            this.deferred_queue.push({ event, args: args as any[], reason: 'explicit_deferred' })
            return
        }

        // Handle events during transitions
        if (this.transition_in_progress) {
            this.deferred_queue.push({ event, args: args as any[], reason: 'during_transition' })
            return
        }

        // Process all deferred items first to ensure FSM is in correct state
        await this.process_deferred_queue()

        // Type-safe state validation (re-check after processing deferred items)
        if (event.from.name !== this.current_state.name) {
            throw new FSMError(
                `Cannot trigger event '${String(event_name)}' from state '${this.current_state.name}'. This event can only be triggered from state '${event.from.name}'`,
                this.current_state.name,
                String(event_name)
            )
        }

        // Check guard condition if present
        if (event.guard && !event.guard(this.current_state.data, event.to.data, ...(args as any))) {
            return // Transition blocked by guard
        }

        // Execute immediate transition
        await this.perform_transition(event, args as any)
    }

    /**
     * Perform the actual state transition
     */
    private async perform_transition<T extends TEvents[number]>(
        event: T,
        args: T extends Event<any, any, infer Args> ? Args : never
    ): Promise<void> {
        this.transition_in_progress = true

        const from_state = this.current_state
        const to_state = event.to

        try {
            // Call exit callback for current state
            if (from_state.on_exit) {
                from_state.on_exit(to_state.name, to_state.data)
            }

            // Call event callback if present
            if (event.callback) {
                await event.callback(from_state.data, to_state.data, ...(args as any))
            }

            // Update current state
            this.current_state = to_state

            // Call enter callback for new state
            if (to_state.on_enter) {
                to_state.on_enter(from_state.name, from_state.data)
            }

            // Note: Deferred events are now processed before each new event in trigger_event()
            // This ensures the FSM is always in the correct state before attempting transitions

        } catch (error) {
            // If transition fails, rollback
            this.current_state = from_state
            throw new FSMError(
                `Transition failed: ${error instanceof Error ? error.message : 'Unknown error'}`,
                from_state.name,
                event.name
            )
        } finally {
            this.transition_in_progress = false
        }
    }

    
    /**
     * Type-safe event finding
     */
    private find_event<T extends TEvents[number]>(event_name: T['name']): T | undefined {
        return this.config.events.find(event =>
            event.name === event_name && event.from.name === this.current_state.name
        ) as T | undefined
    }

    /**
     * Find event by name (for type-safe triggering)
     */
    private find_event_by_name(event_name: string): TEvents[number] | undefined {
        return this.config.events.find(event =>
            event.name === event_name && event.from.name === this.current_state.name
        )
    }

    /**
     * Find event by branded name (for type-safe triggering)
     */
    private find_event_by_branded_name<T extends string>(event_name: EventName<T>): TEvents[number] | undefined {
        return this.config.events.find(event =>
            event.name === event_name && event.from.name === this.current_state.name
        )
    }

    /**
     * Get all possible events from the current state
     */
    get_available_events(): string[] {
        return this.config.events
            .filter(event => event.from.name === this.current_state.name)
            .map(event => event.name)
    }

    /**
     * Type-safe event trigger check
     */
    can_trigger_event<T extends TEvents[number]>(event_name: T['name']): boolean {
        const event = this.find_event<T>(event_name)
        if (!event) return false

        // Check guard condition if present
        if (event.guard) {
            return event.guard(this.current_state.data, event.to.data)
        }

        return true
    }

    /**
     * Process all deferred items from the unified queue
     */
    async process_deferred_queue(options?: {
        only_explicit_deferred?: boolean;
        only_during_transition?: boolean;
    }): Promise<void> {
        const queue = [...this.deferred_queue]
        this.deferred_queue = []

        for (const { event, args, reason } of queue) {
            // Apply filtering if specified
            if (options?.only_explicit_deferred && reason !== 'explicit_deferred') continue
            if (options?.only_during_transition && reason !== 'during_transition') continue

            try {
                // Re-check that the transition is still valid from current state
                if (event.from.name === this.current_state.name) {
                    // Re-check guard condition
                    if (!event.guard || event.guard(this.current_state.data, event.to.data, ...(args as any))) {
                        await this.perform_transition(event, args as any)
                    }
                }
            } catch (error) {
                console.error(`Deferred event '${event.name}' (reason: ${reason}) failed:`, error)
                // Continue processing remaining items even if one fails
            }
        }
    }

    /**
     * Process only explicit deferred transitions (legacy method)
     */
    async process_deferred_transitions(): Promise<void> {
        await this.process_deferred_queue({ only_explicit_deferred: true })
    }

    /**
     * Process only items queued during transitions (legacy method)
     */
    async process_deferred_events(): Promise<void> {
        await this.process_deferred_queue({ only_during_transition: true })
    }

    /**
     * Get deferred queue information
     */
    get_deferred_info(): {
        has_items: boolean;
        total_count: number;
        explicit_deferred: number;
        during_transition: number;
    } {
        const counts = {
            explicit_deferred: 0,
            during_transition: 0,
            total: this.deferred_queue.length
        }

        for (const item of this.deferred_queue) {
            if (item.reason === 'explicit_deferred') counts.explicit_deferred++
            else if (item.reason === 'during_transition') counts.during_transition++
        }

        return {
            has_items: this.deferred_queue.length > 0,
            total_count: this.deferred_queue.length,
            explicit_deferred: counts.explicit_deferred,
            during_transition: counts.during_transition
        }
    }

    
    /**
     * Type-safe state data update
     */
    update_state_data<T extends TStates[number]>(state_name: T['name'], data: Partial<T['data']>): void {
        if (this.current_state.name === state_name) {
            this.current_state.data = { ...this.current_state.data, ...data } as any
        } else {
            throw new FSMError(`Cannot update data for state '${state_name}' when current state is '${this.current_state.name}'`)
        }
    }

    /**
     * Reset to initial state
     */
    reset(): void {
        this.current_state = this.config.initial_state
        this.deferred_queue = []
        this.transition_in_progress = false
    }

    /**
     * Validate the FSM configuration
     */
    private validate_configuration(): void {
        // Check that all states in events exist in the states array
        const state_names = new Set(this.config.states.map(state => state.name))

        this.config.events.forEach(event => {
            if (!state_names.has(event.from.name)) {
                throw new FSMError(`State '${event.from.name}' referenced in event '${event.name}' does not exist in states array`)
            }
            if (!state_names.has(event.to.name)) {
                throw new FSMError(`State '${event.to.name}' referenced in event '${event.name}' does not exist in states array`)
            }
        })

        // Check that initial state exists
        if (!state_names.has(this.config.initial_state.name)) {
            throw new FSMError(`Initial state '${this.config.initial_state.name}' does not exist in states array`)
        }
    }

    /**
     * Get a visual representation of the FSM
     */
    get_debug_info(): {
        current_state: string
        states: string[]
        events: Array<{ name: string; from: string; to: string }>
        available_events: string[]
        deferred_info: {
            has_items: boolean;
            total_count: number;
            explicit_deferred: number;
            during_transition: number;
        }
    } {
        return {
            current_state: this.current_state.name,
            states: this.config.states.map(state => state.name),
            events: this.config.events.map(event => ({
                name: event.name,
                from: event.from.name,
                to: event.to.name
            })),
            available_events: this.get_available_events(),
            deferred_info: this.get_deferred_info()
        }
    }
}

/**
 * Type-safe factory function to create an FSM with better type inference
 */
export function create_fsm<
    TStates extends readonly State<any>[],
    TEvents extends readonly Event<any, any, any>[]
>(config: FSMConfig<TStates, TEvents>): FSM<TStates, TEvents> {
    return new FSM(config)
}

/**
 * Type-safe state creation helper
 */
export function create_state<Data>(name: string, data: Data, options?: {
    on_enter?: (from_state: StateName<any>, from_data: any) => void
    on_exit?: (to_state: StateName<any>, to_data: any) => void
}): State<Data> {
    return {
        name: createStateName(name),
        data,
        on_enter: options?.on_enter,
        on_exit: options?.on_exit
    }
}

/**
 * Type-safe event creation helper
 */
export function create_event<
    FromState extends State<any>,
    ToState extends State<any>,
    Args extends CallbackArgs = []
>(
    name: string,
    from: FromState,
    to: ToState,
    options?: {
        callback?: EventCallback<FromState['data'], ToState['data'], Args>
        guard?: GuardCallback<FromState['data'], ToState['data'], Args>
        transition_type?: TransitionType
    }
): Event<FromState, ToState, Args> {
    return {
        name: createEventName(name),
        from,
        to,
        callback: options?.callback,
        guard: options?.guard,
        transition_type: options?.transition_type
    }
}

/**
 * Type-safe builder pattern for FSM creation
 */
export class FSMBuilder<
    TStates extends readonly State<any>[] = [],
    TEvents extends readonly Event<any, any, any>[] = []
> {
    private states: TStates = [] as unknown as TStates
    private events: TEvents = [] as unknown as TEvents

    add_state<Data>(state: State<Data>): FSMBuilder<[...TStates, State<Data>], TEvents> {
        const new_states = [...(this.states as readonly State<any>[]), state] as [...TStates, State<Data>]
        const builder = new FSMBuilder<[...TStates, State<Data>], TEvents>()
        builder.states = new_states
        builder.events = this.events
        return builder as any
    }

    add_event<FromState extends TStates[number], ToState extends TStates[number], Args extends CallbackArgs = []>(
        event: Event<FromState, ToState, Args>
    ): FSMBuilder<TStates, [...TEvents, Event<FromState, ToState, Args>]> {
        const new_events = [...(this.events as readonly Event<any, any, any>[]), event] as [...TEvents, Event<FromState, ToState, Args>]
        const builder = new FSMBuilder<TStates, [...TEvents, Event<FromState, ToState, Args>]>()
        builder.states = this.states
        builder.events = new_events
        return builder as any
    }

    build<InitialState extends TStates[number]>(initial_state: InitialState): FSM<TStates, TEvents> {
        if (this.states.length === 0) {
            throw new FSMError('At least one state must be added before building FSM')
        }

        return new FSM({
            states: this.states,
            events: this.events,
            initial_state
        })
    }
}

/**
 * Helper to create a new FSM builder
 */
export function create_fsm_builder(): FSMBuilder<[], []> {
    return new FSMBuilder()
}

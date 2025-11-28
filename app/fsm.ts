export type EventCallback<
    FromData extends any,
    ToData extends any,
    Args extends any[]
> = (from: FromData, to: ToData, ...args: Args) => void | Promise<void>

export type GuardCallback<
    FromData extends any,
    ToData extends any,
    Args extends any[]
> = (from: FromData, to: ToData, ...args: Args) => boolean

export type FindEvent<NameToMatch, TEvents> =
    TEvents extends readonly [infer First, ...infer Rest]
    ? First extends Event<infer EventName, any, any>
    ? EventName extends NameToMatch
    ? First
    : Rest extends readonly Event<string, any, any>[]
    ? FindEvent<NameToMatch, Rest>
    : never
    : never
    : never

export type ExtractCallbackArgs<T> = T extends Event<any, any, any, infer Args>
    ? Args
    : never

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
export type TransitionType = 'immediate' | 'deferred'

// Enhanced type-safe state interface
export interface State<Name extends string, Data = {}, FromName extends string = string, ToName extends string = string> {
    name: Name
    data: Data
    on_enter?: (from_state: State<FromName>) => void
    on_exit?: (to_state: State<ToName>) => void
}

// Type-safe event interface with explicit state constraints
export interface Event<
    Name extends string,
    FromState extends State<string>,
    ToState extends State<string>,
    Args extends any[] = any[]
> {
    name: Name
    from: FromState
    to: ToState
    callback?: EventCallback<FromState['data'], ToState['data'], Args>
    guard?: GuardCallback<FromState['data'], ToState['data'], Args>
    transition_type?: TransitionType
}

// Type-safe configuration that enforces event-state relationships
export interface FSMConfig<
    TStates extends readonly State<string>[],
    TEvents extends readonly Event<string, TStates[number], TStates[number]>[]
> {
    states: TStates
    events: TEvents
    initial_state: TStates[number]
}

// Type-safe FSM class with enhanced constraints
export class FSM<
    TStates extends readonly State<string>[],
    TEvents extends readonly Event<string, TStates[number], TStates[number]>[]
> {
    private current_state: TStates[number]
    private deferred_queue: Array<{
        event: TEvents[number];
        args: any[];
        reason: 'during_transition' | 'explicit_deferred'
    }> = []
    private transition_in_progress = false

    constructor(public config: FSMConfig<TStates, TEvents>) {
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
     * Trigger an event by name with type-safe arguments
     */
    async trigger_event<T extends string, Args extends ExtractCallbackArgs<FindEvent<T, TEvents>>>(event_name: T, ...args: Args): Promise<void> {
        // Find the event that can be triggered from current state
        const found_event = this.config.events.find(e =>
            e.name === event_name && e.from.name === this.current_state.name
        )
        if (!found_event) {
            throw new Error(`Event '${String(event_name)}' cannot be triggered from current state '${this.current_state.name}'. Available events: ${this.get_available_events().join(', ')}`)
        }

        // Handle deferred transitions first (before checking transition_in_progress)
        if (found_event.transition_type === 'deferred') {
            this.deferred_queue.push({ event: found_event, args: args as any[], reason: 'explicit_deferred' })
            return
        }

        // Handle events during transitions
        if (this.transition_in_progress) {
            this.deferred_queue.push({ event: found_event, args: args as any[], reason: 'during_transition' })
            return
        }

        // Process all deferred items first to ensure FSM is in correct state
        await this.process_deferred_queue()

        // Type-safe state validation (re-check after processing deferred items)
        if (found_event.from.name !== this.current_state.name) {
            throw new Error(
                `Cannot trigger event '${String(event_name)}' from state '${this.current_state.name}'. This event can only be triggered from state '${found_event.from.name}'`
            )
        }

        // Check guard condition if present
        if (found_event.guard && !found_event.guard(this.current_state.data, found_event.to.data, ...args)) {
            return // Transition blocked by guard
        }

        // Execute immediate transition
        await this.perform_transition(found_event, args as any)
    }

    /**
     * Perform the actual state transition
     */
    private async perform_transition<T extends TEvents[number]>(
        event: T,
        args: T extends Event<string, any, any, infer Args> ? Args : never
    ): Promise<void> {
        this.transition_in_progress = true

        const from_state = this.current_state
        const to_state = event.to

        try {
            // Call exit callback for current state
            if (from_state.on_exit) {
                from_state.on_exit(to_state)
            }

            // Call event callback if present
            if (event.callback) {
                await event.callback(from_state.data, to_state.data, ...args)
            }

            // Update current state
            this.current_state = to_state

            // Call enter callback for new state
            if (to_state.on_enter) {
                to_state.on_enter(from_state)
            }

            // Note: Deferred events are now processed before each new event in trigger_event()
            // This ensures the FSM is always in the correct state before attempting transitions

        } catch (error) {
            // If transition fails, rollback
            this.current_state = from_state
            throw new Error(
                `Transition failed: ${error instanceof Error ? error.message : 'Unknown error'}`
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
        stop_on_error?: boolean;
    }): Promise<{
        processed: number;
        failed: Array<{
            event: string;
            reason: string;
            error: Error;
        }>;
    }> {
        const queue = [...this.deferred_queue]
        this.deferred_queue = []

        const results = {
            processed: 0,
            failed: [] as Array<{ event: string; reason: string; error: Error }>
        }

        for (const { event, args, reason } of queue) {
            // Apply filtering if specified
            if (options?.only_explicit_deferred && reason !== 'explicit_deferred') continue
            if (options?.only_during_transition && reason !== 'during_transition') continue

            try {
                // Re-check that the transition is still valid from current state
                if (event.from.name === this.current_state.name) {
                    // Re-check guard condition
                    if (!event.guard || event.guard(this.current_state.data, event.to.data, ...args)) {
                        await this.perform_transition(event, args as any)
                        results.processed++
                    }
                }
            } catch (error) {
                const errorInfo = {
                    event: event.name,
                    reason,
                    error: error instanceof Error ? error : new Error(String(error))
                }

                if (options?.stop_on_error) {
                    // Re-queue remaining items for later
                    this.deferred_queue.unshift(...queue.slice(queue.indexOf({ event, args, reason })))
                    throw new Error(`Deferred event processing stopped: ${errorInfo.event} failed: ${errorInfo.error.message}`)
                }

                results.failed.push(errorInfo)
                console.error(`Deferred event '${event.name}' (reason: ${reason}) failed:`, error)
            }
        }

        return results
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
            throw new Error(`Cannot update data for state '${state_name}' when current state is '${this.current_state.name}'`)
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
                throw new Error(`State '${event.from.name}' referenced in event '${event.name}' does not exist in states array`)
            }
            if (!state_names.has(event.to.name)) {
                throw new Error(`State '${event.to.name}' referenced in event '${event.name}' does not exist in states array`)
            }
        })

        // Check that initial state exists
        if (!state_names.has(this.config.initial_state.name)) {
            throw new Error(`Initial state '${this.config.initial_state.name}' does not exist in states array`)
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
    TStates extends readonly State<string>[],
    TEvents extends readonly Event<string, TStates[number], TStates[number]>[]
>(config: FSMConfig<TStates, TEvents>): FSM<TStates, TEvents> {
    return new FSM(config)
}

/**
 * Type-safe state creation helper
 */
export function create_state<
    Name extends string,
    Data,
>(name: Name, data: Data, options?: {
    on_enter?: (from_state: State<string>) => void
    on_exit?: (to_state: State<string>) => void
}): State<Name, Data> {
    return {
        name: name,
        data,
        on_enter: options?.on_enter,
        on_exit: options?.on_exit
    } as const
}

/**
 * Type-safe event creation helper with proper argument type inference
 */
export function create_event<
    Name extends string,
    FromState extends State<string>,
    ToState extends State<string>,
    Args extends any[]
>(
    name: Name,
    from: FromState,
    to: ToState,
    options: {
        callback: (from: FromState['data'], to: ToState['data'], ...args: Args) => void | Promise<void>
        guard?: (from: FromState['data'], to: ToState['data'], ...args: Args) => boolean
        transition_type?: TransitionType
    }
): Event<Name, FromState, ToState, Args> {
    return {
        name: name,
        from: from,
        to: to,
        callback: options?.callback,
        guard: options?.guard,
        transition_type: options?.transition_type
    } as const
}

export class FSMBuilder<
    const TStates extends readonly State<string>[] = [],
    const TEvents extends readonly Event<string, TStates[number], TStates[number]>[] = []
> {
    private states: TStates = [] as unknown as TStates
    private events: TEvents = [] as unknown as TEvents

    add_state<const Name extends string,
        const NewStates extends readonly [...TStates, State<Name>],
        const NewEventType extends readonly Event<string, NewStates[number], NewStates[number]>[] = TEvents
    >(state: State<Name>): FSMBuilder<NewStates, NewEventType> {
        const new_states = [...this.states, state]
        const builder = new FSMBuilder<NewStates, NewEventType>()
        builder.states = new_states as unknown as NewStates
        return builder
    }

    add_states<Name extends string>(
        ...states: State<Name>[]
    ): FSMBuilder<readonly [...TStates, ...State<Name>[]], []> {
        const new_states = [...this.states, ...states] as unknown as readonly [...TStates, ...State<Name>[]]
        const builder = new FSMBuilder<readonly [...TStates, ...State<Name>[]], []>()
        builder.states = new_states
        return builder
    }

    add_event<Name extends string,
        FromState extends TStates[number],
        ToState extends TStates[number],
        Args extends any[],
        const NewEventType extends readonly [...TEvents, Event<Name, FromState, ToState, Args>]>(
            event: Event<Name, FromState, ToState, Args>
        ): FSMBuilder<TStates, NewEventType> {
        const new_events = [...this.events, event]
        const builder = new FSMBuilder<TStates, NewEventType>()
        builder.states = this.states
        builder.events = new_events as unknown as NewEventType
        return builder
    }

    build<InitialState extends TStates[number]>(initial_state: InitialState): FSM<TStates, TEvents> {
        if (this.states.length === 0) {
            throw new Error('At least one state must be added before building FSM')
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

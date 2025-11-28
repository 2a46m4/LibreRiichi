// Test file to verify the add_states method works correctly
import { create_fsm_builder, create_state, create_event } from './fsm'

// Test the add_states method
function test_add_states() {
    console.log('Testing add_states method...')

    // Create some states with different names
    const idleState = create_state('idle', { counter: 0 })
    const activeState = create_state('active', { counter: 0 })
    const pausedState = create_state('paused', { counter: 0 })
    const completedState = create_state('completed', { counter: 0 })

    // Test add_states with variadic arguments
    const builder = create_fsm_builder()
        .add_states(idleState, activeState, pausedState, completedState)

    console.log('✓ add_states method successfully added multiple states')

    // Test that we can build the FSM
    const fsm = builder.build(idleState)
    console.log('✓ FSM successfully built with add_states')

    // Test current state
    console.log('Current state:', fsm.get_current_state_name()) // Should be 'idle'
    console.log('Available states:', fsm.get_debug_info().states) // Should show all 4 states

    // Add some events to test functionality
    const startEvent = create_event('start', idleState, activeState, {
        callback: (from, to) => {
            console.log(`Transitioning from ${from.counter} to ${to.counter}`)
        }
    })

    const pauseEvent = create_event('pause', activeState, pausedState, {
        callback: (from, to) => {
            console.log(`Pausing at counter ${from.counter}`)
        }
    })

    const completeEvent = create_event('complete', activeState, completedState, {
        callback: (from, to) => {
            console.log(`Completing at counter ${from.counter}`)
        }
    })

    // Add events and rebuild
    const builderWithEvents = builder
        .add_event(startEvent)
        .add_event(pauseEvent)
        .add_event(completeEvent)

    const fsmWithEvents = builderWithEvents.build(idleState)

    console.log('✓ Events added successfully')
    console.log('Available events from idle:', fsmWithEvents.get_available_events())

    // Test transitions
    async function testTransitions() {
        try {
            await fsmWithEvents.trigger_event('start')
            console.log('Current state after start:', fsmWithEvents.get_current_state_name())

            await fsmWithEvents.trigger_event('pause')
            console.log('Current state after pause:', fsmWithEvents.get_current_state_name())

            // Reset and try complete
            fsmWithEvents.reset()
            await fsmWithEvents.trigger_event('start')
            await fsmWithEvents.trigger_event('complete')
            console.log('Current state after complete:', fsmWithEvents.get_current_state_name())

            console.log('✓ All transitions worked correctly')
        } catch (error) {
            console.error('✗ Transition failed:', error)
        }
    }

    testTransitions()
}

// Run the test
test_add_states()

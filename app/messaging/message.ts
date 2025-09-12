import {ArenaMessage} from "./arena_message";
import { ServerAction, ServerActionType } from "./server_action_generated";
import { ServerResponse, ServerResponseType } from "./server_response_generated";
import { ServerEvent, ServerEventType } from "./server_event_generated";
import {hasOwnProperty} from "@tailwindcss/postcss";

export enum MessageType {
    RESPONSE = 0,
    REQUEST = 1,
    EVENT = 2
}

export interface Message {
    message_type: MessageType;
    message_index: number;
    data: any;
}

export type IncomingMessage = Message;

// Basic type guards
function is_string(value: any): value is string {
    return typeof value === 'string';
}

function is_number(value: any): value is number {
    return typeof value === 'number' && !isNaN(value);
}

function is_boolean(value: any): value is boolean {
    return typeof value === 'boolean';
}

function is_object(value: any): value is object {
    return value !== null && typeof value === 'object' && !Array.isArray(value);
}

function is_array(value: any): value is any[] {
    return Array.isArray(value);
}

// ServerAction validators
function is_initial_message_action(data: any): data is import('./server_action_generated').InitialMessageAction {
    return is_object(data) && "name" in data && is_string(data.name);
}

function is_join_arena_action(data: any): data is import('./server_action_generated').JoinArenaAction {
    return is_object(data) && is_string(data.arena_name);
}

function is_server_arena_action(data: any): data is import('./server_action_generated').ServerArenaAction {
    return is_object(data) && data.arena_action !== undefined;
}

function is_list_arenas_action(data: any): data is import('./server_action_generated').ListArenasAction {
    return is_object(data);
}

function is_create_arena_action(data: any): data is import('./server_action_generated').CreateArenaAction {
    return is_object(data) && is_string(data.arena_name);
}

function is_arena_info_action(data: any): data is import('./server_action_generated').ArenaInfoAction {
    return is_object(data);
}

// ServerResponse validators
function is_generic_response(data: any): data is import('./server_response_generated').GenericResponse {
    return is_object(data) && 
           is_boolean(data.success) && 
           is_string(data.fail_reason);
}

function is_list_arenas_response(data: any): data is import('./server_response_generated').ListArenasResponse {
    return is_object(data) && 
           is_boolean(data.success) && 
           is_array(data.arena_list) && 
           data.arena_list.every((item: any) => is_string(item));
}

function is_arena_info_response(data: any): data is import('./server_response_generated').ArenaInfoResponse {
    return is_object(data) && 
           is_boolean(data.success) && 
           is_string(data.name) && 
           is_array(data.agents) && 
           is_boolean(data.game_started) && 
           is_string(data.date_created);
}

// ServerEvent validators
function is_server_arena_event(data: any): data is import('./server_event_generated').ServerArenaEvent {
    return is_object(data) && data.arena_message !== undefined;
}

// Type validation result
export interface ValidationResult {
    isValid: boolean;
    messageType?: MessageType;
    dataType?: string;
    errors: string[];
}

// Main message validator function
export function validate_message(json: any): ValidationResult {
    const errors: string[] = [];
    
    // Check if json is an object
    if (!is_object(json)) {
        return {
            isValid: false,
            errors: ['Message must be an object']
        };
    }
    
    // Cast to any to access properties
    const obj = json as any;
    
    // Validate message structure
    if (!is_number(obj.message_type)) {
        errors.push('message_type must be a number');
    }
    
    if (!is_number(obj.message_index)) {
        errors.push('message_index must be a number');
    }
    
    if (obj.data === undefined) {
        errors.push('data field is required');
    }
    
    // If basic structure is invalid, return early
    if (errors.length > 0) {
        return {
            isValid: false,
            errors
        };
    }
    
    const messageType = obj.message_type as MessageType;
    const data = obj.data;
    
    // Validate based on message type
    switch (messageType) {
        case MessageType.REQUEST:
            return validate_request_data(data, errors);
        case MessageType.RESPONSE:
            return validate_response_data(data, errors);
        case MessageType.EVENT:
            return validate_event_data(data, errors);
        default:
            return {
                isValid: false,
                errors: [`Unknown message_type: ${messageType}`]
            };
    }
}

function validate_request_data(data: any, errors: string[]): ValidationResult {
    // Check for serveraction_type field to determine specific type
    if (!is_object(data)) {
        return {
            isValid: false,
            messageType: MessageType.REQUEST,
            errors: [...errors, 'Request data must be an object']
        };
    }
    
    const actionType = data.serveraction_type;
    if (!is_number(actionType)) {
        return {
            isValid: false,
            messageType: MessageType.REQUEST,
            errors: [...errors, 'serveraction_type must be a number']
        };
    }
    
    let dataType: string;
    let isValidData: boolean;
    
    switch (actionType) {
        case ServerActionType.InitialMessageAction:
            isValidData = is_initial_message_action(data);
            dataType = 'InitialMessageAction';
            break;
        case ServerActionType.JoinArenaAction:
            isValidData = is_join_arena_action(data);
            dataType = 'JoinArenaAction';
            break;
        case ServerActionType.ServerArenaAction:
            isValidData = is_server_arena_action(data);
            dataType = 'ServerArenaAction';
            break;
        case ServerActionType.ListArenasAction:
            isValidData = is_list_arenas_action(data);
            dataType = 'ListArenasAction';
            break;
        case ServerActionType.CreateArenaAction:
            isValidData = is_create_arena_action(data);
            dataType = 'CreateArenaAction';
            break;
        case ServerActionType.ArenaInfoAction:
            isValidData = is_arena_info_action(data);
            dataType = 'ArenaInfoAction';
            break;
        default:
            return {
                isValid: false,
                messageType: MessageType.REQUEST,
                errors: [...errors, `Unknown ServerActionType: ${actionType}`]
            };
    }
    
    if (!isValidData) {
        errors.push(`Invalid ${dataType} data structure`);
    }
    
    return {
        isValid: errors.length === 0,
        messageType: MessageType.REQUEST,
        dataType,
        errors
    };
}

function validate_response_data(data: any, errors: string[]): ValidationResult {
    if (!is_object(data)) {
        return {
            isValid: false,
            messageType: MessageType.RESPONSE,
            errors: [...errors, 'Response data must be an object']
        };
    }
    
    const responseType = data.serverresponse_type;
    if (!is_number(responseType)) {
        return {
            isValid: false,
            messageType: MessageType.RESPONSE,
            errors: [...errors, 'serverresponse_type must be a number']
        };
    }
    
    let dataType: string;
    let isValidData: boolean;
    
    switch (responseType) {
        case ServerResponseType.GenericResponse:
            isValidData = is_generic_response(data);
            dataType = 'GenericResponse';
            break;
        case ServerResponseType.ListArenasResponse:
            isValidData = is_list_arenas_response(data);
            dataType = 'ListArenasResponse';
            break;
        case ServerResponseType.ArenaInfoResponse:
            isValidData = is_arena_info_response(data);
            dataType = 'ArenaInfoResponse';
            break;
        default:
            return {
                isValid: false,
                messageType: MessageType.RESPONSE,
                errors: [...errors, `Unknown ServerResponseType: ${responseType}`]
            };
    }
    
    if (!isValidData) {
        errors.push(`Invalid ${dataType} data structure`);
    }
    
    return {
        isValid: errors.length === 0,
        messageType: MessageType.RESPONSE,
        dataType,
        errors
    };
}

function validate_event_data(data: any, errors: string[]): ValidationResult {
    if (!is_object(data)) {
        return {
            isValid: false,
            messageType: MessageType.EVENT,
            errors: [...errors, 'Event data must be an object']
        };
    }
    
    const eventType = data.serverevent_type;
    if (!is_number(eventType)) {
        return {
            isValid: false,
            messageType: MessageType.EVENT,
            errors: [...errors, 'serverevent_type must be a number']
        };
    }
    
    let dataType: string;
    let isValidData: boolean;
    
    switch (eventType) {
        case ServerEventType.ServerArenaEvent:
            isValidData = is_server_arena_event(data);
            dataType = 'ServerArenaEvent';
            break;
        default:
            return {
                isValid: false,
                messageType: MessageType.EVENT,
                errors: [...errors, `Unknown ServerEventType: ${eventType}`]
            };
    }
    
    if (!isValidData) {
        errors.push(`Invalid ${dataType} data structure`);
    }
    
    return {
        isValid: errors.length === 0,
        messageType: MessageType.EVENT,
        dataType,
        errors
    };
}

// Example usage and tests
export function test_validator() {
    // Example 1: Valid InitialMessageAction
    const validRequest = {
        message_type: MessageType.REQUEST,
        message_index: 1,
        data: {
            serveraction_type: ServerActionType.InitialMessageAction,
            name: "John"
        }
    };
    
    console.log("Valid InitialMessageAction:", validate_message(validRequest));
    
    // Example 2: Valid GenericResponse
    const validResponse = {
        message_type: MessageType.RESPONSE,
        message_index: 2,
        data: {
            serverresponse_type: ServerResponseType.GenericResponse,
            success: true,
            fail_reason: ""
        }
    };
    
    console.log("Valid GenericResponse:", validate_message(validResponse));
    
    // Example 3: Invalid message - missing required field
    const invalidMessage1 = {
        message_type: MessageType.REQUEST,
        message_index: 3,
        data: {
            serveraction_type: ServerActionType.InitialMessageAction
            // missing 'name' field
        }
    };
    
    console.log("Invalid InitialMessageAction (missing name):", validate_message(invalidMessage1));
    
    // Example 4: Invalid message - wrong structure
    const invalidMessage2 = {
        message_type: "invalid", // should be number
        message_index: 4,
        data: {}
    };
    
    console.log("Invalid message structure:", validate_message(invalidMessage2));
    
    // Example 5: Valid ListArenasResponse
    const validListResponse = {
        message_type: MessageType.RESPONSE,
        message_index: 5,
        data: {
            serverresponse_type: ServerResponseType.ListArenasResponse,
            success: true,
            arena_list: ["arena1", "arena2", "arena3"]
        }
    };
    
    console.log("Valid ListArenasResponse:", validate_message(validListResponse));
    
    // Example 6: Invalid array in response
    const invalidArrayResponse = {
        message_type: MessageType.RESPONSE,
        message_index: 6,
        data: {
            serverresponse_type: ServerResponseType.ListArenasResponse,
            success: true,
            arena_list: ["arena1", 123, "arena3"] // contains number instead of string
        }
    };
    
    console.log("Invalid ListArenasResponse (mixed array):", validate_message(invalidArrayResponse));
}


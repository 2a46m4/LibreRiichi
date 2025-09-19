import {ServerActionMessage, ServerActionType} from "./server_action_generated";
import {ServerResponseMessage, ServerResponseType} from "./server_response_generated";
import {ServerEventMessage, ServerEventType} from "./server_event_generated";

export enum MessageType {
    RESPONSE = 0,
    REQUEST = 1,
    EVENT = 2
}

export type Message = ServerAction | ServerResponse | ServerEvent

export type ServerAction = {
    message_type: MessageType.REQUEST,
    data: ServerActionMessage
}

export type ServerResponse = {
    message_type: MessageType.RESPONSE,
    data: ServerResponseMessage
}

export type ServerEvent = {
    message_type: MessageType.EVENT,
    data: ServerEventMessage
}

export type IncomingMessage = Message & {message_index: number};

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
    return is_object(data) && "arena_name" in data && is_string(data.arena_name);
}

function is_server_arena_action(data: any): data is import('./server_action_generated').ServerArenaAction {
    return is_object(data) && "arena_action" in data && data.arena_action !== undefined;
}

function is_list_arenas_action(data: any): data is import('./server_action_generated').ListArenasAction {
    return is_object(data);
}

function is_create_arena_action(data: any): data is import('./server_action_generated').CreateArenaAction {
    return is_object(data) && "arena_name" in data && is_string(data.arena_name);
}

function is_arena_info_action(data: any): data is import('./server_action_generated').ArenaInfoAction {
    return is_object(data);
}

// ServerResponse validators
function is_generic_response(data: any): data is import('./server_response_generated').GenericResponse {
    return is_object(data) && 
           "success" in data && is_boolean(data.success) && 
           "fail_reason" in data && is_string(data.fail_reason);
}

function is_list_arenas_response(data: any): data is import('./server_response_generated').ListArenasResponse {
    return is_object(data) && 
           "success" in data && is_boolean(data.success) && 
           "arena_list" in data && is_array(data.arena_list) && 
           data.arena_list.every((item: any) => is_string(item));
}

function is_arena_info_response(data: any): data is import('./server_response_generated').ArenaInfoResponse {
    return is_object(data) && 
           "success" in data && is_boolean(data.success) && 
           "name" in data && is_string(data.name) && 
           "agents" in data && is_array(data.agents) && 
           "game_started" in data && is_boolean(data.game_started) && 
           "date_created" in data && is_string(data.date_created);
}

// ServerEvent validators
function is_server_arena_event(data: any): data is import('./server_event_generated').ServerArenaEvent {
    return is_object(data) && "arena_message" in data && data.arena_message !== undefined;
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
    
    const actionType = "serveraction_type" in data ? data.serveraction_type : undefined;
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
    
    const responseType = "serverresponse_type" in data ? data.serverresponse_type : undefined;
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
    
    const eventType = "serverevent_type" in data ? data.serverevent_type : undefined;
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

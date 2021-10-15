import { ActionType } from "../types/type"
import { createContext } from '../utils/createContext'
export type InitialState = typeof initialState;
export type AppReducer = typeof reducer;

const initialState = {
    isAuthenticated: false
}

const reducer = (state: InitialState, action: ActionType): InitialState => {
    switch (action.type) {
        case 'SET_IS_AUTHENTICATED':
            return { ...state, isAuthenticated: action.payload.status }
        default:
            return state
    }
}

const actions = {
    setIsAuth: (dispatch: any) => (status: boolean): void => {
        dispatch({ type: 'SET_IS_AUTHENTICATED', payload: { status } })
    },
};

export type UnboundedActions = typeof actions;
export type getActions<T extends { [key: string]: (dispatch: any) => (...args: any[]) => any }> = {
    [key in keyof T]: ReturnType<T[key]>;
  };
export type BoundedActions = getActions<typeof actions>;

export const { Context, Provider } = createContext(reducer,  actions, initialState)

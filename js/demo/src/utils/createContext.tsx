import React, { FC, useMemo, useReducer } from 'react';
import { BoundedActions, AppReducer, InitialState, UnboundedActions } from '../context/app-context';
type AppContextType = {
    state: InitialState;
    actions: BoundedActions;
};
export const createContext = (reducer: AppReducer, actions: UnboundedActions, initialState: InitialState) => {

    const Context = React.createContext<AppContextType>({ state: initialState, actions });

    const Provider: FC = ({ children }) => {
        const [state, dispatch] = useReducer(reducer, initialState);
        const cachedActions:BoundedActions = useMemo(() => {
            let boundActions: BoundedActions = {} as BoundedActions;
            let key: keyof UnboundedActions;
            for (key in actions) {
                boundActions[key] = actions[key](dispatch) as any;
            };
            return boundActions
        }, [])

        return (
            <Context.Provider value={{ state, actions: cachedActions } as any}>
                {children}
            </Context.Provider>
        )
    }
    return { Context, Provider };
}
import { BehaviorSubject } from 'rxjs';
import IdentifoAuth from '../IdentifoAuth';
import { IdentifoConfig } from '../types/types';
import { StateLogin, States } from './model';
export declare class CDK {
    auth: IdentifoAuth;
    states: {
        login: StateLogin;
    };
    state: BehaviorSubject<States>;
    constructor(authConfig?: IdentifoConfig);
}

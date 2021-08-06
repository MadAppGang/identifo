import '@babel/polyfill';
import React from 'react';
import Enzyme, { shallow } from 'enzyme';
import EnzymeAdapter from 'enzyme-adapter-react-16';

global.React = React;
global.shallow = shallow;

Enzyme.configure({ adapter: new EnzymeAdapter() });

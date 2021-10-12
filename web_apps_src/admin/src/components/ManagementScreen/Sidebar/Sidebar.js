import React from 'react';
import StaticFilesIcon from '~/components/icons/StaticFilesIcon.svg';
import UsersIcon from '~/components/icons/UsersIcon.svg';
import AdminIcon from '~/components/icons/AdminIcon.svg';
import DatabaseIcon from '~/components/icons/DatabaseIcon.svg';
import ApplicationsIcon from '~/components/icons/ApplicationsIcon.svg';
import ExternalServicesIcon from '~/components/icons/ExternalServicesIcon.svg';
import LoginTypesIcon from '~/components/icons/LoginTypesIcon.svg';
import MultiFactorAuthIcon from '~/components/icons/MultiFactorAuthIcon.svg';
import HostedPagesIcon from '~/components/icons/HostedPagesIcon.svg';
import GearIcon from '~/components/icons/GearIcon.svg';
import AppleIcon from '~/components/icons/AppleIcon.svg';
import Section from './Section';

const sections = [
  {
    title: 'Server',
    path: '/management',
    exact: true,
    Icon: GearIcon,
  },
  {
    title: 'Account',
    path: '/management/account',
    Icon: AdminIcon,
  },
  {
    title: 'Users',
    path: '/management/users',
    Icon: UsersIcon,
  },
  {
    title: 'Applications',
    path: '/management/applications',
    Icon: ApplicationsIcon,
  },
  {
    title: 'Storages',
    path: '/management/database',
    Icon: DatabaseIcon,
  },
  {
    title: 'Login Settings',
    path: '/management/settings',
    Icon: LoginTypesIcon,
  },
  {
    title: 'External Services',
    path: '/management/email_integration',
    Icon: ExternalServicesIcon,
  },
  {
    title: 'Static Files',
    path: '/management/static',
    Icon: StaticFilesIcon,
  },
  {
    title: 'Apple Integration',
    path: '/management/apple',
    Icon: AppleIcon,
  },
  {
    title: 'Multi-factor Auth',
    path: '/management/multi-factor_auth',
    Icon: MultiFactorAuthIcon,
    disabled: true,
  },
  {
    title: 'Hosted Pages',
    path: '/management/hosted_pages',
    Icon: HostedPagesIcon,
  },
];

const Sidebar = () => (
  <aside className="iap-management-sidebar">
    {sections.map(section => <Section key={section.title} {...section} />)}
  </aside>
);

export default Sidebar;
